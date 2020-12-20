package smallest_test

import (
	"testing"

	"github.com/google/wuffs/lib/compression"
	"github.com/takumakei/go-grpc-zstd/internal/bench"
	"github.com/takumakei/go-grpc-zstd/zstd"
)

func init() {
	zstd.SetLevel(compression.LevelSmallest)
}

func BenchmarkCompress(b *testing.B) {
	bench.BenchmarkCompress(b)
}

func BenchmarkDecompress(b *testing.B) {
	bench.BenchmarkDecompress(b)
}
