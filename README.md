# httpbuffer
--
    import "vimagination.zapto.org/httpbuffer"

Package httpbuffer provides a buffer for HTTP requests so that the
Content-Length may be set and compression applied for dynamic pages.

## Usage

```go
var (
	// BufferSize determines the initial size of the buffer
	BufferSize = 128 << 10
)
```

#### func  Register

```go
func Register(e Encoding)
```
Register registers the encoding for the buffers to use. Should not be used
passed initialisation.

#### type Encoding

```go
type Encoding interface {
	// Open takes a buffer and returns an encoder-wrapped buffer.
	Open(io.Writer) io.Writer

	// Close returns the encoder-wrapped buffer to flush/close and release
	// resources.
	Close(io.Writer)

	// Name returns the identifier for the encoding algorithm.
	Name() string
}
```

Encoding represents a type that applies a Coding to a byte stream.

#### type Handler

```go
type Handler struct {
	http.Handler
}
```

Handler wraps a http.Handler and provides a buffer and possible gzip
compression. It buffers the Writes and sends the Content-Length header before
Writing the buffer to the client.

#### func (Handler) ServeHTTP

```go
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request)
```
ServeHTTP implements the http.Handler interface.
