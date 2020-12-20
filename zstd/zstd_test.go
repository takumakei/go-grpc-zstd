package zstd_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/takumakei/go-grpc-zstd/zstd"
	"google.golang.org/grpc/encoding"
)

func TestName(t *testing.T) {
	want := "zstd"
	if zstd.Name != want {
		t.Fatalf("zstd.Name = %q, want %q", zstd.Name, want)
	}
}

func TestCompressor(t *testing.T) {
	c := encoding.GetCompressor(zstd.Name)

	src := "hello world"
	buf := new(bytes.Buffer)

	zw, _ := c.Compress(buf)
	zw.Write([]byte(src))
	zw.Close()

	zr, _ := c.Decompress(buf)
	dst, _ := ioutil.ReadAll(zr)

	if res := string(dst); res != src {
		t.Fatalf("%q != %q", res, src)
	}
}

func TestCompressor_pipe(t *testing.T) {
	c := encoding.GetCompressor(zstd.Name)

	r, w := io.Pipe()
	zw, err := c.Compress(w)
	if err != nil {
		t.Fatal(err)
	}
	zr, err := c.Decompress(r)
	if err != nil {
		t.Fatal(err)
	}
	dst := new(bytes.Buffer)
	go io.Copy(dst, zr)
	src := "hello world"
	zw.Write([]byte(src))
	zw.Close()
	w.Close()
	if res := dst.String(); res != src {
		t.Fatalf("%q != %q", res, src)
	}
}
