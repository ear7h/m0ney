package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"github.com/ear7h/m0ney/data"
	"github.com/ear7h/m0ney/log"
	"net/http"
	"strings"
	"time"
)

const (
	QUOTES_URL       = "https://api.robinhood.com/quotes/?symbols="
	NASDAQ_HOURS_URL = "https://api.robinhood.com/markets/XNAS/hours/"
)

var SYMBOLS []string = []string{"AAL", "AAPL", "ADBE", "ADI", "ADP", "ADSK", "AKAM", "ALXN", "AMD", "AMAT", "AMGN", "AMZN", "ATVI", "AVGO", "BIDU", "BIIB", "BMRN", "CA", "CELG", "CERN", "CHKP", "CHTR", "CTRP", "CTAS", "CSCO", "CTXS", "CMCSA", "COST", "CSX", "CTSH", "DISCA", "DISCK", "DISH", "DLTR", "EA", "EBAY", "ESRX", "EXPE", "FAST", "FB", "FISV", "FOX", "FOXA", "GILD", "GOOG", "GOOGL", "HAS", "HSIC", "HOLX", "ILMN", "INCY", "INTC", "INTU", "ISRG", "JBHT", "JD", "KLAC", "KHC", "LBTYK", "LILA", "LBTYA", "QCOM", "QVCA", "MELI", "MAR", "MAT", "MDLZ", "MNST", "MSFT", "MU", "MXIM", "MYL", "NCLH", "NFLX", "NTES", "NVDA", "PAYX", "PCLN", "PYPL", "QCOM", "REGN", "ROST", "SHPG", "SIRI", "SWKS", "SBUX", "SYMC", "TSCO", "TXN", "TMUS", "ULTA", "VIAB", "VOD", "VRTX", "WBA", "WDC", "XRAY", "IDXX", "LILAK", "LRCX", "MCHP", "ORLY", "PCAR", "STX", "TSLA", "VRSK", "WYNN", "XLNX"}

func insertPrices() {
	addr := QUOTES_URL + strings.Join(SYMBOLS, ",")
	log.Enter(log.OK, addr)
	resp, err := http.Get(addr)
	if err != nil {
		log.Enter(log.ERROR, err)
		return
	}
	defer resp.Body.Close()

	resData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Enter(log.ERROR, err)
		return
	}

	var dat map[string][]data.RhQuote

	err = json.Unmarshal(resData, &dat)
	if err != nil {
		log.Enter(log.ERROR, err)
		return
	}

	for _, v := range dat["results"] {
		err := data.InsertRhQuote(v)
		if err != nil {
			log.Enter(log.ERROR, err)
			return
		}
	}
}

func getMarketHours() (time.Time, time.Time) {

	url := NASDAQ_HOURS_URL + time.Now().Format("2006-01-02") + "/"
	//tag for goto statement
L:
	//get today's market info
	resp, err := http.Get(url)
	if err != nil {
		log.Enter(log.ERROR, err)
		return time.Time{}, time.Time{}
	}
	defer resp.Body.Close()

	byt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Enter(log.ERROR, err)
		return time.Time{}, time.Time{}
	}

	var m map[string]interface{}

	err = json.Unmarshal(byt, &m)
	if err != nil {
		log.Enter(log.ERROR, err)
		return time.Time{}, time.Time{}
	}

	//if the market is open on fetched day and the close time is in the future
	//then parse time strings
	if m["is_open"].(bool) {
		openTime, err := time.Parse(time.RFC3339, m["opens_at"].(string))
		closeTime, err := time.Parse(time.RFC3339, m["closes_at"].(string))
		if err != nil {
			log.Enter(log.ERROR, err)
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

func insertDataSets(d time.Duration) {

	_, err := data.DB.Exec("INSERT INTO sets " +
		"(`symbol`, `start`, `end`, `scale`, `table`) " +
		"SELECT `symbol`, min(`updated_at`), max(`updated_at`), ?, 'moment' FROM `moment` " +
		"WHERE DATE(`updated_at`) = DATE(NOW()) GROUP BY `symbol`, DATE(`updated_at`) ;", d)

	if err != nil {
		log.Enter(log.ERROR, err)
		return
	}

}

func dayLoop(start, end time.Time) {
	//add data set after completion of day loop
	defer insertDataSets(time.Second)

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

//checks config to see if the retriever should run another loop
func shouldRunToday() (retrieveAnother bool) {
	byt, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Enter(log.ERROR, err)
		return false
	}

	v := map[string]bool{}

	err = json.Unmarshal(byt, &v)
	if err != nil {
		log.Enter(log.ERROR, err)
		return false
	}


	retrieveAnother = v["retrieveAnother"]
	log.Enter(log.OK, "retrieveAnother = ", retrieveAnother)

	return

}

//program entry point
func momentRetriever() error {

	//program loop
	for true {
		s, e := getMarketHours()

		for (time.Time{} == s) || (time.Time{} == e) {
			log.Enter(log.WARNING, "could not market hrs trying again in 10 seconds")
			time.Sleep(10 * time.Second)
			s, e = getMarketHours()
		}

		log.Enter(log.OK, "market hrs: ", s, " - ", e)

		if !shouldRunToday() {
			time.Sleep(time.Until(e))
			continue
		}

		dayLoop(s, e)
	}

	return nil
}
