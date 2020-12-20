// +build cgo

package zstd

import (
	"github.com/takumakei/go-grpc-zstd/cgo/zstd"
)

const (
	Name      = zstd.Name
	PoweredBy = zstd.PoweredBy
)

var SetLevel = zstd.SetLevel
