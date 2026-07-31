# httpbuffer

[![CI](https://github.com/MJKWoolnough/httpbuffer/actions/workflows/go-checks.yml/badge.svg)](https://github.com/MJKWoolnough/httpbuffer/actions)
[![Go Reference](https://pkg.go.dev/badge/vimagination.zapto.org/httpbuffer.svg)](https://pkg.go.dev/vimagination.zapto.org/httpbuffer)

--
    import "vimagination.zapto.org/httpbuffer"

Package httpbuffer provides a buffer for HTTP requests so that the `Content-Length` may be set and compression applied for dynamic pages.

## Highlights

 - Buffer HTTP responses before sending them to the client.
 - Automatically sets `Content-Length` header.
 - Supports optional compression which is automatically applied based on `Accept-Encoding` header.

## Usage

```go
package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"

	"vimagination.zapto.org/httpbuffer"
)

func handler(w http.ResponseWriter, r *http.Request) {
	for range 1000 {
		io.WriteString(w, "Hello, World!\n")
	}
}

func main() {
	srv := httptest.NewServer(http.HandlerFunc(handler))
	defer srv.Close()

	resp, err := srv.Client().Get(srv.URL)
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println(resp.ContentLength)
	io.Copy(os.Stdout, io.LimitReader(resp.Body, 14))

	srv = httptest.NewServer(httpbuffer.New(http.HandlerFunc(handler), httpbuffer.Gzip(gzip.DefaultCompression)))
	defer srv.Close()

	req, _ := http.NewRequest(http.MethodGet, srv.URL, nil)
	req.Header.Set("Accept-Encoding", "identity")

	resp, err = srv.Client().Do(req)
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println(resp.ContentLength)
	io.Copy(os.Stdout, io.LimitReader(resp.Body, 14))

	req, _ = http.NewRequest(http.MethodGet, srv.URL, nil)
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err = srv.Client().Do(req)
	if err != nil {
		fmt.Println(err)

		return
	}

	g, _ := gzip.NewReader(resp.Body)

	fmt.Println(resp.ContentLength)
	io.Copy(os.Stdout, io.LimitReader(g, 14))

	// Output:
	// -1
	// Hello, World!
	// 14000
	// Hello, World!
	// 87
	// Hello, World!
}
```

## Documentation

Full API docs can be found at:

https://pkg.go.dev/vimagination.zapto.org/httpbuffer
