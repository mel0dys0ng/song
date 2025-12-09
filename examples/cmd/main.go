package main

import (
	"github.com/mel0dys0ng/song/core/cobras"
	"github.com/mel0dys0ng/song/examples/cmd/internal"
)

const (
	name = "cmdexample"
)

func main() {
	cobras.New(name).RegisterExecute(func(c cobras.CbasInterface) {
		c.RegisterRoot(internal.NewRoot)
	})
}
