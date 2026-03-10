package https

import "github.com/mel0dys0ng/song/internal/core/https"

func ResponseOptionCode(code int64) https.ResponseOption {
	return func(rsp *https.ResponseData) {
		rsp.Code = code
	}
}

func ResponseOptionMsg(msg string) https.ResponseOption {
	return func(rsp *https.ResponseData) {
		rsp.Msg = msg
	}
}

func ResponseOptionData(data any) https.ResponseOption {
	return func(rsp *https.ResponseData) {
		rsp.Data = data
	}
}

func ResponseOptionTypeJSON() https.ResponseOption {
	return func(rsp *https.ResponseData) {
		rsp.Type = https.ResponseTypeJSON
	}
}

func ResponseOptionTypeJSONP() https.ResponseOption {
	return func(rsp *https.ResponseData) {
		rsp.Type = https.ResponseTypeJSONP
	}
}

func ResponseOptionTypeAsciiJSON() https.ResponseOption {
	return func(rsp *https.ResponseData) {
		rsp.Type = https.ResponseTypeAsciiJSON
	}
}

func ResponseOptionTypeHTML() https.ResponseOption {
	return func(rsp *https.ResponseData) {
		rsp.Type = https.ResponseTypeHTML
	}
}

func ResponseOptionTypeStream() https.ResponseOption {
	return func(rsp *https.ResponseData) {
		rsp.Type = https.ResponseTypeSTREAM
	}
}
