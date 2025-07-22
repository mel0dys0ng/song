package clients

import (
	"time"

	"github.com/mel0dys0ng/song/clients/internal/resty"
)

func RestyOptionDebug(b bool) resty.Option {
	return resty.Debug(b)
}

func RestyOptionBaseURL(s string) resty.Option {
	return resty.BaseURL(s)
}

func RestyOptionDid(s string) resty.Option {
	return resty.Did(s)
}

func RestyOptionTypeIntranet() resty.Option {
	return resty.Type(resty.Intranet)
}

func RestyOptionTypeExtranet() resty.Option {
	return resty.Type(resty.Extranet)
}

func RestyOptionTrace(b bool) resty.Option {
	return resty.Trace(b)
}

func RestyOptionTimeout(t time.Duration) resty.Option {
	return resty.Timeout(t)
}

func RestyOptionRetryCount(i int) resty.Option {
	return resty.RetryCount(i)
}

func RestyOptionRetryWaitTime(t time.Duration) resty.Option {
	return resty.RetryWaitTime(t)
}

func RestyOptionRetryWaitMaxTime(t time.Duration) resty.Option {
	return resty.RetryWaitMaxTime(t)
}

func RestyOptionSignTTL(t time.Duration) resty.Option {
	return resty.SignTTL(t)
}

func RestyOptionSignSecret(s string) resty.Option {
	return resty.SignSecret(s)
}
