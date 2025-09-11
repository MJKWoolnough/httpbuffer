// Package deflate provides an Encoder for the httpbuffer package that uses
// deflate compression
package deflate // import "vimagination.zapto.org/httpbuffer/deflate"

import (
	"compress/flate"
	"io"
	"sync"

	"vimagination.zapto.org/httpbuffer"
)

type flateWriter struct {
	*flate.Writer
}

func (f flateWriter) WriteString(str string) (int, error) {
	return f.Write([]byte(str))
}

var (
	// Compression sets the compression level for the deflate encoder.
	Compression = flate.BestCompression

	pool = sync.Pool{
		New: func() interface{} {
			d, _ := flate.NewWriter(nil, Compression)
			return flateWriter{d}
		},
	}
)

type encoding struct{}

func (encoding) Open(w io.Writer) io.Writer {
	d := pool.Get().(flateWriter)
	d.Reset(w)

	return d
}

func (encoding) Close(w io.Writer) {
	d := w.(flateWriter)

	d.Close()
	pool.Put(w)
}

func (encoding) Name() string {
	return "deflate"
}

func init() {
	httpbuffer.Register(encoding{})
}
