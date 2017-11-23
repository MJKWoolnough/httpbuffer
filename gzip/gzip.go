// Package gzip provides an Encoder for the httpbuffer package that uses gzip
// compression.
package gzip

import (
	"compress/gzip"
	"io"
	"sync"

	"github.com/MJKWoolnough/httpbuffer"
)

var (
	// Compression sets the compression level for the gzip encoder
	Compression = gzip.BestCompression

	pool = sync.Pool{
		New: func() interface{} {
			g, _ := gzip.NewWriterLevel(nil, Compression)
			return g
		},
	}
)

type encoding struct{}

func (encoding) Open(w io.Writer) io.Writer {
	g := pool.Get().(*gzip.Writer)
	g.Reset(w)
	return g
}

func (encoding) Close(w io.Writer) {
	g := w.(*gzip.Writer)
	g.Close()
	pool.Put(w)
}

func (encoding) Name() string {
	return "gzip"
}

type encodingX struct {
	encoding
}

func (encodingX) Name() string {
	return "x-gzip"
}

func init() {
	httpbuffer.Register(encoding{})
	httpbuffer.Register(encodingX{})
}
