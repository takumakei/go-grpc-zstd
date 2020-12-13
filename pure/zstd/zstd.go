// Package zstd implements and registers the pure go zstd compressor
// for gRPC-Go during the initialization.
package zstd

import (
	"io"
	"sync"

	"github.com/google/wuffs/lib/compression"
	"github.com/klauspost/compress/zstd"
	"google.golang.org/grpc/encoding"
)

const Name = "zstd"

func init() {
	c := &compressor{}
	c.poolCompressor.New = func() interface{} {
		return &writer{pool: &c.poolCompressor}
	}
	c.poolDecompressor.New = func() interface{} {
		return &reader{pool: &c.poolDecompressor}
	}
	encoding.RegisterCompressor(c)
}

var (
	compressionLevel = compression.LevelDefault
)

func SetLevel(level compression.Level) {
	compressionLevel = level
}

type writer struct {
	*zstd.Encoder
	pool *sync.Pool
}

func (c *compressor) Compress(w io.Writer) (io.WriteCloser, error) {
	z := c.poolCompressor.Get().(*writer)
	if z.Encoder == nil {
		e, err := zstd.NewWriter(w, zstd.WithEncoderLevel(zstdCompressionLevel(compressionLevel)))
		if err != nil {
			c.poolCompressor.Put(z)
			return nil, err
		}
		z.Encoder = e
	} else {
		z.Encoder.Reset(w)
	}
	return z, nil
}

func (z *writer) Close() error {
	defer z.pool.Put(z)
	return z.Encoder.Close()
}

type reader struct {
	*zstd.Decoder
	pool *sync.Pool
}

func (c *compressor) Decompress(r io.Reader) (io.Reader, error) {
	z := c.poolDecompressor.Get().(*reader)
	d, err := zstd.NewReader(r)
	if err != nil {
		c.poolDecompressor.Put(z)
		return nil, err
	}
	z.Decoder = d
	return z, nil
}

func (z *reader) Read(p []byte) (n int, err error) {
	n, err = z.Decoder.Read(p)
	if err == io.EOF {
		z.Decoder.Close()
		z.Decoder = nil
		z.pool.Put(z)
	}
	return n, err
}

func (c *compressor) Name() string { return Name }

type compressor struct {
	poolCompressor   sync.Pool
	poolDecompressor sync.Pool
}

func zstdCompressionLevel(level compression.Level) zstd.EncoderLevel {
	return zstd.EncoderLevelFromZstd(int(level.Interpolate(1, 2, 3, 15, 22)))
}
