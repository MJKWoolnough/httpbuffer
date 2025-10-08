# httpbuffer

[![CI](https://github.com/MJKWoolnough/httpbuffer/actions/workflows/go-checks.yml/badge.svg)](https://github.com/MJKWoolnough/httpbuffer/actions)
[![Go Reference](https://pkg.go.dev/badge/vimagination.zapto.org/httpbuffer.svg)](https://pkg.go.dev/vimagination.zapto.org/httpbuffer)
[![Go Report Card](https://goreportcard.com/badge/vimagination.zapto.org/httpbuffer)](https://goreportcard.com/report/vimagination.zapto.org/httpbuffer)

--
    import "vimagination.zapto.org/httpbuffer"

Package httpbuffer provides a buffer for HTTP requests so that the `Content-Length` may be set and compression applied for dynamic pages.

## Highlights

 - Buffer HTTP responses before sending them to the client.
 - Automatically sets `Content-Length` header.
 - Supports optional compression which is automatically applied based on `Accept-Encoding` header.
 - Import `vimagination.zapto.org/httpbuffer/{brotli,deflate,gzip}` to support compression.

## Usage

```go
package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"vimagination.zapto.org/httpbuffer"
	_ "vimagination.zapto.org/httpbuffer/gzip"
)

func handler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, World!")
}

func main() {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	handler(w, r)

	fmt.Println(w.Result().ContentLength)

	w = httptest.NewRecorder()
	buf := httpbuffer.Handler{Handler: http.HandlerFunc(handler)}
	buf.ServeHTTP(w, r)

	fmt.Println(w.Result().ContentLength)
	fmt.Println(w.Body)

	// Output:
	// -1
	// 13
	// Hello, World!
}
```

## Documentation

Full API docs can be found at:

https://pkg.go.dev/vimagination.zapto.org/httpbuffer
