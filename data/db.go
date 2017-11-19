package data

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	SQL_TIME = "2006-01-02 15:04:05"
)

type MoneyDB struct {
	*sql.DB
	lastPart  string
	lastStage string
}

func NewMoneyDB(user, pass, host, port, schema string) (mdb *MoneyDB, err error) {
	url := fmt.Sprintf("%s:%s@(%s:%s)/%s", user, pass, host, port, schema)
	url += "?multiStatements=true&parseTime=true"
	fmt.Println("connecting to: ", url)
	db, err := sql.Open("mysql", url)
	if err != nil {
		return
	}

	err = db.Ping()
	for err != nil {
		fmt.Println("db ping fail, trying again in 20s...")
		time.Sleep(20 * time.Second)
		err = db.Ping()
	}

	mdb = &MoneyDB{DB: db, lastPart: "", lastStage: ""}
	return
}

func (m *MoneyDB) Partitions() (arr []Partition, err error) {

	var rows *sql.Rows
	rows, err = m.Query(`SELECT count(*) FROM	partitions;
SELECT id, name, week_of FROM partitions;
`)
	if err != nil {
		return
	}
	defer rows.Close()

	// get the count to make a set length array
	var count int
	for rows.Next() {
		if err = rows.Scan(&count); err != nil {
			return
		}
	}

	if !rows.NextResultSet() {
		return
	}

	arr = make([]Partition, count)

	for i := 0; rows.Next(); i++ {
		var id int
		var name string
		var ts time.Time

		rows.Scan(&id, &name, &ts)
		arr[i] = Partition{
			ID:     id,
			Name:   name,
			WeekOf: ts,
		}
	}
	return
}

func (m *MoneyDB) Runs() (arr []Run, err error) {
	var rows *sql.Rows
	rows, err = m.Query(`
SELECT count(*) FROM runs;
SELECT id, symbol, start, end, partition_name FROM runs;`)
	if err != nil {
		fmt.Println("couldn't get runs", err)
		return
	}
	defer rows.Close()

	// get the count to make a set length array
	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return
		}
	}

	if count == 0 {
		fmt.Println("no rows in runs table")
		return
	}

	if !rows.NextResultSet() {
		fmt.Println("expected a second result set")
		return
	}

	arr = make([]Run, count)

	for i := 0; rows.Next(); i++ {
		v := Run{}

		err = rows.Scan(&(v.ID), &(v.Symbol), &(v.Start), &(v.End), &(v.PartitionName))
		if err != nil {
			fmt.Println("error scanning run rows\n", err)
			return
		}

		arr[i] = v
	}

	return
}

// In memory tables allow for temporary performance boosts,
// namely the UNIQUE directive to filter duplicate rows
func (m *MoneyDB) NewMemoryTable(name string) (actualName string, err error) {
	tStr := time.Now().Format("2006_01_02")

	actualName = fmt.Sprintf("mem_%s_%s", tStr, name)

	q := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
  ask_price double DEFAULT NULL,
  ask_size int(11) DEFAULT NULL,
  bid_price double DEFAULT NULL,
  bid_size int(11) DEFAULT NULL,
  last_trade_price double DEFAULT NULL,
  symbol varchar(8) NOT NULL,
  trading_halted tinyint(1) DEFAULT NULL,
  updated_at datetime NOT NULL,
  UNIQUE (ask_price, ask_size, bid_price, bid_size,   last_trade_price,  symbol, trading_halted, updated_at)
) ENGINE = MEMORY;`, actualName)

	_, err = m.Exec(q)
	return
}

// A partition table is a regular table which holds data by the week
//
func (m *MoneyDB) NewPartitionTable(t time.Time) (name string, err error) {
	y, w := t.ISOWeek()
	name = fmt.Sprintf("part_%d%d", y, w)

	q := fmt.Sprintf(`INSERT INTO partitions (name, week_of) VALUE ('%s', date('%s'));
CREATE TABLE %s (
  ask_price double DEFAULT NULL,
  ask_size int(11) DEFAULT NULL,
  bid_price double DEFAULT NULL,
  bid_size int(11) DEFAULT NULL,
  last_trade_price double DEFAULT NULL,
  symbol varchar(8) NOT NULL,
  trading_halted tinyint(1) DEFAULT NULL,
  updated_at datetime NOT NULL
) ENGINE = InnoDB;`, name, t.Format(SQL_TIME), name)
	_, err = m.Exec(q)
	return
}

// Insert quotes assumes the passed quotes are the latest pulled quotes
// MEANING: it only works real time
func (m *MoneyDB) InsertQuotes(q []Quote) (err error) {
	tblName := time.Now().Format("mem_2006_01_02_stage")

	// if the last stage created wasn't created today
	// transfer it to a partition, create today's stage and set it as today's stage
	if m.lastStage != tblName {
		go m.transferStageToPartition(m.lastStage)

		_, err = m.NewMemoryTable("stage")
		if err != nil {
			fmt.Println("error making mem table\n", err)
			return
		}
		m.lastStage = tblName
	}

	query := fmt.Sprintf(`INSERT IGNORE INTO %s
(ask_price, ask_size, bid_price, bid_size,
last_trade_price, symbol, trading_halted, updated_at)
VALUES `, m.lastStage)

	for _, v := range q {
		query += fmt.Sprintf("(%f, %d, %f, %d, %f, '%s', %t, '%s'),", v.AskPrice, v.AskSize, v.BidPrice, v.BidSize,
			v.LastTradePrice, v.Symbol, v.TradingHalted, v.UpdatedAt.Format(SQL_TIME))
	}

	byt := []byte(query)
	byt[len(query)-1] = ';'

	_, err = m.Exec(string(byt))

	return
}

// 1) finds the partition for the stage's date
// 2) registers the staged runs in the `runs` table
// 3) moves all staged rows into the fetched partition(1)
// 4) drops the stage table
// It will only work with auto generated stage tables
// named with the following format "mem_2006_01_02_stage"
//
func (m *MoneyDB) transferStageToPartition(stageName string) (name string) {
	t, err := time.Parse("mem_2006_01_02_stage", stageName)
	if err != nil {
		fmt.Println("error parsing last stage\n", err)
		return
	}
get_part:
	// get the name of the partition for the stage's date
	rows, err := m.Query(`SELECT name
FROM partitions as p
WHERE p.week_of = date(?)`, t.Format("2006-01-02"))
	if err != nil {
		fmt.Println("couldn't retrieve name\n", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&name)
		if err != nil {
			fmt.Println("couldn't scan name\n", err)
			return
		}
	}

	// partition hasn't been created
	if name == "" {
		fmt.Println("partition does not exist, creating...")
		_, err = m.NewPartitionTable(t)
		if err != nil {
			fmt.Println("coulndn't make partition table\n", err)
			return
		}
		goto get_part
	}

	q := fmt.Sprintf(`INSERT INTO runs(symbol, start, end, partition_name) (SELECT symbol, min(updated_at), max(updated_at), '%s' as partition_name from %s GROUP BY symbol)`, name, stageName)

	// register staged runs in the runs table
	_, err = m.Exec(q)
	if err != nil {
		fmt.Println("couldn't trasnfer from stage\n", err)
		return
	}

	// transfer to the partition
	q = fmt.Sprintf(`INSERT INTO %s (SELECT * FROM %s);`, name, stageName)
	_, err = m.Exec(q)
	if err != nil {
		fmt.Println("couldn't trasnfer from stage\n", err)
		return
	}

	// drop the stage
	q = fmt.Sprintf(`DROP TABLE %s;`, stageName)
	_, err = m.Exec(q)
	if err != nil {
		fmt.Println("couldn't trasnfer from stage\n", err)
		return
	}

	return
}

func (m *MoneyDB) Nightly() {
	m.lastPart = m.transferStageToPartition(m.lastStage)
	m.lastStage = ""
}

func (m *MoneyDB) Run(i int) (run []Quote, err error) {
	var rows *sql.Rows
	q := fmt.Sprintf(`SELECT symbol, partition_name, start, end FROM runs WHERE id = %d;`, i)
	rows, err = m.Query(q)
	if err != nil {
		return
	}

	var symbol, partitionName string
	var start, end time.Time

	for rows.Next() {
		err = rows.Scan(&symbol, &partitionName, &start, &end)
		if err != nil {
			fmt.Println("couldn't scan row")
			return
		}
	}
	rows.Close()

	startStr, endStr := start.Format(SQL_TIME), end.Format(SQL_TIME)

	q = fmt.Sprintf(`
SELECT count(*)
FROM %s
WHERE updated_at >= timestamp('%s') AND updated_at <= timestamp('%s') AND symbol = '%s';
SELECT ask_price, ask_size, bid_price, bid_size, last_trade_price, symbol, trading_halted, updated_at
FROM %s
WHERE updated_at >= timestamp('%s') AND updated_at <= timestamp('%s') AND symbol = '%s';`,
		partitionName, startStr, endStr, symbol,
		partitionName, startStr, endStr, symbol)

	rows, err = m.Query(q)
	if err != nil {
		return
	}

	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return
		}
	}

	if !rows.NextResultSet() {
		return
	}

	run = make([]Quote, count)
	for i := 0; rows.Next(); i++ {
		v := Quote{}

		err = rows.Scan(&(v.AskPrice), &(v.AskSize),
			&(v.BidPrice), &(v.BidSize), &(v.LastTradePrice),
			&(v.Symbol), &(v.TradingHalted), &(v.UpdatedAt))
		if err != nil {
			return
		}

		run[i] = v
	}

	return
}
