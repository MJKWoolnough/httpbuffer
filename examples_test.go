package httpbuffer_test

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

func Example() {
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
