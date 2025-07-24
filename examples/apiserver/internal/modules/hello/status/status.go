package status

import "github.com/mel0dys0ng/song/pkgs/erlogs"

var (
	ServerErrorMsg          = "系统网络异常,请稍后重试"
	TooFrequentOperationMsg = "操作太频繁,请稍后重试"
	InvalidArgumentsMsg     = "请求参数错误"
	InvalidRequestMsg       = "请求参数错误"
	InvalidArguments        = erlogs.Common.WithStatus(40001, InvalidArgumentsMsg)
	Unauthorized            = erlogs.Common.WithStatus(40003, "请登录")
	Forbidden               = erlogs.Common.WithStatus(40004, "请求被拒绝")
	Unknown                 = erlogs.Common.WithStatus(50000, ServerErrorMsg)
	InvalidRequest          = erlogs.Common.WithStatus(50001, InvalidRequestMsg)
)
