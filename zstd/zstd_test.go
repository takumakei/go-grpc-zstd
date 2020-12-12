package zstd_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/google/wuffs/lib/compression"
	"github.com/takumakei/go-grpc-zstd/zstd"
	"github.com/takumakei/go-grpc-zstd/zstd/internal/bench"
	"google.golang.org/grpc/encoding"
)

func Test(t *testing.T) {
	src := make([]byte, 1024*1024)
	bench.Rand(src)

	buf := new(bytes.Buffer)

	levels := []struct {
		Name  string
		Level compression.Level
	}{
		{"(default)", compression.LevelDefault}, // this must be 1st.
		{"fastest", compression.LevelFastest},
		{"fast", compression.LevelFast},
		{"small", compression.LevelSmall},
		{"smallest", compression.LevelSmallest},
		{"default", compression.LevelDefault}, // the last affects other tests.
	}

	for _, level := range levels {
		if level.Name != "(default)" {
			zstd.SetLevel(level.Level)
		}

		c := encoding.GetCompressor(zstd.Name)

		buf.Reset()
		w, err := c.Compress(buf)
		if err != nil {
			t.Error(level.Name, err)
			continue
		}
		w.Write(src)
		w.Close()

		t.Log(level.Name, buf.Len())

		r, err := c.Decompress(buf)
		if err != nil {
			t.Error(level.Name, err)
			continue
		}
		dst, err := ioutil.ReadAll(r)
		if err != nil {
			t.Error(level.Name, err)
			continue
		}

		if bytes.Compare(src, dst) != 0 {
			t.Error(level.Name, "compare", len(src), len(dst))
		}
	}
}
