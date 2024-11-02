package main

import (
	"fmt"
	"net/http"
)

func middlewareHandleFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("middlewareHandleFunc start: %s\n", r.URL)

		next.ServeHTTP(w, r)

		fmt.Printf("middlewareHandleFunc end: %s\n", r.URL)
	})
}

func basepathHandleFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("basepath response\n"))
}

func main() {
	mux := http.NewServeMux()

	mux.Handle("/basepath", middlewareHandleFunc(http.HandlerFunc(basepathHandleFunc)))

	http.ListenAndServe(":8000", mux)
}
