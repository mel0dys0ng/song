package internal

type Option struct {
	Apply func(*ErLog)
}

func Log(b bool) Option {
	return Option{
		Apply: func(e *ErLog) {
			e.SetLog(b)
		},
	}
}
