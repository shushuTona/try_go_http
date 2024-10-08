package main

import (
	"fmt"
	"net/http"
)

type Handle1 struct{}

func (h *Handle1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Handle1.ServeHTTP()\n")
}

func main() {
	server := http.Server{
		Addr: ":8000",
	}

	http.Handle("/", &Handle1{})
	http.Handle("/handle", &Handle1{})
	http.HandleFunc("/handleFunc", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "HandleFunc callback args\n")
	})

	server.ListenAndServe()
}
