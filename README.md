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
