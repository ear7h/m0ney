package main

import "net/http"

func get(w http.ResponseWriter, r *http.Request) {

}

func delete(w http.ResponseWriter, r *http.Request) {

}

func main() {
	http.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("auth")[0] !=

		switch r.Method {
		case http.MethodGet:
			get(w, r)
		case http.MethodDelete:
			delete(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":8080", nil)
}
