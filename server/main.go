package main

import (
	"net/http"
	"encoding/json"
	"fmt"
	"m0ney/data"
	"strconv"
	"strings"
)

func handleList(w http.ResponseWriter, r *http.Request) {
	arr := data.GetSets()

	byt, err := json.Marshal(arr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error getting data sets"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(byt)
}

func handleSet(w http.ResponseWriter, r *http.Request) {
	arr := strings.Split(r.URL.Path, "/")
	if len(arr) != 3 {
		http.Error(w, "path must be /set/{integer > 1}", http.StatusBadRequest)
		return
	}


	i, err := strconv.ParseInt(arr[2], 10, 64)
	if err != nil || i < 1 {
		http.Error(w, "path must be /set/{integer > 1}", http.StatusBadRequest)
		return
	}

	set := data.GetSet(int(i))

	byt, err := json.Marshal(set)
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
	m.HandleFunc("/set/", handleSet)


	//start server
	//http.ListenAndServe() is a blocking call
	err := http.ListenAndServe(":8080", m)
	panic(err)
}
