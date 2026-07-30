package httpbuffer

import (
	"compress/flate"
	"compress/gzip"
	"io"
)

type Option func(*Handler)

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

type deflateEncoding int

func (d deflateEncoding) Open(w io.Writer) io.Writer {
	dw, _ := flate.NewWriter(w, int(d))

	return dw
}

func (deflateEncoding) Close(w io.Writer) {
	w.(*flate.Writer).Close()
}

func (deflateEncoding) Name() string {
	return "deflate"
}

func Deflate(compressionLevel int) Option {
	if compressionLevel < flate.HuffmanOnly || compressionLevel > flate.BestCompression {
		compressionLevel = flate.DefaultCompression
	}

	return func(h *Handler) {
		h.encodings = append(h.encodings, "deflate")
		h.compressors["deflate"] = deflateEncoding(compressionLevel)
	}
}
