// Package brotli provides an Encoder for the httpbuffer package that uses
// brotli compression.
package brotli

import (
	"io"

	"github.com/MJKWoolnough/httpbuffer"
	"github.com/google/brotli/go/cbrotli"
)

var (
	// Compression sets the compression options for the brotli encoder
	Compression = cbrotli.WriterOptions{
		Quality: 4,
	}
)

type encoding struct{}

func (encoding) Open(w io.Writer) io.Writer {
	return cbrotli.NewWriter(w, Compression)
}

func (encoding) Close(w io.Writer) {
	w.(*cbrotli.Writer).Close()
}

func (encoding) Name() string {
	return "br"
}

func init() {
	httpbuffer.Register(encoding{})
}
