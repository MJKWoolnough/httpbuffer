// Package deflate provides an Encoder for the httpbuffer package that uses
// deflate compression
package deflate

import (
	"compress/flate"
	"io"
	"sync"

	"github.com/MJKWoolnough/httpbuffer"
)

var (
	// Compression sets the compression level for the deflate encoder
	Compression = flate.BestCompression

	pool = sync.Pool{
		New: func() interface{} {
			d, _ := flate.NewWriter(nil, Compression)
			return d
		},
	}
)

type encoding struct{}

func (encoding) Open(w io.Writer) io.Writer {
	d := pool.Get().(*flate.Writer)
	d.Reset(w)
	return d
}

func (encoding) Close(w io.Writer) {
	d := w.(*flate.Writer)
	d.Close()
	pool.Put(w)
}

func (encoding) Name() string {
	return "deflate"
}

func init() {
	httpbuffer.Register(encoding{})
}
