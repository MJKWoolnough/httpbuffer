# deflate
--
    import "vimagination.zapto.org/httpbuffer/deflate"

Package deflate provides an Encoder for the httpbuffer package that uses deflate compression

Import into `vimagination.zapto.org/httpbuffer` to add deflate compression support.

## Usage

```go
var (
	// Compression sets the compression level for the deflate encoder.
	Compression = flate.BestCompression
)
```
