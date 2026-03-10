package metas

const (
	KindAPI       KindType = "api"       // HTTP API
	KindJob       KindType = "job"       // 定时任务
	KindTool      KindType = "tool"      // 工具
	KindMessaging KindType = "messaging" // 消息队列
)

// 应用类型
type KindType string

func (kt KindType) Validate() bool {
	switch kt {
	case KindAPI, KindJob, KindTool, KindMessaging:
		return true
	default:
		return false
	}
}

func (kt KindType) String() string {
	return string(kt)
}
