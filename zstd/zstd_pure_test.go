// +build !cgo

package zstd_test

import (
	"testing"

	"github.com/takumakei/go-grpc-zstd/zstd"
)

func TestPoweredBy(t *testing.T) {
	want := "github.com/klauspost/compress/zstd"
	if zstd.PoweredBy != want {
		t.Fatalf("zstd.PoweredBy = %q, want %q", zstd.PoweredBy, want)
	}
}
