# zstd
--
    import "vimagination.zapto.org/httpbuffer/zstd"

Package zstd provides an Encoder for the httpbuffer package that uses zstd
compression.

## Usage

```go
var (
	// Compression sets the compression level for the zstd encoder.
	Compression = zstd.SpeedDefault
)
```
