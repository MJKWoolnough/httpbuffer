# brotli
--
    import "vimagination.zapto.org/httpbuffer/brotli"

Package brotli provides an Encoder for the httpbuffer package that uses brotli
compression.

## Usage

```go
var Compression = cbrotli.WriterOptions{
	Quality: 4,
}
```
Compression sets the compression options for the brotli encoder.
