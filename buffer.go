// Package httpbuffer provides a buffer for HTTP requests so that the
// Content-Length may be set and compression applied for dynamic pages.
package httpbuffer // import "vimagination.zapto.org/httpbuffer"

import (
	"io"
	"net/http"
	"strconv"
	"sync"

	"vimagination.zapto.org/httpencoding"
	"vimagination.zapto.org/httpwrap"
	"vimagination.zapto.org/memio"
)

var (
	// BufferSize determines the initial size of the buffer
	BufferSize = 128 << 10

	responsePool = sync.Pool{
		New: func() interface{} {
			return &responseWriter{
				Buffer: make(memio.Buffer, 0, BufferSize),
			}
		},
	}
)

// Handler wraps a http.Handler and provides a buffer and possible gzip
// compression. It buffers the Writes and sends the Content-Length header
// before Writing the buffer to the client.
type Handler struct {
	http.Handler
}

type encodingType struct {
	Encoding
}

func (e *encodingType) Handle(encoding httpencoding.Encoding) (ok bool) {
	e.Encoding, ok = encodings[encoding]
	return ok
}

// ServeHTTP implements the http.Handler interface.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var encoding encodingType

	httpencoding.HandleEncoding(r, &encoding)

	if encoding.Encoding == nil {
		httpencoding.InvalidEncoding(w)
		return
	}
	httpencoding.ClearEncoding(r)

	resp := responsePool.Get().(*responseWriter)

	resp.Writer = encoding.Open(&resp.Buffer)
	sw, _ := resp.Writer.(httpwrap.StringWriter)
	h.Handler.ServeHTTP(
		httpwrap.Wrap(w, httpwrap.OverrideWriter(resp), httpwrap.OverrideHeaderWriter(resp), httpwrap.OverrideStringWriter(sw), httpwrap.OverrideFlusher(nil), httpwrap.OverrideHijacker(nil)),
		r,
	)
	encoding.Close(resp.Writer)
	if resp.Written == 0 {
	} else if enc := encoding.Name(); enc != "" {
		w.Header().Set("Content-Encoding", enc)
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(resp.Buffer)))
	if resp.Status > 0 {
		w.WriteHeader(resp.Status)
	}

	if resp.Written > 0 {
		w.Write(resp.Buffer)
	}

	*resp = responseWriter{
		Buffer: resp.Buffer[:0],
	}

	responsePool.Put(resp)
}

type responseWriter struct {
	Status  int
	Writer  io.Writer
	Written int64
	Buffer  memio.Buffer
}

func (r *responseWriter) Write(p []byte) (int, error) {
	n, err := r.Writer.Write(p)
	r.Written += int64(n)
	return n, err
}

func (r *responseWriter) WriteHeader(s int) {
	r.Status = s
}
