package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/checkcookie", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("cookie")
		if err != nil {
			http.Error(w, fmt.Sprintf("error: not found cookie"), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "cookie: %#v\n", cookie)
	})

	mux.HandleFunc("/checkcookies", func(w http.ResponseWriter, r *http.Request) {
		cookies := r.Cookies()
		response := ""
		for _, cookie := range cookies {
			response = fmt.Sprintf("%s%s: %s\n", response, cookie.Name, cookie.Value)
		}

		fmt.Fprintf(w, response)
	})

	mux.HandleFunc("/setcookie", func(w http.ResponseWriter, r *http.Request) {
		cookies, err := http.ParseCookie("cookie=XXXXX; cookie2=YYYYY")
		if err != nil {
			http.Error(w, "cookie error", http.StatusInternalServerError)
			return
		}

		for _, cookie := range cookies {
			http.SetCookie(w, cookie)
		}

		fmt.Fprint(w, "setcookie")
	})

	http.ListenAndServe(":8000", mux)
}
