package strjngs

import "strings"

// MaskPhoneNumber 对手机号码进行隐私处理
func MaskPhoneNumber(phone string) string {
	if len(phone) != 11 {
		return phone // 如果不是 11 位，直接返回原字符串
	}
	return phone[:3] + "****" + phone[7:]
}

// MaskEmail 对邮箱进行隐私处理
func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email // 如果不是标准邮箱格式，直接返回原字符串
	}

	prefix := parts[0]
	domain := parts[1]

	// 对前缀部分进行处理
	if len(prefix) <= 4 {
		// 如果前缀长度小于等于 4，只保留第一个字符
		return string(prefix[0]) + "****@" + domain
	} else {
		// 保留前 3 位和后 1 位
		return prefix[:3] + "****" + string(prefix[len(prefix)-1]) + "@" + domain
	}
}
