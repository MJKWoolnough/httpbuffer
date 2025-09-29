# brotli
--
    import "vimagination.zapto.org/httpbuffer/brotli"

Package brotli provides an Encoder for the httpbuffer package that uses brotli
compression.

Import into `vimagination.zapto.org/httpbuffer` to add brotli compression support.

## Usage

```go
var Compression = cbrotli.WriterOptions{
	Quality: 4,
}
```
Compression sets the compression options for the brotli encoder.
