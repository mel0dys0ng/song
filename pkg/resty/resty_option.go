package resty

import (
	"time"

	"github.com/mel0dys0ng/song/internal/core/clients/resty"
)

func OptionDebug(b bool) resty.Option {
	return resty.Debug(b)
}

func OptionBaseURL(s string) resty.Option {
	return resty.BaseURL(s)
}

func OptionDid(s string) resty.Option {
	return resty.Did(s)
}

func OptionTypeIntranet() resty.Option {
	return resty.Type(resty.Intranet)
}

func OptionTypeExtranet() resty.Option {
	return resty.Type(resty.Extranet)
}

func OptionTrace(b bool) resty.Option {
	return resty.Trace(b)
}

func OptionTimeout(t time.Duration) resty.Option {
	return resty.Timeout(t)
}

func OptionRetryCount(i int) resty.Option {
	return resty.RetryCount(i)
}

func OptionRetryWaitTime(t time.Duration) resty.Option {
	return resty.RetryWaitTime(t)
}

func OptionRetryWaitMaxTime(t time.Duration) resty.Option {
	return resty.RetryWaitMaxTime(t)
}

func OptionSignTTL(t int) resty.Option {
	return resty.SignTTL(t)
}

func OptionSignSecret(s string) resty.Option {
	return resty.SignSecret(s)
}
