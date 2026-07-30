package httpbuffer

import (
	"compress/gzip"
	"io"
)

type Option func(*Handler)

type gzipWriter struct {
	*gzip.Writer
}

func (g gzipWriter) WriteString(str string) (int, error) {
	return g.Write([]byte(str))
}

type gzipEncoding int

func (g gzipEncoding) Open(w io.Writer) io.Writer {
	gw, _ := gzip.NewWriterLevel(w, int(g))

	return gw
}

func (gzipEncoding) Close(w io.Writer) {
	w.(*gzip.Writer).Close()
}

func (gzipEncoding) Name() string {
	return "gzip"
}

func Gzip(compressionLevel int) Option {
	if compressionLevel < gzip.HuffmanOnly || compressionLevel > gzip.BestCompression {
		compressionLevel = gzip.DefaultCompression
	}

	return func(h *Handler) {
		h.encodings = append(h.encodings, "gzip")
		h.compressors["gzip"] = gzipEncoding(compressionLevel)
	}
}
