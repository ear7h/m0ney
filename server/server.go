package server

import (
	"encoding/json"
	"fmt"
	"github.com/ear7h/m0ney/data"
	"net/http"
	"strconv"
	"strings"
)

var _db *data.MoneyDB

func init() {
	var err error
	_db, err = data.NewMoneyDB("root", "", "db", "3306", "money")
	if err != nil {
		panic(err)
	}
}

func handleList(w http.ResponseWriter, r *http.Request) {
	arr, err := _db.Runs()
	if err != nil {
		http.Error(w, "couldn't retrieve runs", http.StatusInternalServerError)
		return
	}

	byt, err := json.Marshal(arr)
	if err != nil {
		http.Error(w, "Error marshaling data run list", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(byt)
}

func handleRun(w http.ResponseWriter, r *http.Request) {
	arr := strings.Split(r.URL.Path, "/")
	if len(arr) != 3 {
		http.Error(w, "path must be /run/{integer >= 1}", http.StatusBadRequest)
		return
	}

	i, err := strconv.ParseInt(arr[2], 10, 64)
	if err != nil || i < 1 {
		http.Error(w, "path must be /run/{integer >= 1}", http.StatusBadRequest)
		return
	}

	run, err := _db.Run(int(i))
	if err != nil {
		http.Error(w, "could not retrieve run " + arr[2], http.StatusInternalServerError)
	}

	byt, err := json.Marshal(run)
	if err != nil {
		http.Error(w, "couldn't marshal response", http.StatusInternalServerError)
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

	//get the data set
	m.HandleFunc("/run/", handleRun)

	//start server
	//http.ListenAndServe() is a blocking call
	err := http.ListenAndServe(":8080", m)
	panic(err)
}
