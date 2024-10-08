package main

import (
	"fmt"
	"net/http"
)

const SERVER_PORT = ":8000"

func NotWorkingServer() error {
	server := http.Server{
		Addr: SERVER_PORT,
	}

	return server.ListenAndServe()
}

type TestHandler struct{}

func (h *TestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "SimpleServer\n")
}

func SimpleServer() error {
	server := http.Server{
		Addr:    SERVER_PORT,
		Handler: &TestHandler{},
	}

	return server.ListenAndServe()
}

func main() {
	SimpleServer()
}
