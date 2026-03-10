package strjngs

import (
	"crypto/md5"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"html"
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/cast"
)

// IndexOfSplit 返回字符串切割后的片段
func IndexOfSplit(data, depr string, index int) string {
	slices := strings.Split(data, depr)
	length := len(slices)
	if length == 0 {
		return ""
	}

	index %= length
	if index < 0 {
		index += length
	}

	return slices[index]
}

// IndexesOfSplit 返回字符串切割后的多个指定片段
func IndexesOfSplit(data, depr string, indexes ...int) (res []string) {
	slices := strings.Split(data, depr)
	length := len(slices)
	if length == 0 {
		return
	}

	for _, index := range indexes {
		if index %= length; index < 0 {
			index += length
		}

		if index < length {
			res = append(res, slices[index])
		}
	}

	return
}

// JSONMarshal 尝试将任意数据转换为JSON字符串
func JSONMarshal[T any](data T) (res string, err error) {
	bytes, err := json.Marshal(&data)
	if err == nil && len(bytes) > 0 {
		res = string(bytes)
	}
	return
}

// JSONUnmarshal 尝试将JSON字符串转换为任意数据
func JSONUnmarshal[T any](data string, v *T) (err error) {
	if len(data) > 0 {
		err = json.Unmarshal([]byte(data), v)
	}
	return
}

// ConstantTimeCompare 使用常数时间进行字符串比较
func ConstantTimeCompare(a, b string) (res bool) {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// GenerateStableUniqueStr 生成稳定唯一的字符串，通过对参数进行排序后计算MD5哈希。
func GenerateStableUniqueStr(arguments ...any) string {
	type elem struct {
		val any
		str string
	}

	// 预转换所有参数为字符串，减少排序时的重复转换
	elems := make([]elem, len(arguments))
	for i, v := range arguments {
		elems[i] = elem{val: v, str: cast.ToString(v)}
	}

	// 使用预生成的字符串进行稳定排序（降序）
	sort.SliceStable(elems, func(i, j int) bool {
		return elems[i].str > elems[j].str
	})

	// 提取排序后的原始值
	sortedArgs := make([]any, len(elems))
	for i := range elems {
		sortedArgs[i] = elems[i].val
	}

	// 序列化并计算哈希
	bytes, _ := json.Marshal(sortedArgs)
	h := md5.New()
	h.Write(bytes)

	return fmt.Sprintf("%x", h.Sum(nil))
}

// IsEmail 校验邮箱是否合法
func IsEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return re.MatchString(email)
}

// IsMobile 校验手机号是否合法
func IsMobile(mobile string) bool {
	pattern := `^1[3456789]\d{9}$`
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return re.MatchString(mobile)
}

// FilterUserString 对用户输入的字符串进行过滤，移除可能的恶意内容
func FilterUserString(input string) string {
	// 移除URL
	input = regexp.MustCompile(`https?://[^\s]+`).ReplaceAllString(input, "")

	// 移除连续的空格
	input = regexp.MustCompile(`\s+`).ReplaceAllString(input, " ")

	// 对HTML特殊字符进行转义，防止XSS攻击
	input = html.EscapeString(input)

	return strings.TrimSpace(input)
}
