package mysql

import "github.com/mel0dys0ng/song/pkg/erlogs"

func BaseELOptions() []erlogs.Option {
	return []erlogs.Option{
		erlogs.OptionKindSystem(),
		erlogs.ClientsMySQLBiz,
	}
}
