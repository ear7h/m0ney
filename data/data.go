package data

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"time"
)

const (
	SQL_TIME = "2006-01-02 15:04:05"
)

var DB *sql.DB

//open db
func init() {
	m := mysqlCreds

	//user, password, host, port, database
	Open(m["user"], m["password"], m["host"], m["port"], m["database"])
}

func Open(user, password, host, port, database string) {

	url := fmt.Sprint(user, ":", password, "@(", host, ":", port, ")/", database)

	fmt.Println(url)

	var err error
	DB, err = sql.Open("mysql", url)
	if err != nil {
		panic(err)
	}

	err = DB.Ping()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("database connection successful")
	}
}

func GetSets() []Set {
	rows, err := DB.Query("SELECT `id`, `symbol`, `start`, `end`, `scale`, `table` FROM sets;")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	arr := []Set{}

	for rows.Next() {
		var (
			id int
			symbol string
			start string
			end string
			scale int
			table string
		)

		rows.Scan(&id, &symbol, &start, &end, &scale, &table)

		ts, _ := time.Parse(SQL_TIME, start)
		te, _ := time.Parse(SQL_TIME, end)

		v := Set{
			ID: id,
			Symbol: symbol,
			Start: ts,
			End: te,
			Scale: time.Duration(scale),
			Table: Table(table),
		}

		arr = append(arr, v)
	}

	return arr
}

func InsertRhQuote(r RhQuote) error {
	v := r.ToMoment()

	_, err := DB.Exec("INSERT INTO moment " +
		"(`ask_price`, `ask_size`, `bid_price`, `bid_size`," +
		"`last_trade_price`, `symbol`, `trading_halted`, `updated_at`) " +
		"VALUES (?,?,?,?,?,?,?,?);",
		v.AskPrice, v.AskSize, v.BidPrice, v.BidSize,
		v.LastTradePrice, v.Symbol, v.TradingHalted,v.UpdatedAt,
	)

	if err != nil {
		return err
	}
	return nil
}