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
type Dataset struct {
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
	PreviousCloseDate time.Time `json:"previous_close_date"`
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
	PreviousCloseDate string    `json:"previous_close_date"`
	Symbol            string    `json:"symbol"`
	TradingHalted     bool      `json:"trading_halted"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (r RhQuote) ToMoment() Moment {

	askPrice, err := strconv.ParseFloat(r.AskPrice, 64)
	bidPrice, err := strconv.ParseFloat(r.BidPrice, 64)
	lastTradePrice, err := strconv.ParseFloat(r.LastTradePrice, 64)
	previousCloseDate, err := time.Parse("2006-01-02", r.PreviousCloseDate)
	if err != nil {
		panic(err)
	}

	return Moment{
		AskPrice: askPrice,
		AskSize: r.AskSize,
		BidPrice: bidPrice,
		BidSize: r.BidSize,
		LastTradePrice: lastTradePrice,
		PreviousCloseDate: previousCloseDate,
		Symbol: r.Symbol,
		TradingHalted: r.TradingHalted,
		UpdatedAt: r.UpdatedAt,

	}
}
