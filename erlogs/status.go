package erlogs

var (
	ComBiz  = Biz(0, "common")
	Common  = New(Log(true), TypeBiz(), ComBiz)
	Ok      = Common.WithStatus(0, "ok")
	Unknown = Common.WithStatus(50000, "系统网络异常，请稍后重试")
)
