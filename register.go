package httpbuffer

import (
	"io"
	"slices"

	"vimagination.zapto.org/httpencoding"
)

var encodings = map[httpencoding.Encoding]Encoding{
	"": identity{},
}

type identity struct{}

func (identity) Open(w io.Writer) io.Writer {
	return w
}

func (identity) Close(io.Writer) {}

func (identity) Name() string {
	return ""
}

// Encoding represents a type that applies a Coding to a byte stream.
type Encoding interface {
	// Open takes a buffer and returns an encoder-wrapped buffer.
	Open(io.Writer) io.Writer

	// Close returns the encoder-wrapped buffer to flush/close and release
	// resources.
	Close(io.Writer)

	// Name returns the identifier for the encoding algorithm.
	Name() string
}

var order = []httpencoding.Encoding{""}

// Register registers the encoding for the buffers to use. Should not be used
// passed initialisation.
func Register(e Encoding) {
	encodings[httpencoding.Encoding(e.Name())] = e
	order = slices.Insert(order, len(order)-1, httpencoding.Encoding(e.Name()))
}
