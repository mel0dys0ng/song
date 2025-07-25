package main

import "github.com/mel0dys0ng/song/core/cobras"

const (
	app = "song"
)

func main() {
	cobras.New(app).RegisterExecute(func(c cobras.CbasInterface) {

	})
}
