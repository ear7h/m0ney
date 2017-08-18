package main

import "net/http"

func main() {
	m := http.NewServeMux()

	//list available data sets
	m.HandleFunc("/list", handleList)

	//session
	//get (with auth) - returns values
	//post - returns key and value for session
	m.HandleFunc("/session", handleSession)



}
