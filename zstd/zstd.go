// Package zstd provides an Encoder for the httpbuffer package that uses zstd
// compression.
package zstd // import "vimagination.zapto.org/httpbuffer/zstd"

import (
	"io"
	"sync"

	"github.com/klauspost/compress/zstd"
	"vimagination.zapto.org/httpbuffer"
)

type zstdWriter struct {
	*zstd.Encoder
}

func (z zstdWriter) WriteString(str string) (int, error) {
	return z.Write([]byte(str))
}

var (
	// Compression sets the compression level for the zstd encoder.
	Compression = zstd.SpeedDefault

	pool = sync.Pool{
		New: func() interface{} {
			z, _ := zstd.NewWriter(nil, zstd.WithEncoderLevel(Compression))

			return zstdWriter{z}
		},
	}
)

type encoding struct{}

func (encoding) Open(w io.Writer) io.Writer {
	z := pool.Get().(zstdWriter)

	z.Reset(w)

	return z
}

func (encoding) Close(w io.Writer) {
	z := w.(zstdWriter)

	z.Close()
	pool.Put(w)
}

func (encoding) Name() string {
	return "zstd"
}

func init() {
	httpbuffer.Register(encoding{})
}
