package main

import (
	"net/http"
	"time"
	"encoding/json"
	"fmt"
	"crypto/rand"
	mrand "math/rand"
	"encoding/base64"
	"io/ioutil"
)

//dummy funcion
func getSets() []dataset {
	arr := []dataset{
		{
			Ticker: "AAPL",
			Start: time.Now(),
			End: time.Now().Add(time.Second),
		}, {
			Ticker: "AAPL",
			Start: time.Now().Add(20 * time.Hour),
			End: time.Now().Add(25 * time.Hour),
		},
	}

	return arr
}

func getMoment(t time.Time, p time.Duration) float32 {
	return mrand.Float32()
}


//data sets are continuous on the specified time interval. the moment table has an interval of 1 minute
type dataset struct {
	Ticker string `json:"ticker"`
	Interval time.Duration `json:"interval"`
	Start time.Time `json:"start"`
	End time.Time `json:"end"`
}

//this represents a practice/training session
type session struct {
	sessStart time.Time
	CurrentTime time.Time `json:"current_time"`
	Interval time.Duration `json:"interval"`
	Ticker string `json:"ticker"`
	BidPrice float32 `json:"bid_price"`
}

func (s *session) next() session {
	s.BidPrice = getMoment(s.CurrentTime, s.Interval)

	s.CurrentTime = s.CurrentTime.Add(s.Interval)

	return *s
}

var SESSIONS map[string]*session = map[string]*session{}

func handleList(w http.ResponseWriter, r *http.Request) {
	arr := getSets()

	byt, err := json.Marshal(arr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error getting data sets"))
	}

	w.WriteHeader(http.StatusOK)
	w.Write(byt)
}

func sessionFromDataset(req dataset) (string, error) {
	token := make([]byte, 10)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	tokenStr := base64.StdEncoding.EncodeToString(token)


	for SESSIONS[tokenStr] != nil {
		_, err = rand.Read(token)
		if err != nil {
			return "", err
		}

		tokenStr = base64.StdEncoding.EncodeToString(token)
	}

	SESSIONS[tokenStr] = &session{
		sessStart: time.Now(),
		Interval: req.Interval,
		CurrentTime: req.Start,
		Ticker: req.Ticker,
	}

	return tokenStr, nil
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

	req := dataset{}
	err = json.Unmarshal(byt, &req)
	if err != nil {
		http.Error(w, "could not parse body", http.StatusBadRequest)
		return
	}

	str, err := sessionFromDataset(req)
	if err != nil {
		http.Error(w, "session could not be created", http.StatusInternalServerError)
		return
	}

	v := struct {
		Token string `json:"token"`
	}{
		Token: str,
	}

	fmt.Println(str)

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
	}

	token := r.URL.Path[len("/session/"):]

	if SESSIONS[token] == (&session{}) {
		http.Error(w, "session " + token +" doesn't exist", http.StatusNotFound)
		return
	}

	ret := SESSIONS[token].next()

	byt, err := json.Marshal(ret)
	if err != nil {
		http.Error(w, "could not marshal response", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(byt)

}

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
