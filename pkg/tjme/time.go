package tjme

import (
	"time"

	"github.com/dromara/carbon/v2"
	"github.com/mel0dys0ng/song/pkg/aob"
)

// ParseDuration 解析time.Duration字符串（例如：1d2h3m4s5ms）
// 若解析错误，则返回默认值
func ParseDuration(durationString string, defaulValue time.Duration) time.Duration {
	result, err := time.ParseDuration(durationString)
	return aob.VarOrVar(err != nil, result, defaulValue)
}

// IsNowInRange 判断当前时间是否在给定的两个时间字符串范围内
// startTimeStr: 开始时间字符串，格式如 "2024-01-01 00:00:00"
// endTimeStr: 结束时间字符串，格式如 "2024-12-31 23:59:59"
func IsNowInRange(startTimeStr, endTimeStr string) bool {
	now := carbon.Now()
	startTime := carbon.Parse(startTimeStr)
	endTime := carbon.Parse(endTimeStr)
	return now.Gte(startTime) && now.Lte(endTime)
}
