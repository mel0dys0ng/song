package song

import (
	"github.com/mel0dys0ng/song/internal/core/cobras"
)

type (
	CommandInterface        = cobras.CommandInterface
	CommandRunnerInterface  = cobras.CommandRunnerInterface
	CommandBuilderInterface = cobras.CommandBuilderInterface
	CommandChildInterface   = cobras.CommandChildInterface
	Empty                   = cobras.Empty
)

// New return a new cobra command runner
func New(name string) CommandRunnerInterface {
	return cobras.New(name)
}
