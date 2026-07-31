package httpbuffer

import (
	"compress/flate"
	"compress/gzip"
	"io"

	"github.com/klauspost/compress/zstd"
	"github.com/molecule-man/go-brrr"
)

// Option represents a type that can be passed to the New function to specify
// which compressors can be used to compress the buffered data.
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

// Gzip enables the use of gzip to compress the data that will be sent to the
// client when gzip is allowed in the Accept-Encoding header.
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

// Deflate enables the use of flate to compress the data that will be sent to
// the client when deflate is allowed in the Accept-Encoding header.
func Deflate(compressionLevel int) Option {
	if compressionLevel < flate.HuffmanOnly || compressionLevel > flate.BestCompression {
		compressionLevel = flate.DefaultCompression
	}

	return func(h *Handler) {
		h.encodings = append(h.encodings, "deflate")
		h.compressors["deflate"] = deflateEncoding(compressionLevel)
	}
}

type brotliEncoding int

func (b brotliEncoding) Open(w io.Writer) io.Writer {
	dw, _ := brrr.NewWriter(w, int(b))

	return dw
}

func (brotliEncoding) Close(w io.Writer) {
	w.(*brrr.Writer).Close()
}

func (brotliEncoding) Name() string {
	return "br"
}

// Brotli enables the use of brotli to compress the data that will be sent to
// the client when br is allowed in the Accept-Encoding header.
func Brotli(compressionLevel int) Option {
	if compressionLevel < 0 {
		compressionLevel = brrr.BestSpeed
	} else if compressionLevel > 11 {
		compressionLevel = brrr.BestCompression
	}

	return func(h *Handler) {
		h.encodings = append(h.encodings, "br")
		h.compressors["br"] = brotliEncoding(compressionLevel)
	}
}

type zstdEncoding int

func (z zstdEncoding) Open(w io.Writer) io.Writer {
	dw, _ := zstd.NewWriter(w, zstd.WithEncoderLevel(zstd.EncoderLevel(z)))

	return dw
}

func (zstdEncoding) Close(w io.Writer) {
	w.(*zstd.Encoder).Close()
}

func (zstdEncoding) Name() string {
	return "zstd"
}

// Zstd enables the use of zstandard to compress the data that will be sent to
// the client when zstd is allowed in the Accept-Encoding header.
func Zstd(compressionLevel zstd.EncoderLevel) Option {
	if compressionLevel < zstd.SpeedFastest {
		compressionLevel = zstd.SpeedFastest
	} else if compressionLevel > zstd.SpeedBestCompression {
		compressionLevel = zstd.SpeedBestCompression
	}

	return func(h *Handler) {
		h.encodings = append(h.encodings, "zstd")
		h.compressors["zstd"] = zstdEncoding(compressionLevel)
	}
}
