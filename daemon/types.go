package main

import (
	"github.com/ear7h/m0ney/data"
	"strconv"
	"time"
)

type rhQuote struct {
	AskPrice       string    `json:"ask_price"`
	AskSize        int       `json:"ask_size"`
	BidPrice       string    `json:"bid_price"`
	BidSize        int       `json:"bid_size"`
	LastTradePrice string    `json:"last_trade_price"`
	Symbol         string    `json:"symbol"`
	TradingHalted  bool      `json:"trading_halted"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (r rhQuote) ToQuote() data.Quote {

	askPrice, err := strconv.ParseFloat(r.AskPrice, 64)
	bidPrice, err := strconv.ParseFloat(r.BidPrice, 64)
	lastTradePrice, err := strconv.ParseFloat(r.LastTradePrice, 64)

	if err != nil {
		panic(err)
	}

	return data.Quote{
		AskPrice:       askPrice,
		AskSize:        r.AskSize,
		BidPrice:       bidPrice,
		BidSize:        r.BidSize,
		LastTradePrice: lastTradePrice,
		Symbol:         r.Symbol,
		TradingHalted:  r.TradingHalted,
		UpdatedAt:      r.UpdatedAt,
	}
}
