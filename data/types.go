package data

import (
	"time"
)

//data sets are continuous on the specified time interval. the moment table has an interval of 1 minute

// A partition is a division of the same data across multiple tables
// in this case, each partition will be one week's data
type Partition struct {
	ID     int
	Name   string
	WeekOf time.Time
}

type Run struct {
	ID            int           `json:"id"`
	Symbol        string        `json:"symbol"`
	Start         time.Time     `json:"start"`
	End           time.Time     `json:"end"`
	PartitionName string `json:"partition_id"`
}

type Quote struct {
	AskPrice       float64   `json:"ask_price"`
	AskSize        int       `json:"ask_size"`
	BidPrice       float64   `json:"bid_price"`
	BidSize        int       `json:"bid_size"`
	LastTradePrice float64   `json:"last_trade_price"`
	Symbol         string    `json:"symbol"`
	TradingHalted  bool      `json:"trading_halted"`
	UpdatedAt      time.Time `json:"updated_at"`
}
