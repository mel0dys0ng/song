package main

import (
	"github.com/mel0dys0ng/song/core/cobras"
	"github.com/mel0dys0ng/song/examples/apiserver/internal"
)

const (
	// 应用名称
	name = "song-apiserver-example"
)

func main() {
	cobras.New(name).RegisterExecute(func(c cobras.CbasInterface) {
		c.RegisterRoot(internal.NewHttpServer)
	})
}
