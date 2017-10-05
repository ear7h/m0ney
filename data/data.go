package data

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"m0ney/log"
	"os"
)

const (
	SQL_TIME = "2006-01-02 15:04:05"
)

var DB *sql.DB

//open db
func init() {
	user := "root"
	password := ""
	host := ""
	port :=     "3306"
	database := "stocks"

	if os.Getenv("EAR7H_ENV") == "prod" {
			host =  "db"
	}

	url := fmt.Sprint(user, ":", password, "@(", host, ":", port, ")/", database)

	fmt.Println(url)

	var err error
	DB, err = sql.Open("mysql", url)
	if err != nil {
		panic(err)
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
			id     int
			symbol string
			start  string
			end    string
			scale  int
			table  string
		)

		rows.Scan(&id, &symbol, &start, &end, &scale, &table)

		ts, _ := time.Parse(SQL_TIME, start)
		te, _ := time.Parse(SQL_TIME, end)

		v := Set{
			ID:     id,
			Symbol: symbol,
			Start:  ts,
			End:    te,
			Scale:  time.Duration(scale),
			Table:  Table(table),
		}

		arr = append(arr, v)
	}

	return arr
}

func GetSet(i int) (ret []Moment) {

	rows, err := DB.Query("SELECT symbol, date(start) FROM sets WHERE id = ?", i)
	if err != nil {
		log.Enter(3, err)
		return
	}

	rows.Next()

	var sym string
	var start string

	err = rows.Scan(&sym, &start)
	if err != nil {
		log.Enter(3, err)
	}
	rows.Close()

	rows, err = DB.Query(`SELECT
	avg(ask_price) as ask_price, avg(ask_size) as ask_size,
	avg(bid_price) as bid_price, avg(bid_size) as bid_size,
    avg(last_trade_price) as last_trade_price,
    ? as symbol,
    min(updated_at) as updated_at
FROM moment
WHERE
	symbol = ?
	AND date(updated_at) = date(?)
GROUP BY hour(updated_at) asc, minute(updated_at) asc;`, sym, sym, start)
	if err != nil {
		log.Enter(3, err)
		return []Moment{}
	}
	defer rows.Close()


	for rows.Next() {
		var (
			askPrice       float64
			askSize        float64
			bidPrice       float64
			bidSize        float64
			lastTradePrice float64
			symbol         string
			updatedAt      string
		)

		err = rows.Scan(
			&askPrice,
			&askSize,
			&bidPrice,
			&bidSize,
			&lastTradePrice,
			&symbol,
			&updatedAt,
		)
		if err != nil {
			log.Enter(3, err)
			return
		}

		ua, _ := time.Parse(SQL_TIME, updatedAt)

		ret = append(ret, Moment{
			AskPrice:       askPrice,
			AskSize:        int(askSize),
			BidPrice:       bidPrice,
			BidSize:        int(bidSize),
			LastTradePrice: lastTradePrice,
			Symbol:         symbol,
			UpdatedAt:      ua,
		})
	}

	return ret
}

func InsertRhQuote(r RhQuote) error {
	v := r.ToMoment()

	_, err := DB.Exec("INSERT INTO moment "+
		"(`ask_price`, `ask_size`, `bid_price`, `bid_size`,"+
		"`last_trade_price`, `symbol`, `trading_halted`, `updated_at`) "+
		"VALUES (?,?,?,?,?,?,?,?);",
		v.AskPrice, v.AskSize, v.BidPrice, v.BidSize,
		v.LastTradePrice, v.Symbol, v.TradingHalted, v.UpdatedAt,
	)

	if err != nil {
		return err
	}
	return nil
}
