package data

import (
	"time"
	"strconv"
)

type Table string

const (
	TABLE_MOMENT     = Table("moment")
	TABLE_HISTORICAL = Table("historical")
)

//data sets are continuous on the specified time interval. the moment table has an interval of 1 minute
type Set struct {
	ID     int           `json:"id"`
	Symbol string        `json:"symbol"`
	Start  time.Time     `json:"start"`
	End    time.Time     `json:"end"`
	Scale  time.Duration `json:"scale"`
	Table  Table         `json:"table"`
}

type Moment struct {
	AskPrice          float64   `json:"ask_price"`
	AskSize           int       `json:"ask_size"`
	BidPrice          float64   `json:"bid_price"`
	BidSize           int       `json:"bid_size"`
	LastTradePrice    float64   `json:"last_trade_price"`
	Symbol            string    `json:"symbol"`
	TradingHalted     bool      `json:"trading_halted"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type RhQuote struct {
	AskPrice          string    `json:"ask_price"`
	AskSize           int       `json:"ask_size"`
	BidPrice          string    `json:"bid_price"`
	BidSize           int       `json:"bid_size"`
	LastTradePrice    string    `json:"last_trade_price"`
	Symbol            string    `json:"symbol"`
	TradingHalted     bool      `json:"trading_halted"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (r RhQuote) ToMoment() Moment {

	askPrice, err := strconv.ParseFloat(r.AskPrice, 64)
	bidPrice, err := strconv.ParseFloat(r.BidPrice, 64)
	lastTradePrice, err := strconv.ParseFloat(r.LastTradePrice, 64)

	if err != nil {
		panic(err)
	}

	return Moment{
		AskPrice: askPrice,
		AskSize: r.AskSize,
		BidPrice: bidPrice,
		BidSize: r.BidSize,
		LastTradePrice: lastTradePrice,
		Symbol: r.Symbol,
		TradingHalted: r.TradingHalted,
		UpdatedAt: r.UpdatedAt,

	}
}
