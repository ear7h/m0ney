package daemon

import (
	"time"
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"strings"
	"m0ney/data"
	"os"
)

const (
	QUOTES_URL = "https://api.robinhood.com/quotes/?symbols="
	NASDAQ_HOURS_URL = "https://api.robinhood.com/markets/XNAS/hours/"
)

var SYMBOLS []string = []string{"AAL", "AAPL", "ADBE", "ADI", "ADP", "ADSK", "AKAM", "ALXN", "AMD", "AMAT", "AMGN", "AMZN", "ATVI", "AVGO", "BIDU", "BIIB", "BMRN", "CA", "CELG", "CERN", "CHKP", "CHTR", "CTRP", "CTAS", "CSCO", "CTXS", "CMCSA", "COST", "CSX", "CTSH", "DISCA", "DISCK", "DISH", "DLTR", "EA", "EBAY", "ESRX", "EXPE", "FAST", "FB", "FISV", "FOX", "FOXA", "GILD", "GOOG", "GOOGL", "HAS", "HSIC", "HOLX", "ILMN", "INCY", "INTC", "INTU", "ISRG", "JBHT", "JD", "KLAC", "KHC", "LBTYK", "LILA", "LBTYA", "QCOM", "QVCA", "MELI", "MAR", "MAT", "MDLZ", "MNST", "MSFT", "MU", "MXIM", "MYL", "NCLH", "NFLX", "NTES", "NVDA", "PAYX", "PCLN", "PYPL", "QCOM", "REGN", "ROST", "SHPG", "SIRI", "SWKS", "SBUX", "SYMC", "TSCO", "TXN", "TMUS", "ULTA", "VIAB", "VOD", "VRTX", "WBA", "WDC", "XRAY", "IDXX", "LILAK", "LRCX", "MCHP", "ORLY", "PCAR", "STX", "TSLA", "VRSK", "WYNN", "XLNX"}


func insertPrices() {

	addr := QUOTES_URL + strings.Join(SYMBOLS, ",")
	fmt.Println(addr)
	res, err := http.Get(addr)
	if err != nil {
		fmt.Fprint(os.Stderr,"could not get")
		fmt.Println(err)
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

func addDataSets(sym []string, d time.Duration) {


	for _, v := range sym {
		_, err := data.DB.Exec("INSERT INTO sets (`symbol`, `start`, `end`, `scale`, `table`) SELECT `symbol`, min(`updated_at`), max(`updated_at`), ?, 'moment' FROM moment WHERE `symbol` = ? AND DATE(`updated_at`) = DATE(NOW()) GROUP BY `symbol`;", d, v)
		if err != nil {
			panic(err)
		}
	}

}

func dayLoop(start, end time.Time) {
	fmt.Println("market open at: ", start)

	//add data set after completion of day loop
	defer func () {
		addDataSets(SYMBOLS, time.Second)
	}()

	fmt.Println("retriever starting")

	for time.Now().Before(end) {
		time.Sleep(time.Second)
		insertPrices()
	}

	fmt.Println("retriver finished")

}

//program loop
func momentRetriever() error {

	//day loop
	for start, end := getMarketHours(); true; {
		dayLoop(start, end)
	}

	return nil
}


