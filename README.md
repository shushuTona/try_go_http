# try_go_http

[net/http](https://pkg.go.dev/net/http) パッケージを使用したWebサーバーの実装。

`http.Server` がサーバー自体の構造体で、 `Handler` フィールドに `http.Handler` インターフェースを満たす構造体（＝ `ServeHTTP` を実装している構造体）を定義する。

### TODO

下記の流れ確認する。

1. `func (srv *http.Server) ListenAndServe() error` 
1. `net.Listen("tcp", addr)` で `net.Listener` 生成
1. `func (srv *Server) Serve(l net.Listener) error` でサーバー起動
1. `func (c *conn) serve(ctx context.Context)` 

## ハンドラを設定する

### ServeMux ・ DefaultServeMux

- `http.Server` の `Handler` フィールドが `nil` の場合、 `http.DefaultServeMux` （= `http.ServeMux` :  `mutex` を内包していて複数のパスを登録できるようになっている）が使用される。
    - （`serverHandler.ServeHTTP` で `Handler` が `nil` の場合 `http.DefaultServeMux` を使用してサーバーを起動している）

- `http.ServeMux` は `ServeHTTP(w ResponseWriter, r *Request)` が実装されている（＝ `http.Handler` インターフェースを満たす）ため、下記のように `http.Server` の `Handler` に指定して使用する。
    ```go
        mux := http.NewServeMux()

        mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprintf(w, "test\n")
        })

        server := http.Server{
            Addr:    ":8000",
            Handler: mux,
        }

        server.ListenAndServe()
    ```

### パスの登録方法

- `http.Handle` を使用することで `http.DefaultServeMux` にパスに対してハンドラを設定することができる。（ `*ServeMux.register` で `http.DefaultServeMux` の `tree`, `index`, `patterns` に指定のパスを登録する）

- `http.HandleFunc` は関数の引数に直接ハンドラを指定することができる。
（`http.DefaultServeMux` への登録方法は `http.Handle` と同じ）

    - → `http.Handle` の方がハンドラ指定が `http.Handler` インターフェースで抽象化されているから、テストなどがしやすそう。

--- 

## Request

```go
	mux.HandleFunc("/checkreq", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}

		fmt.Printf("Method: %#v\n", r.Method)
		fmt.Printf("URL: %#v\n", r.URL)
		fmt.Printf("Proto: %#v\n", r.Proto)
		fmt.Printf("Header: %#v\n", r.Header)
		fmt.Printf("Body: %#v\n", string(body))
		fmt.Printf("ContentLength: %#v\n", r.ContentLength)
		fmt.Printf("Host: %#v\n", r.Host)
		fmt.Printf("Pattern: %#v\n", r.Pattern)
		fmt.Printf("RemoteAddr: %#v\n", r.RemoteAddr)
		fmt.Printf("RequestURI: %#v\n", r.RequestURI)
	})
```

### GET

```bash
curl http://localhost:8000/checkreq
```

```
Method: "GET"
URL: &url.URL{Scheme:"", Opaque:"", User:(*url.Userinfo)(nil), Host:"", Path:"/checkreq", RawPath:"", OmitHost:false, ForceQuery:false, RawQuery:"", Fragment:"", RawFragment:""}
Proto: "HTTP/1.1"
Header: http.Header{"Accept":[]string{"*/*"}, "User-Agent":[]string{"curl/7.88.1"}}
Body: ""
ContentLength: 0
Host: "localhost:8000"
Pattern: "/checkreq"
RemoteAddr: "127.0.0.1:49124"
RequestURI: "/checkreq"
```

### POST

```bash
curl -X POST -H "Content-Type: application/json" -d '{"Name":"tanaka", "Age":"20"}' localhost:8000/checkreq
```

```
Method: "POST"
URL: &url.URL{Scheme:"", Opaque:"", User:(*url.Userinfo)(nil), Host:"", Path:"/checkreq", RawPath:"", OmitHost:false, ForceQuery:false, RawQuery:"", Fragment:"", RawFragment:""}
Proto: "HTTP/1.1"
Header: http.Header{"Accept":[]string{"*/*"}, "Content-Length":[]string{"29"}, "Content-Type":[]string{"application/json"}, "User-Agent":[]string{"curl/7.88.1"}}
Body: "{\"Name\":\"tanaka\", \"Age\":\"20\"}"
ContentLength: 29
Host: "localhost:8000"
Pattern: "/checkreq"
RemoteAddr: "127.0.0.1:36822"
RequestURI: "/checkreq"
```

### path parameter

Go 1.22で [`http.Request.PathValue`](https://pkg.go.dev/net/http#Request.PathValue) が追加されたことでWebフレームワークやルーティングライブラリを使用しなくてもパスパラメータを簡単に取得できるようになった。

```go
	mux.HandleFunc("/page/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("r.PathValue: %s\n", r.PathValue("id"))
		fmt.Printf("Method: %#v\n", r.Method)
		fmt.Printf("URL: %#v\n", r.URL)
		fmt.Printf("Proto: %#v\n", r.Proto)
		fmt.Printf("Header: %#v\n", r.Header)
		fmt.Printf("ContentLength: %#v\n", r.ContentLength)
		fmt.Printf("Host: %#v\n", r.Host)
		fmt.Printf("Pattern: %#v\n", r.Pattern)
		fmt.Printf("RemoteAddr: %#v\n", r.RemoteAddr)
		fmt.Printf("RequestURI: %#v\n", r.RequestURI)
	})
```

```bash
curl http://localhost:8000/page/100
```

```
r.PathValue: 100
Method: "GET"
URL: &url.URL{Scheme:"", Opaque:"", User:(*url.Userinfo)(nil), Host:"", Path:"/page/100", RawPath:"", OmitHost:false, ForceQuery:false, RawQuery:"", Fragment:"", RawFragment:""}
Proto: "HTTP/1.1"
Header: http.Header{"Accept":[]string{"*/*"}, "User-Agent":[]string{"curl/7.88.1"}}
ContentLength: 0
Host: "localhost:8000"
Pattern: "/page/{id}"
RemoteAddr: "127.0.0.1:51772"
RequestURI: "/page/100"
```

### path query

パスに指定されているクエリは `http.Request` の `URL` フィールドを使用して下記の流れで取得する。

```go
	mux.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("r.URL.Query(): %#v\n", r.URL.Query())
		fmt.Printf("r.URL.Query() query1: %#v\n", r.URL.Query().Get("query1"))
		fmt.Printf("r.URL.Query() query2: %#v\n", r.URL.Query().Get("query2"))
		PrintRequest(r)
	})
```

```bash
curl "http://localhost:8000/query?query1=1&query2=2"
```

```
r.URL.Query(): url.Values{"query1":[]string{"1"}, "query2":[]string{"2"}}
r.URL.Query() query1: "1"
r.URL.Query() query2: "2"
Method: "GET"
URL: &url.URL{Scheme:"", Opaque:"", User:(*url.Userinfo)(nil), Host:"", Path:"/query", RawPath:"", OmitHost:false, ForceQuery:false, RawQuery:"query1=1&query2=2", Fragment:"", RawFragment:""}
Proto: "HTTP/1.1"
Header: http.Header{"Accept":[]string{"*/*"}, "User-Agent":[]string{"curl/7.88.1"}}
ContentLength: 0
Host: "localhost:8000"
Pattern: "/query"
RemoteAddr: "127.0.0.1:59506"
RequestURI: "/query?query1=1&query2=2"
```

## Response

ハンドラ内で使用することができる `http.ResponseWriter` は `io.Writer` を満たしているので、下記のように一般的な書き込み処理を使用することでレスポンス内容を生成することができる。

```go
	mux.HandleFunc("/checkresponse", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "checkresponse\n")
	})
```

```bash
curl -D - -X POST -H "Content-Type: application/json" -d '{"Name":"tanaka", "Age":"20"}' localhost:8000/checkresponse
```

```
HTTP/1.1 200 OK
Date: Sat, 02 Nov 2024 10:00:18 GMT
Content-Length: 14
Content-Type: text/plain; charset=utf-8

checkresponse
```

### Header

`w.Header()` で取得した `http.Header` の `Set` メソッドでレスポンスヘッダーを設定する。

```go
	mux.HandleFunc("/setheader", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Custom-Header", "test")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "{\"id\": 100, \"name\": \"test\"}")
	})
```

```bash
curl -D - -X POST -H "Content-Type: application/json" -d '{"Name":"tanaka", "Age":"20"}' localhost:8000/setheader
```

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Custom-Header: test
Date: Sat, 02 Nov 2024 10:00:36 GMT
Content-Length: 28

{"id": 100, "name": "test"}
```

## cookie

### Request

`r.Cookie` メソッドでcookie名を指定して `*http.Cookie` を取得する。

```go
	mux.HandleFunc("/checkcookie", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("cookie")
		if err != nil {
			http.Error(w, fmt.Sprintf("error: not found cookie"), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "cookie: %#v\n", cookie)
	})
```

```bash
curl -X Get --cookie 'cookie="XXXXX"' -D - localhost:8000/checkcookie
```

```bash
HTTP/1.1 200 OK
Date: Sat, 02 Nov 2024 14:12:56 GMT
Content-Length: 263
Content-Type: text/plain; charset=utf-8

cookie: &http.Cookie{Name:"cookie", Value:"XXXXX", Quoted:true, Path:"", Domain:"", Expires:time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), RawExpires:"", MaxAge:0, Secure:false, HttpOnly:false, SameSite:0, Partitioned:false, Raw:"", Unparsed:[]string(nil)}
```

`r.Cookies` メソッドは `[]*http.Cookie` を取得することができる。

```go
	mux.HandleFunc("/checkcookies", func(w http.ResponseWriter, r *http.Request) {
		cookies := r.Cookies()
		response := ""
		for _, cookie := range cookies {
			response = fmt.Sprintf("%s%s: %s\n", response, cookie.Name, cookie.Value)
		}

		fmt.Fprintf(w, response)
	})
```

```bash
curl -X Get --cookie 'cookie="XXXXX";cookie2=YYYYY' -D - localhost:8000/checkcookies
```

```
HTTP/1.1 200 OK
Date: Sat, 02 Nov 2024 14:22:06 GMT
Content-Length: 29
Content-Type: text/plain; charset=utf-8

cookie: XXXXX
cookie2: YYYYY
```

### Response

`http.ParseCookie` メソッドでcookieの文字列から `[]*http.Cookie` を生成して `http.SetCookie` メソッドでレスポンスにcookieを設定する。

```go
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
```

```bash
curl -X Get -D - localhost:8000/setcookie
```

```
HTTP/1.1 200 OK
Set-Cookie: cookie=XXXXX
Set-Cookie: cookie2=YYYYY
Date: Sat, 02 Nov 2024 14:51:55 GMT
Content-Length: 9
Content-Type: text/plain; charset=utf-8
```
