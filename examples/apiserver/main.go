package main

import (
	"github.com/mel0dys0ng/song/examples/apiserver/internal"
	"github.com/mel0dys0ng/song/pkgs/cobras"
)

const (
	name = "apiserver-example"
)

func main() {
	cobras.New(name).RegisterExecute(func(c cobras.CbasInterface) {
		c.RegisterRoot(internal.NewHttpServer)
	})
}
