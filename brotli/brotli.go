// Package brotli provides an Encoder for the httpbuffer package that uses
// brotli compression.
//
// Deprecated: Use httpencoding.New instead, and supply the httpencoding.Brotli
// option.
package brotli // import "vimagination.zapto.org/httpbuffer/brotli"

import (
	"io"
	"sync"

	"github.com/molecule-man/go-brrr"
	"vimagination.zapto.org/httpbuffer"
)

type brotliWriter struct {
	*brrr.Writer
}

func (b brotliWriter) WriteString(str string) (int, error) {
	return b.Write([]byte(str))
}

var (
	// Compression sets the compression level for the brotli encoder.
	Compression = brrr.BestCompression

	pool = sync.Pool{
		New: func() interface{} {
			b, _ := brrr.NewWriter(nil, Compression)

			return brotliWriter{b}
		},
	}
)

type encoding struct{}

func (encoding) Open(w io.Writer) io.Writer {
	b := pool.Get().(brotliWriter)

	b.Reset(w)

	return b
}

func (encoding) Close(w io.Writer) {
	b := w.(brotliWriter)

	b.Close()
	pool.Put(w)
}

func (encoding) Name() string {
	return "br"
}

func init() {
	httpbuffer.Register(encoding{})
}
