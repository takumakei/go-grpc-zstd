package default_level

import (
	"testing"

	"github.com/google/wuffs/lib/compression"
	"github.com/takumakei/go-grpc-zstd/pure/zstd/internal/bench"
)

func init() {
	bench.SetLevel(compression.LevelDefault)
}

func BenchmarkCompress(b *testing.B) {
	bench.BenchmarkCompress(b, 1024*1024)
}

func BenchmarkDecompress(b *testing.B) {
	bench.BenchmarkDecompress(b, 1024*1024)
}
