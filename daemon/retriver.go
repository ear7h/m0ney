package main

import (
	"time"
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"strings"
	"m0ney/data"
)

const (
	QUOTES_URL = "https://api.robinhood.com/quotes/?symbols="
	NASDAQ_HOURS_URL = "https://api.robinhood.com/markets/XNAS/hours/"
)

var SYMBOLS []string
var errcount = 0

func init() {
	config, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	var m map[string][]string
	err = json.Unmarshal(config, &m)
	if err != nil {
		panic(err)
	}

	SYMBOLS = m["symbols"]

	fmt.Println("retriever initiated")
}

func insertPrices() {

	addr := QUOTES_URL + strings.Join(SYMBOLS, ",")
	fmt.Println(addr)
	res, err := http.Get(addr)
	if err != nil {
		fmt.Println(err)
		errcount ++

	}

	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}


	var dat map[string][]data.RhQuote

	err = json.Unmarshal(resData, &dat)
	if err != nil {
		panic(err)
	}

	for _, v := range dat["results"] {
		err := data.InsertRhQuote(v)
		if err != nil {
			panic(err)
		}
	}
}

func getMarketHours() (time.Time, time.Time) {

	//get today's market info
	res, err := http.Get(NASDAQ_HOURS_URL + time.Now().Format("2006-01-02") + "/")
	if err != nil {
		panic(err)
	}

	byt, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var m map[string]interface{}

	err = json.Unmarshal(byt, &m)
	if err != nil {
		panic(err)
	}


	//if it's open today and the close time is in the future
	//then assign a close time and return when open
	if m["is_open"].(bool) {
		openTime, err := time.Parse(time.RFC3339, m["opens_at"].(string))
		closeTime, err := time.Parse(time.RFC3339, m["closes_at"].(string))
		if err != nil {
			panic(err)
		}

		if closeTime.After(time.Now()) {
			return openTime, closeTime
		}

	}

	nextOpenURL := m["next_open_hours"].(string)

	//get next open day market info
	res, err = http.Get(nextOpenURL)
	if err != nil {
		panic(err)
	}

	byt, err = ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(byt, &m)
	if err != nil {
		panic(err)
	}

	openTime, err := time.Parse(time.RFC3339, m["opens_at"].(string))
	closeTime, err := time.Parse(time.RFC3339, m["closes_at"].(string))
	if err != nil {
		panic(err)
	}

	return openTime, closeTime

}

func addDataSets(start, end time.Time, scale time.Duration, sym []string) {
	s := data.Dataset{
		Start: start,
		End: end,
		Table: "moment",
	}

	for _, v := range sym {
		s.Ticker = v
		data.InsertDataset(s)
	}

}

func dayLoop(start, end time.Time) {
	fmt.Println("market open at: ", start)

	if true {//start.Before(time.Now()) {
		go addDataSets(time.Now(), end, time.Second, SYMBOLS)
	} else {
		go addDataSets(start, end, time.Second, SYMBOLS)
		time.Sleep(time.Until(start))
	}
	fmt.Println("starting")

	for time.Now().Before(end) {
		time.Sleep(time.Second)
		insertPrices()
	}

	fmt.Println("finished")

}

//program loop
func momentRetriever() error {

	//day loop
	for start, end := getMarketHours(); true; {
		dayLoop(start, end)
	}

	return nil
}


