package main

import (
	"encoding/json"
	"fmt"
	"github.com/ear7h/m0ney/data"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	QUOTES_URL       = "https://api.robinhood.com/quotes/?symbols="
	NASDAQ_HOURS_URL = "https://api.robinhood.com/markets/XNAS/hours/"
)

var _symbols []string = []string{"AAL", "AAPL", "ADBE", "ADI", "ADP", "ADSK", "AKAM", "ALXN", "AMD", "AMAT", "AMGN", "AMZN", "ATVI", "AVGO", "BIDU", "BIIB", "BMRN", "CA", "CELG", "CERN", "CHKP", "CHTR", "CTRP", "CTAS", "CSCO", "CTXS", "CMCSA", "COST", "CSX", "CTSH", "DISCA", "DISCK", "DISH", "DLTR", "EA", "EBAY", "ESRX", "EXPE", "FAST", "FB", "FISV", "FOX", "FOXA", "GILD", "GOOG", "GOOGL", "HAS", "HSIC", "HOLX", "ILMN", "INCY", "INTC", "INTU", "ISRG", "JBHT", "JD", "KLAC", "KHC", "LBTYK", "LILA", "LBTYA", "QCOM", "QVCA", "MELI", "MAR", "MAT", "MDLZ", "MNST", "MSFT", "MU", "MXIM", "MYL", "NCLH", "NFLX", "NTES", "NVDA", "PAYX", "PCLN", "PYPL", "QCOM", "REGN", "ROST", "SHPG", "SIRI", "SWKS", "SBUX", "SYMC", "TSCO", "TXN", "TMUS", "ULTA", "VIAB", "VOD", "VRTX", "WBA", "WDC", "XRAY", "IDXX", "LILAK", "LRCX", "MCHP", "ORLY", "PCAR", "STX", "TSLA", "VRSK", "WYNN", "XLNX"}

var _db *data.MoneyDB

func init() {
	var err error
	_db, err = data.NewMoneyDB("root", "", "db", "3306", "money")
	if err != nil {
		panic(err)
	}
}

func insertPrices() {
	addr := QUOTES_URL + strings.Join(_symbols, ",")
	// log.Println("getting ", addr)
	resp, err := http.Get(addr)
	if err != nil {
		log.Println("error error getting prices", err)
		return
	}
	defer resp.Body.Close()

	resData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("error reading body ", err)
		return
	}

	dat := struct {
		Results []rhQuote `json:"results"`
	}{}

	err = json.Unmarshal(resData, &dat)
	if err != nil {
		log.Println("error unmarshaling json", err)
		return
	}

	arr := make([]data.Quote, len(dat.Results))
	for k, v := range dat.Results {
		arr[k] = v.ToQuote()
	}

	err = _db.InsertQuotes(arr)
	if err != nil {
		log.Print()
	}
}

func getMarketHours() (time.Time, time.Time) {

	url := NASDAQ_HOURS_URL + time.Now().Format("2006-01-02") + "/"
	//tag for goto statement
L:
	//get today's market info
	resp, err := http.Get(url)
	if err != nil {
		log.Println("error getting request ", err)
		return time.Time{}, time.Time{}
	}
	defer resp.Body.Close()

	byt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("error reading body ", err)
		return time.Time{}, time.Time{}
	}

	var m map[string]interface{}

	err = json.Unmarshal(byt, &m)
	if err != nil {
		log.Println("error marhsaling json ", err)
		return time.Time{}, time.Time{}
	}

	//if the market is open on fetched day and the close time is in the future
	//then parse time strings
	if m["is_open"].(bool) {
		openTime, err := time.Parse(time.RFC3339, m["opens_at"].(string))
		closeTime, err := time.Parse(time.RFC3339, m["closes_at"].(string))
		if err != nil {
			log.Println("error parsing market hours", err)
			return time.Time{}, time.Time{}
		}

		if closeTime.After(time.Now()) {
			return openTime, closeTime
		}

	}
	//else get the next open hours
	url = m["next_open_hours"].(string)
	//and fetch again
	goto L
}

func dayLoop(start, end time.Time) {
	//add data set after completion of day loop
	defer _db.Nightly()

	//if the market opens in the future
	//then wait
	if time.Now().Before(start) {
		time.Sleep(start.Sub(time.Now()))
	}

	fmt.Println("retriever starting")

	//fetch and insert prices until the close time is in the past
	for time.Now().Before(end) {
		insertPrices()
		time.Sleep(time.Second)
	}

	fmt.Println("retriver finished")
}

//program entry point
func main() {
	//program loop
	for true {
		s, e := getMarketHours()

		if (time.Time{} == s) || (time.Time{} == e) {
			log.Println("could not get market hrs trying again in 10 seconds")
			time.Sleep(10 * time.Second)
			s, e = getMarketHours()
		}

		log.Println("market hrs: ", s, " - ", e)

		dayLoop(s, e)
	}
}
