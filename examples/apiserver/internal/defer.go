package internal

import "github.com/mel0dys0ng/song/pkgs/https"

func SetupDefers() https.Option {
	return https.Defers()
}
