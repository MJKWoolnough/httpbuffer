# httpbuffer
--
    import "github.com/MJKWoolnough/httpbuffer"

Package httpbuffer provides a buffer for HTTP requests so that the
Content-Length may be set.

## Usage

#### type Handler

```go
type Handler struct {
	http.Handler
	Compress bool
}
```

Handler wraps a http.Handler and provides a buffer and possible gzip
compression. It buffers the Writes and sends the Content-Length header before
Writing the buffer to the client.

The compress flag, when true, enables gzip compression for the output.

#### func (Handler) ServeHTTP

```go
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request)
```
ServeHTTP implements the http.Handler interface
