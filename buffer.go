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
				responseWriter: responseWriter{
					Buffer: bytes.NewBuffer(make([]byte, 0, BufferSize)),
				},
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
	resp.Status = 200

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
	resp.Pusher = nil
	resp.Writer = nil
	resp.ResponseWriter = nil
	if enc := encoding.Name(); enc != "" {
		w.Header().Set("Content-Encoding", enc)
	}

	w.Header().Set("Content-Length", strconv.Itoa(resp.Buffer.Len()))
	if resp.Status != 200 {
		w.WriteHeader(resp.Status)
	}

	w.Write(resp.Buffer.Bytes())

	resp.Buffer.Reset()
	responsePool.Put(resp)
}

type responseWriter struct {
	http.ResponseWriter
	Status int
	Writer io.Writer
	Buffer *bytes.Buffer
}

func (r *responseWriter) Write(p []byte) (int, error) {
	return r.Writer.Write(p)
}

func (r *responseWriter) WriteHeader(s int) {
	r.Status = s
}

type responsePusherWriter struct {
	responseWriter
	http.Pusher
}
