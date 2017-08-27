package main

import (
	"net/http"
	"time"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"m0ney/data"
	"m0ney/sessions"
	"database/sql"
	"m0ney/daemon"
)

func handleList(w http.ResponseWriter, r *http.Request) {
	arr := data.GetSets()

	byt, err := json.Marshal(arr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error getting data sets"))
	}

	w.WriteHeader(http.StatusOK)
	w.Write(byt)
}

func handleSessionCreate(w http.ResponseWriter, r *http.Request) {
	//post request creates session
	if r.Method != http.MethodPost {
		http.Error(w, "send a POST request with dataset id from /list or get request to /session/{token}", http.StatusMethodNotAllowed)
		return
	}

	byt, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "could not read body", http.StatusInternalServerError)
		return
	}

	req := data.Dataset{}
	err = json.Unmarshal(byt, &req)
	if err != nil {
		http.Error(w, "could not parse body", http.StatusBadRequest)
		return
	}

	//the session request can specify and id or everything else
	//if it only specifies the id we need to fill in everything else
	if req.ID == 0 {
		rows, err := data.DB.Query("SELECT `symbol`, `start`, `end`, `scale`, `table` FROM sets WHERE id = ?;", req.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "specified id could not be found", http.StatusNotFound)
			} else {
				panic(err)
			}

			return
		}
		defer rows.Close()

		//scan values from query and fill (replace) req with data.Dataset values
		var (
			symbol string
			start string
			end string
			scale int
			table string
		)

		rows.Scan(&symbol, &start, &end, &scale, &table)

		ts, _ := time.Parse(data.SQL_TIME, start)
		te, _ := time.Parse(data.SQL_TIME, end)

		req = data.Dataset{
			Symbol: symbol,
			Start: ts,
			End: te,
			Scale: time.Duration(scale),
			Table: data.Table(table),
		}
	}


	//create session from request
	token := sessions.Create(&sessions.Session{
		SessStart: time.Now(),
		Scale: req.Scale,
		CurrentTime: req.Start,
		Ticker: req.Symbol,
	})


	v := map[string]string{"token": token}

	retByt, err := json.Marshal(v)
	if err != nil {
		http.Error(w, "could not marshal response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(retByt)

}

func handleSession(w http.ResponseWriter, r *http.Request)  {

	if r.Method != http.MethodGet {
		http.Error(w, "must be get request", http.StatusMethodNotAllowed)
		return
	}

	token := r.URL.Path[len("/session/"):]

	sess := sessions.Get(token)

	if sess == nil {
		http.Error(w, "session " + token +" doesn't exist", http.StatusNotFound)
		return
	}

	ret := sess.Next()

	byt, err := json.Marshal(ret)
	if err != nil {
		http.Error(w, "could not marshal response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(byt)
}

//entry point for server
func main() {
	fmt.Println("starting")

	m := http.NewServeMux()

	//list available data sets
	m.HandleFunc("/list", handleList)

	//session
	//get (with auth) - returns values
	//post - returns key and value for session
	m.HandleFunc("/session", handleSessionCreate)
	m.HandleFunc("/session/", handleSession)

	//start daemon
	go daemon.Main()

	//start server
	//http.ListenAndServe() is a blocking call
	err := http.ListenAndServe(":8080", m)
	panic(err)
}
