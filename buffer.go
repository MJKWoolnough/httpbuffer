// Package httpbuffer provides a buffer for HTTP requests so that the
// Content-Length may be set.
package httpbuffer

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
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
//
// The compress flag, when true, enables gzip compression for the output.
type Handler struct {
	http.Handler
	Compress bool
}

// ServeHTTP implements the http.Handler interface
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var compressed bool

	if h.Compress {
		accepts := r.Header.Get("Accept-Encoding")
		for {
			p := strings.IndexByte(accepts, ',')
			if p < 0 {
				if strings.TrimSpace(accepts) == "gzip" {
					compressed = true
				}
				break
			}
			if strings.TrimSpace(accepts[:p]) == "gzip" {
				compressed = true
				break
			}
			accepts = accepts[p+1:]
		}
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

	if compressed {
		g := gzipPool.Get().(*gzip.Writer)
		g.Reset(buf)
		resp.Writer = g

		h.Handler.ServeHTTP(rw, r)

		g.Close()

		*resp = responsePusherWriter{}
		gzipPool.Put(g)
		w.Header().Set("Content-Encoding", "gzip")
	} else {
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
