// Package httpbuffer provides a buffer for HTTP requests so that the
// Content-Length may be set and compression applied for dynamic pages.
package httpbuffer // import "vimagination.zapto.org/httpbuffer"

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"sync"

	"vimagination.zapto.org/httpencoding"
	"vimagination.zapto.org/httpwrap"
)

var (
	// BufferSize determines the initial size of the buffer
	BufferSize = 128 << 10

	responsePool = sync.Pool{
		New: func() any {
			return new(responseWriter)
		},
	}
)

// Handler wraps a http.Handler and provides a buffer and possible gzip
// compression. It buffers the Writes and sends the Content-Length header
// before Writing the buffer to the client.
type Handler struct {
	http.Handler
	encodings   []httpencoding.Encoding
	compressors map[httpencoding.Encoding]Encoding
}

func New(handler http.Handler, opts ...Option) *Handler {
	h := &Handler{
		Handler:     handler,
		compressors: make(map[httpencoding.Encoding]Encoding),
	}

	for _, opt := range opts {
		opt(h)
	}

	h.encodings = append(h.encodings, "")
	h.compressors[""] = identity{}

	return h
}

// ServeHTTP implements the http.Handler interface.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	order := order
	encodings := encodings

	if h.encodings != nil {
		order = h.encodings
		encodings = h.compressors
	}

	enc, ok := httpencoding.Negotiate(r, order...)
	if !ok {
		httpencoding.InvalidEncoding(w)

		return
	}

	httpencoding.ClearEncoding(r)

	encoding := encodings[enc]
	resp := responsePool.Get().(*responseWriter)
	resp.Writer = encoding.Open(&resp.Buffer)

	h.Handler.ServeHTTP(
		httpwrap.Wrap(w, httpwrap.OverrideWriter(resp), httpwrap.OverrideHeaderWriter(resp), httpwrap.OverrideStringWriter(nil), httpwrap.OverrideFlusher(nil), httpwrap.OverrideHijacker(nil)),
		r,
	)

	encoding.Close(resp.Writer)

	if resp.written {
		if enc != "" {
			w.Header().Set("Content-Encoding", string(enc))
		}

		w.Header().Set("Content-Length", strconv.Itoa(resp.Buffer.Len()))
	}

	if resp.Status > 0 {
		w.WriteHeader(resp.Status)
	} else if !resp.written {
		w.WriteHeader(http.StatusNoContent)
	}

	if resp.written {
		w.Write(resp.Buffer.Bytes())
	}

	resp.Status = 0
	resp.Writer = nil
	resp.written = false
	resp.Buffer.Reset()
	responsePool.Put(resp)
}

type responseWriter struct {
	written bool
	Status  int
	io.Writer
	Buffer bytes.Buffer
}

func (r *responseWriter) WriteHeader(s int) {
	r.Status = s
}

func (r *responseWriter) Write(p []byte) (int, error) {
	r.written = r.written || len(p) > 0

	return r.Writer.Write(p)
}
