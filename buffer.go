// Package httpbuffer provides a buffer for HTTP requests so that the
// Content-Length may be set and compression applied for dynamic pages.
package httpbuffer

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"sync"

	"github.com/MJKWoolnough/httpencoding"
)

var (
	// BufferSize determines the initial size of the buffer
	BufferSize = 128 << 10

	responsePool = sync.Pool{
		New: func() interface{} {
			return &responsePusherWriter{
				Buffer: bytes.NewBuffer(make([]byte, 0, BufferSize)),
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

func (e *encodingType) Handle(encoding string) (ok bool) {
	e.Encoding, ok = encodings[encoding]
	return ok
}

// ServeHTTP implements the http.Handler interface
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var encoding encodingType

	httpencoding.HandleEncoding(r, &encoding)

	if encoding.Encoding == nil {
		httpencoding.InvalidEncoding(w)
		return
	}
	httpencoding.ClearEncoding(r)

	resp := responsePool.Get().(*responsePusherWriter)
	resp.ResponseWriter = w

	var rw http.ResponseWriter
	if pusher, ok := w.(http.Pusher); ok {
		resp.Pusher = pusher
		rw = resp
	} else {
		rw = &resp.responseWriter
	}

	resp.Writer = encoding.Open(resp.Buffer)
	h.Handler.ServeHTTP(rw, r)
	encoding.Close(resp.Writer)
	if resp.Written == 0 {
		resp.Buffer.Reset()
	} else if enc := encoding.Name(); enc != "" {
		w.Header().Set("Content-Encoding", enc)
	}

	w.Header().Set("Content-Length", strconv.Itoa(resp.Buffer.Len()))
	if resp.Status > 0 {
		w.WriteHeader(resp.Status)
	}

	if resp.Written > 0 {
		w.Write(resp.Buffer.Bytes())
		resp.Buffer.Reset()
	}

	*resp = responsePusherWriter{
		Buffer: resp.Buffer,
	}

	responsePool.Put(resp)
}

type responseWriter struct {
	http.ResponseWriter
	Status  int
	Writer  io.Writer
	Written int64
}

func (r *responseWriter) Write(p []byte) (int, error) {
	n, err := r.Writer.Write(p)
	r.Written += int64(n)
	return n, err
}

func (r *responseWriter) WriteHeader(s int) {
	r.Status = s
}

type responsePusherWriter struct {
	responseWriter
	http.Pusher
	Buffer *bytes.Buffer
}
