package httpbuffer_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"vimagination.zapto.org/httpbuffer"
	_ "vimagination.zapto.org/httpbuffer/gzip"
)

func TestBuffer(t *testing.T) {
	for n, test := range [...]struct {
		Buffers
		code     int
		compress bool
		output   string
		length   int
	}{
		{
			Buffers: Buffers{},
			code:    http.StatusNoContent,
			output:  "",
			length:  0,
		},
		{
			Buffers: Buffers{[]byte("data")},
			code:    http.StatusOK,
			output:  "data",
			length:  4,
		},
		{
			Buffers: Buffers{[]byte("hello, "), []byte("world")},
			code:    http.StatusOK,
			output:  "hello, world",
			length:  12,
		},
		{
			Buffers:  Buffers{[]byte("hello, "), []byte("world")},
			code:     http.StatusOK,
			compress: true,
			output:   "hello, world",
			length:   12,
		},
	} {

		server := httptest.NewServer(httpbuffer.Handler{
			Handler: test.Buffers,
		})

		var buf strings.Builder

		r, _ := http.NewRequest(http.MethodGet, server.URL, nil)

		if !test.compress {
			r.Header.Set("Accept-Encoding", "identity")
		}

		if result, err := server.Client().Do(r); err != nil {
			t.Errorf("test %d: unexpected error: %v", n+1, err)
		} else if result.StatusCode != test.code {
			t.Errorf("test %d: expecting code %d, got %d", n+1, test.code, result.StatusCode)
		} else if result.Uncompressed != test.compress {
			t.Errorf("test %d: unexpected Uncompressed to be %v, got %v", n+1, test.compress, result.Uncompressed)
		} else if _, err = io.Copy(&buf, result.Body); err != nil {
			t.Errorf("test %d: unexpected error copying body: %v", n+1, err)
		} else if output := buf.String(); output != test.output {
			t.Errorf("test %d: expecting output %q, got %q", n+1, test.output, output)
		} else if len(output) != test.length {
			t.Errorf("test %d: expecting content length %d, got %d", n+1, test.length, len(output))
		}

		server.Close()
	}
}

type Buffers [][]byte

func (b Buffers) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	for _, p := range b {
		w.Write(p)
	}
}
