package zstd_test

import (
	"github.com/takumakei/go-grpc-zstd/zstd"
	"google.golang.org/grpc"
)

func Example() {
	conn, err := grpc.Dial(
		"example.com:9000",
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(zstd.Name)),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	// output:
}
