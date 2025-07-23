package cobras

import (
	"github.com/mel0dys0ng/song/pkg/cobras/internal"
)

type (
	CommandInterface = internal.CommandInterface
	CbassInterface   = internal.CbassInterface
	CbasInterface    = internal.CbasInterface
	CbaInterface     = internal.CbaInterface
	EmptyCommand     = internal.EmptyCommand
)

// New name root command name.
func New(name string) CbassInterface {
	return internal.New(name)
}
