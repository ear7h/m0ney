package main

import (
	"net/http"
	"time"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"m0ney/data"
	"m0ney/sessions"
)


//open db
func init() {
	byt, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	var m map[string]string

	err = json.Unmarshal(byt, &m)
	if err != nil {
		panic(err)
	}
	//user, password, host, port, database
	data.Open(m["user"], m["password"], m["host"], m["port"], m["database"])
}

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

func sessionFromDataset(req data.Dataset) (string, error) {

	token := sessions.Create(&sessions.Session{
		SessStart: time.Now(),
		CurrentTime: req.Start,
		EndTime: req.End,
		Scale: req.Scale,
		Ticker: req.Ticker,
		Table: req.Table,
	})

	return token, nil
}

func handleSessionCreate(w http.ResponseWriter, r *http.Request) {
	//post request creates session
	if r.Method != http.MethodPost {
		http.Error(w, "send a POST request with dataset from /list or get request to /session/{token}", http.StatusMethodNotAllowed)
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


	//create session from request
	token := sessions.Create(&sessions.Session{
		SessStart: time.Now(),
		Scale: req.Scale,
		CurrentTime: req.Start,
		Ticker: req.Ticker,
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


	http.ListenAndServe(":8080", m)
}
