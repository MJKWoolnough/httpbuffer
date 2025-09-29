# gzip
--
    import "vimagination.zapto.org/httpbuffer/gzip"

Package gzip provides an Encoder for the httpbuffer package that uses gzip compression.

Import into `vimagination.zapto.org/httpbuffer` to add gzip compression support.

## Usage

```go
var (
	// Compression sets the compression level for the gzip encoder.
	Compression = gzip.BestCompression
)
```
