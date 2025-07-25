package main

import (
	"github.com/mel0dys0ng/song/core/cobras"
	"github.com/mel0dys0ng/song/examples/apiserver/internal"
)

const (
	name = "apiserver-example"
)

func main() {
	cobras.New(name).RegisterExecute(func(c cobras.CbasInterface) {
		c.RegisterRoot(internal.NewHttpServer)
	})
}
