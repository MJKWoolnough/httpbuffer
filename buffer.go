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
}

type encodingType struct {
	Encoding
}

func (e *encodingType) Handle(encoding httpencoding.Encoding) (ok bool) {
	if httpencoding.IsWildcard(encoding) && !httpencoding.IsDisallowedInWildcard(encoding, "") {
		e.Encoding = encodings[""]

		return false
	}

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

	h.Handler.ServeHTTP(
		httpwrap.Wrap(w, httpwrap.OverrideWriter(resp), httpwrap.OverrideHeaderWriter(resp), httpwrap.OverrideStringWriter(nil), httpwrap.OverrideFlusher(nil), httpwrap.OverrideHijacker(nil)),
		r,
	)

	encoding.Close(resp.Writer)

	written := resp.Buffer.Len()

	if written == 0 {
	} else if enc := encoding.Name(); enc != "" {
		w.Header().Set("Content-Encoding", enc)
	}

	w.Header().Set("Content-Length", strconv.Itoa(resp.Buffer.Len()))

	if resp.Status > 0 {
		w.WriteHeader(resp.Status)
	} else if written == 0 {
		w.WriteHeader(http.StatusNoContent)
	}

	if written > 0 {
		w.Write(resp.Buffer.Bytes())
	}

	resp.Status = 0
	resp.Writer = nil
	resp.Buffer.Reset()
	responsePool.Put(resp)
}

type responseWriter struct {
	Status int
	io.Writer
	Buffer bytes.Buffer
}

func (r *responseWriter) WriteHeader(s int) {
	r.Status = s
}
