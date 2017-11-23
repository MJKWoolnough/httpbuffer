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
	bufferPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, BufferSize))
		},
	}

	responsePool = sync.Pool{
		New: func() interface{} {
			return new(responsePusherWriter)
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

	buf := bufferPool.Get().(*bytes.Buffer)
	resp := responsePool.Get().(*responsePusherWriter)
	resp.ResponseWriter = w

	var rw http.ResponseWriter
	if pusher, ok := w.(http.Pusher); ok {
		resp.Pusher = pusher
		rw = resp
	} else {
		rw = &resp.responseWriter
	}

	resp.Writer = encoding.Open(buf)
	h.Handler.ServeHTTP(rw, r)
	encoding.Close(resp.Writer)
	*resp = responsePusherWriter{}
	if enc := encoding.Name(); enc != "" {
		w.Header().Set("Content-Encoding", enc)
	}

	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))

	w.Write(buf.Bytes())

	responsePool.Put(resp)

	buf.Reset()
	bufferPool.Put(buf)
}

type responseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (r *responseWriter) Write(p []byte) (int, error) {
	return r.Writer.Write(p)
}

type responsePusherWriter struct {
	responseWriter
	http.Pusher
}
