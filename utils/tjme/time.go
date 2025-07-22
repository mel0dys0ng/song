package tjme

import (
	"time"

	"github.com/mel0dys0ng/song/utils/sljces"
)

// ParseDuration 解析time.Duration字符串（例如：1d2h3m4s5ms）
// 若提供了defaultValues，则取第一个为默认值
func ParseDuration(durationString string, defaulValues ...time.Duration) time.Duration {
	result, err := time.ParseDuration(durationString)
	if err != nil {
		return sljces.First(defaulValues, 0)
	}
	return result
}
