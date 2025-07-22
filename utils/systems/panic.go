package systems

import (
	"fmt"

	"github.com/song/utils/caller"
)

func Panic(err string) {
	panic(fmt.Sprintf("%s, occurred on %s", err, caller.Location(4, false)))
}

func Panicf(f string, v ...any) {
	msg := fmt.Sprintf(f, v...)
	panic(fmt.Sprintf("%s, occurred on %s", msg, caller.Location(4, false)))
}
