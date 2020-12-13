package no_level

import (
	"testing"

	"github.com/takumakei/go-grpc-zstd/pure/zstd/internal/bench"
)

func BenchmarkCompress(b *testing.B) {
	bench.BenchmarkCompress(b, 1024*1024)
}

func BenchmarkDecompress(b *testing.B) {
	bench.BenchmarkDecompress(b, 1024*1024)
}
