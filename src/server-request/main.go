package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/checkreq", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("==============================================")
		fmt.Printf("Method: %#v\n", r.Method)
		fmt.Printf("URL: %#v\n", r.URL)
		fmt.Printf("Proto: %#v\n", r.Proto)
		fmt.Printf("Header: %#v\n", r.Header)
		fmt.Printf("Body: %#v\n", r.Body)
		fmt.Printf("ContentLength: %#v\n", r.ContentLength)
		fmt.Printf("Host: %#v\n", r.Host)
		fmt.Printf("Pattern: %#v\n", r.Pattern)
		fmt.Printf("RemoteAddr: %#v\n", r.RemoteAddr)
		fmt.Printf("RequestURI: %#v\n", r.RequestURI)
	})

	server := http.Server{
		Addr:    ":8000",
		Handler: mux,
	}

	server.ListenAndServe()
}
