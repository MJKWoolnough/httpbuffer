// Package httpbuffer provides a buffer for HTTP requests so that the
// Content-Length may be set and compression applied for dynamic pages.
package httpbuffer

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strconv"
	"sync"

	"github.com/MJKWoolnough/httpencoding"
)

var (
	bufferPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 128<<10))
		},
	}

	responsePool = sync.Pool{
		New: func() interface{} {
			return new(responsePusherWriter)
		},
	}
	gzipPool = sync.Pool{
		New: func() interface{} {
			g, _ := gzip.NewWriterLevel(nil, gzip.BestCompression)
			return g
		},
	}
)

// Handler wraps a http.Handler and provides a buffer and possible gzip
// compression. It buffers the Writes and sends the Content-Length header
// before Writing the buffer to the client.
type Handler struct {
	http.Handler
}

const (
	encodingIdentity = iota + 1
	encodingGzip
)

type encodingType uint8

func (e *encodingType) Handle(encoding string) bool {
	switch encoding {
	case "":
		*e = encodingIdentity
	case "gzip", "x-gzip":
		*e = encodingGzip
	default:
		return false
	}
	return true
}

// ServeHTTP implements the http.Handler interface
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var encoding encodingType

	httpencoding.HandleEncoding(r, &encoding)

	if encoding == 0 {
		httpencoding.InvalidEncoding(w)
		return
	}

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

	switch encoding {
	case encodingGzip:
		g := gzipPool.Get().(*gzip.Writer)
		g.Reset(buf)
		resp.Writer = g

		h.Handler.ServeHTTP(rw, r)

		g.Close()

		*resp = responsePusherWriter{}
		gzipPool.Put(g)
		w.Header().Set("Content-Encoding", "gzip")
	default:
		resp.Writer = buf
		h.Handler.ServeHTTP(rw, r)
		*resp = responsePusherWriter{}
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
