package zstd

import (
	"github.com/takumakei/go-grpc-zstd/zstd"
	"google.golang.org/grpc"
)

func Example() {
	conn, err := grpc.Dial(
		"localhost:9000",
		grpc.WithDefaultCallOptions(grpc.UseCompressor(zstd.Name)),
		grpc.WithInsecure(),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	// Output:
}
