package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/checkreq", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		fmt.Printf("Body: %#v\n", string(body))

		PrintRequest(r)
	})

	mux.HandleFunc("/page/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("r.PathValue: %s\n", r.PathValue("id"))
		PrintRequest(r)
	})

	mux.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("r.URL.Query(): %#v\n", r.URL.Query())
		fmt.Printf("r.URL.Query() query1: %#v\n", r.URL.Query().Get("query1"))
		fmt.Printf("r.URL.Query() query2: %#v\n", r.URL.Query().Get("query2"))
		PrintRequest(r)
	})

	server := http.Server{
		Addr:    ":8000",
		Handler: mux,
	}

	server.ListenAndServe()
}

func PrintRequest(r *http.Request) {
	fmt.Printf("Method: %#v\n", r.Method)
	fmt.Printf("URL: %#v\n", r.URL)
	fmt.Printf("Proto: %#v\n", r.Proto)
	fmt.Printf("Header: %#v\n", r.Header)
	fmt.Printf("ContentLength: %#v\n", r.ContentLength)
	fmt.Printf("Host: %#v\n", r.Host)
	fmt.Printf("Pattern: %#v\n", r.Pattern)
	fmt.Printf("RemoteAddr: %#v\n", r.RemoteAddr)
	fmt.Printf("RequestURI: %#v\n", r.RequestURI)
}
