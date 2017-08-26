package data

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
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

func GetSets() []Dataset {
	rows, err := DB.Query("SELECT * FROM sets;")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	arr := []Dataset{}

	for rows.Next() {
		v := Dataset{}
		rows.Scan(&v.ID, &v.Symbol, &v.Scale, &v.Start, &v.End)

		arr = append(arr, v)
	}

	return arr
}

func InsertRhQuote(r RhQuote) error {
	v := r.ToMoment()

	_, err := DB.Exec("INSERT INTO MOMENT " +
		"(ask_price, ask_size, bid_price, bid_size," +
		"last_trade_price, symbol, updated_at) " +
		"VALUES (?,?,?,?,?,?,?);",
		v.AskPrice, v.AskSize, v.BidPrice, v.BidSize,
		v.LastTradePrice, v.Symbol, v.UpdatedAt,
	)

	if err != nil {
		return err
	}
	return nil
}

func InsertDataset(v Dataset) error {

	_, err := DB.Exec("INSERT INTO sets" +
		"(symbol, scale, start, end, `table`) " +
		"VALUES (?,?,?,?,?);",
		v.Symbol, int(v.Scale), v.Start, v.End, string(v.Table),
	)

	return err
}