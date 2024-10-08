// Package brotli provides an Encoder for the httpbuffer package that uses
// brotli compression.
package brotli // import "vimagination.zapto.org/httpbuffer/brotli"

import (
	"io"

	"github.com/google/brotli/go/cbrotli"
	"vimagination.zapto.org/httpbuffer"
)

type cbrotliWriter struct {
	*cbrotli.Writer
}

func (c cbrotliWriter) WriteString(str string) (int, error) {
	return c.Write([]byte(str))
}

// Compression sets the compression options for the brotli encoder.
var Compression = cbrotli.WriterOptions{
	Quality: 4,
}

type encoding struct{}

func (encoding) Open(w io.Writer) io.Writer {
	return cbrotliWriter{cbrotli.NewWriter(w, Compression)}
}

func (encoding) Close(w io.Writer) {
	w.(cbrotliWriter).Close()
}

func (encoding) Name() string {
	return "br"
}

func init() {
	httpbuffer.Register(encoding{})
}
