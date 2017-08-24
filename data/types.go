package data

import (
	"time"
)

type Table string

const (
	TABLE_MOMENT     = Table("moment")
	TABLE_HISTORICAL= Table("historical")
)

//data sets are continuous on the specified time interval. the moment table has an interval of 1 minute
type Dataset struct {
	ID     int           `json:"id"`
	Ticker string        `json:"ticker"`
	Scale  time.Duration `json:"interval"`
	Start  time.Time     `json:"start"`
	End    time.Time     `json:"end"`
	Table  Table         `json:"table"`
}

type Moment struct {
	AskPrice          float32   `json:"ask_price"`
	AskSize           int       `json:"ask_size"`
	BidPrice          float32   `json:"bid_price"`
	BidSize           int       `json:"bid_size"`
	LastTradePrice    float32   `json:"last_trade_price"`
	PreviousCloseDate time.Time `json:"previous_close_date"`
	Symbol            string    `json:"symbol"`
	UpdatedAt         time.Time `json:"updated_at"`
}
