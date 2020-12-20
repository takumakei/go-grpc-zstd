package bench

import (
	"bytes"
	"crypto/rand"
	"io/ioutil"
	"testing"

	"github.com/google/wuffs/lib/compression"
	"github.com/takumakei/go-grpc-zstd/zstd"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/encoding/gzip"
)

var BlockSize = 1024 * 1024

// SetLevel update the compression level.
//
// NOTE: this function must only be called during initialization time (i.e. in an init() function),
// and is not thread-safe.
func SetLevel(level compression.Level) {
	zstd.SetLevel(level)
	gzip.SetLevel(GzipLevel(level))
}

func BenchmarkCompress(b *testing.B) {
	b.Run("zero", func(b *testing.B) {
		src := make([]byte, BlockSize)
		Zero(src)
		Compress(b, src)
	})

	b.Run("rand", func(b *testing.B) {
		src := make([]byte, BlockSize)
		Rand(src)
		Compress(b, src)
	})
}

func BenchmarkDecompress(b *testing.B) {
	b.Run("zero", func(b *testing.B) {
		src := make([]byte, BlockSize)
		Zero(src)
		Decompress(b, src)
	})

	b.Run("rand", func(b *testing.B) {
		src := make([]byte, BlockSize)
		Rand(src)
		Decompress(b, src)
	})
}

func GzipLevel(level compression.Level) int {
	return int(level.Interpolate(1, 2, 6, 9, 9))
}

func Zero(p []byte) {
	for i := range p {
		p[i] = 0
	}
}

func Rand(p []byte) {
	rand.Read(p)
}

func Compress(b *testing.B, src []byte) {
	list := []string{gzip.Name, zstd.Name}

	for _, name := range list {
		b.Run(name, func(b *testing.B) {
			c := encoding.GetCompressor(name)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				w, err := c.Compress(ioutil.Discard)
				if err != nil {
					b.Error(err)
				}
				w.Write(src)
				w.Close()
			}
		})
	}
}

func Decompress(b *testing.B, src []byte) {
	list := []string{gzip.Name, zstd.Name}

	for _, name := range list {
		b.Run(name, func(b *testing.B) {
			c := encoding.GetCompressor(name)
			buf := new(bytes.Buffer)
			w, err := c.Compress(buf)
			if err != nil {
				b.Error(err)
			}
			w.Write(src)
			w.Close()
			enc := buf.Bytes()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				r, err := c.Decompress(bytes.NewReader(enc))
				if err != nil {
					b.Error(err)
				}
				ioutil.ReadAll(r)
			}
		})
	}
}
