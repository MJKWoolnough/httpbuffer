package httpbuffer

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
			Buffers:  Buffers{},
			code:     http.StatusNoContent,
			compress: true,
			output:   "",
			length:   0,
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
		var buf strings.Builder

		server := httptest.NewServer(New(test.Buffers, Gzip(gzip.DefaultCompression)))
		r, _ := http.NewRequest(http.MethodGet, server.URL, nil)

		if !test.compress {
			r.Header.Set("Accept-Encoding", "identity")
		}

		if result, err := server.Client().Do(r); err != nil {
			t.Errorf("test %d: unexpected error: %v", n+1, err)
		} else if result.StatusCode != test.code {
			t.Errorf("test %d: expecting code %d, got %d", n+1, test.code, result.StatusCode)
		} else if result.Uncompressed != test.compress && test.length > 0 {
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

func TestInvalidEncoding(t *testing.T) {
	server := httptest.NewServer(New(Buffers{}, Gzip(gzip.DefaultCompression)))
	r, _ := http.NewRequest(http.MethodGet, server.URL, nil)
	r.Header.Set("Accept-Encoding", "identity;q=0")

	if result, err := server.Client().Do(r); err != nil {
		t.Errorf("unexpected error: %s", err)
	} else if result.StatusCode != http.StatusNotAcceptable {
		t.Errorf("expecting response code %d, got %d", http.StatusNotAcceptable, result.StatusCode)
	}
}

func TestCustomStatusCode(t *testing.T) {
	server := httptest.NewServer(New(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusConflict)
		io.WriteString(w, "bad")
	}), Gzip(gzip.DefaultCompression)))
	r, _ := http.NewRequest(http.MethodGet, server.URL, nil)

	if result, err := server.Client().Do(r); err != nil {
		t.Errorf("unexpected error: %s", err)
	} else if result.StatusCode != http.StatusConflict {
		t.Errorf("expecting response code %d, got %d", http.StatusConflict, result.StatusCode)
	} else if data, err := io.ReadAll(result.Body); err != nil {
		t.Errorf("unexpected error: %s", err)
	} else if string(data) != "bad" {
		t.Errorf("expecting output %q, got %q", "bad", data)
	}
}

type Buffers [][]byte

func (b Buffers) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	for _, p := range b {
		w.Write(p)
	}
}
