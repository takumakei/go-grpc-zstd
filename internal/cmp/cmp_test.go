package cmp

import (
	"io/ioutil"
	"math/rand"
	"testing"

	datadog "github.com/DataDog/zstd"
	"github.com/google/wuffs/lib/cgozstd"
	"github.com/google/wuffs/lib/compression"
	klauspost "github.com/klauspost/compress/zstd"
)

func BenchmarkZstd(b *testing.B) {
	src := make([]byte, 1024*1024)
	rand.Read(src)

	levels := []struct {
		Name string
		Goog compression.Level
		DDog int
		Klau klauspost.EncoderLevel
	}{
		{
			"default",
			compression.LevelDefault,
			datadog.DefaultCompression,
			klauspost.SpeedDefault,
		},
		{
			"smallest",
			compression.LevelSmallest,
			datadog.BestCompression,
			klauspost.SpeedBestCompression,
		},
		{
			"fastest",
			compression.LevelFastest,
			datadog.BestSpeed,
			klauspost.SpeedFastest,
		},
	}

	for _, level := range levels {
		b.Run(level.Name, func(b *testing.B) {
			b.Run("goog", func(b *testing.B) {
				w := new(cgozstd.Writer)
				w.Reset(ioutil.Discard, nil, level.Goog)

				b.SetBytes(int64(len(src)))
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					w.Write(src)
				}
				w.Close()
			})

			b.Run("ddog", func(b *testing.B) {
				w := datadog.NewWriterLevel(ioutil.Discard, level.DDog)

				b.SetBytes(int64(len(src)))
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					w.Write(src)
				}
				w.Close()
			})

			b.Run("pure", func(b *testing.B) {
				w, err := klauspost.NewWriter(ioutil.Discard, klauspost.WithEncoderLevel(level.Klau))
				if err != nil {
					b.Error(err)
				}

				b.SetBytes(int64(len(src)))
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					w.Write(src)
				}
				w.Close()
			})
		})
	}
}
