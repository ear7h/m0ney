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

func init() {
	var err error
	DB, err = sql.Open("mysql", "root@/stocks")
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
	rows, err := DB.Query("SELECT * FROM SETS;")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	arr := []Dataset{}

	for rows.Next() {
		v := Dataset{}
		rows.Scan(&v.Ticker, &v.Scale, &v.Start, &v.End)

		arr = append(arr, v)
	}

	return arr
}

func InsertMoment(v Moment) error {
	_, err := DB.Exec("INSERT INTO MOMENT (" +
		"ask_price, ask_size, bid_price, bid_size," +
		"last_trade_price, symbol, updated_at" +
		") values (?,?,?,?,?,?,?)",
		v.AskPrice, v.AskSize, v.BidPrice, v.BidSize,
		v.LastTradePrice, v.Symbol, v.UpdatedAt,
	)

	if err != nil {
		return err
	}
	return nil
}