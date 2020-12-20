// Package zstd implements and registers the zstd compressor
// for gRPC-Go during the initialization.
package zstd

import (
	"io"
	"sync"

	"github.com/google/wuffs/lib/cgozstd"
	"github.com/google/wuffs/lib/compression"
	"google.golang.org/grpc/encoding"
)

const (
	Name      = "zstd"
	PoweredBy = "github.com/google/wuffs/lib/cgozstd"
)

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
	cgozstd.Writer
	pool *sync.Pool
}

func (c *compressor) Compress(w io.Writer) (io.WriteCloser, error) {
	z := c.poolCompressor.Get().(*writer)
	err := z.Writer.Reset(w, nil, compressionLevel)
	return z, err
}

func (z *writer) Close() error {
	defer z.pool.Put(z)
	return z.Writer.Close()
}

type reader struct {
	cgozstd.Reader
	pool *sync.Pool
}

func (c *compressor) Decompress(r io.Reader) (io.Reader, error) {
	z := c.poolDecompressor.Get().(*reader)
	if err := z.Reset(r, nil); err != nil {
		c.poolDecompressor.Put(z)
		return nil, err
	}
	return z, nil
}

func (z *reader) Read(p []byte) (n int, err error) {
	n, err = z.Reader.Read(p)
	if err == io.EOF {
		z.pool.Put(z)
	}
	return n, err
}

func (c *compressor) Name() string { return Name }

type compressor struct {
	poolCompressor   sync.Pool
	poolDecompressor sync.Pool
}
