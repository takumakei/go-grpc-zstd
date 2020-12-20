package bench_test

import (
	"testing"

	"github.com/takumakei/go-grpc-zstd/internal/bench"
)

func BenchmarkCompress(b *testing.B) {
	bench.BenchmarkCompress(b)
}

func BenchmarkDecompress(b *testing.B) {
	bench.BenchmarkDecompress(b)
}
