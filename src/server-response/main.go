package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/checkresponse", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "checkresponse\n")
	})

	mux.HandleFunc("/serheader", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "{\"id\": 100, \"name\": \"test\"}")
	})

	http.ListenAndServe(":8000", mux)
}
