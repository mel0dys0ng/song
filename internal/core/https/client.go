package https

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ClientDeviceIDHeaderKey = "X-Song-Device-Id"
	// ClientInfoContextKey 上下文中客户端信息结构体的key
	ClientInfoContextKey = "X-Song-Client-Info"
	// ForwardedForHeader X-Forwarded-For头部字段
	ForwardedForHeader = "X-Forwarded-For"
	// RealIPHeader X-Real-Ip头部字段
	RealIPHeader = "X-Real-Ip"
	// DeviceIDAlt1Header X-DeviceId头部字段
	DeviceIDAlt1Header = "X-DeviceId"
	// DeviceIDAlt2Header Device-ID头部字段
	DeviceIDAlt2Header = "Device-ID"
	// DeviceIDAlt3Header X-Device-Id头部字段
	DeviceIDAlt3Header = "X-Device-Id"

	// 操作系统常量
	OSWindows = "Windows"
	OSMacOS   = "macOS"
	OSIOS     = "iOS"
	OSAndroid = "Android"
	OSLinux   = "Linux"
	OSUnix    = "Unix"
	OSUnknown = "Unknown"

	// 浏览器常量
	BrowserChrome           = "Chrome"
	BrowserEdge             = "Edge"
	BrowserFirefox          = "Firefox"
	BrowserSafari           = "Safari"
	BrowserInternetExplorer = "Internet Explorer"
	BrowserOpera            = "Opera"
	BrowserSafariMobile     = "Safari Mobile"
	BrowserChromeMobile     = "Chrome Mobile"
	BrowserUCBrowser        = "UC Browser"
	BrowserSamsung          = "Samsung Browser"
	BrowserUnknown          = "Unknown"

	// 客户端类型常量
	ClientTypeUnknown = iota
	ClientTypeWeb
	ClientTypeMobileApp
	ClientTypeDesktopApp
)

// ClientInfo 客户端信息结构体
type ClientInfo struct {
	IP             string `json:"ip"`
	DeviceID       string `json:"device_id"`
	UserAgent      string `json:"user_agent"`
	Method         string `json:"method"`
	Path           string `json:"path"`
	OS             string `json:"os"`
	OSVersion      string `json:"os_version"`
	Browser        string `json:"browser"`
	BrowserVersion string `json:"browser_version"`
	ClientType     int    `json:"client_type"`    // 客户端类型 (Web, Mobile App, Desktop App)
	ClientVersion  string `json:"client_version"` // 客户端版本号
	AppBuild       string `json:"app_build"`      // 应用构建号
}

// IP getter
func (c *ClientInfo) GetIP() string {
	if c == nil {
		return ""
	}
	return c.IP
}

// DeviceID getter
func (c *ClientInfo) GetDeviceID() string {
	if c == nil {
		return ""
	}
	return c.DeviceID
}

// UserAgent getter
func (c *ClientInfo) GetUserAgent() string {
	if c == nil {
		return ""
	}
	return c.UserAgent
}

// Method getter
func (c *ClientInfo) GetMethod() string {
	if c == nil {
		return ""
	}
	return c.Method
}

// Path getter
func (c *ClientInfo) GetPath() string {
	if c == nil {
		return ""
	}
	return c.Path
}

// OS getter
func (c *ClientInfo) GetOS() string {
	if c == nil {
		return ""
	}
	return c.OS
}

// OSVersion getter
func (c *ClientInfo) GetOSVersion() string {
	if c == nil {
		return ""
	}
	return c.OSVersion
}

// Browser getter
func (c *ClientInfo) GetBrowser() string {
	if c == nil {
		return ""
	}
	return c.Browser
}

// BrowserVersion getter
func (c *ClientInfo) GetBrowserVersion() string {
	if c == nil {
		return ""
	}
	return c.BrowserVersion
}

// ClientType getter
func (c *ClientInfo) GetClientType() int {
	if c == nil {
		return ClientTypeUnknown
	}
	return c.ClientType
}

// ClientVersion getter
func (c *ClientInfo) GetClientVersion() string {
	if c == nil {
		return ""
	}
	return c.ClientVersion
}

// AppBuild getter
func (c *ClientInfo) GetAppBuild() string {
	if c == nil {
		return ""
	}
	return c.AppBuild
}

// IsWindows 判断是否为Windows操作系统
func (c *ClientInfo) IsWindows() bool {
	if c == nil {
		return false
	}
	return c.OS == OSWindows
}

// IsMacOS 判断是否为macOS操作系统
func (c *ClientInfo) IsMacOS() bool {
	if c == nil {
		return false
	}
	return c.OS == OSMacOS
}

// IsIOS 判断是否为iOS操作系统
func (c *ClientInfo) IsIOS() bool {
	if c == nil {
		return false
	}
	return c.OS == OSIOS
}

// IsAndroid 判断是否为Android操作系统
func (c *ClientInfo) IsAndroid() bool {
	if c == nil {
		return false
	}
	return c.OS == OSAndroid
}

// IsLinux 判断是否为Linux操作系统
func (c *ClientInfo) IsLinux() bool {
	if c == nil {
		return false
	}
	return c.OS == OSLinux
}

// IsMobileOS 判断是否为移动操作系统
func (c *ClientInfo) IsMobileOS() bool {
	if c == nil {
		return false
	}
	return c.IsIOS() || c.IsAndroid()
}

// IsDesktopOS 判断是否为桌面操作系统
func (c *ClientInfo) IsDesktopOS() bool {
	if c == nil {
		return false
	}
	return c.IsWindows() || c.IsMacOS() || c.IsLinux()
}

// IsChrome 判断是否为Chrome浏览器
func (c *ClientInfo) IsChrome() bool {
	if c == nil {
		return false
	}
	return c.Browser == BrowserChrome || c.Browser == BrowserChromeMobile
}

// IsSafari 判断是否为Safari浏览器
func (c *ClientInfo) IsSafari() bool {
	if c == nil {
		return false
	}
	return c.Browser == BrowserSafari || c.Browser == BrowserSafariMobile
}

// IsFirefox 判断是否为Firefox浏览器
func (c *ClientInfo) IsFirefox() bool {
	if c == nil {
		return false
	}
	return c.Browser == BrowserFirefox
}

// IsEdge 判断是否为Edge浏览器
func (c *ClientInfo) IsEdge() bool {
	if c == nil {
		return false
	}
	return c.Browser == BrowserEdge
}

// IsMobileBrowser 判断是否为移动浏览器
func (c *ClientInfo) IsMobileBrowser() bool {
	if c == nil {
		return false
	}
	return c.Browser == BrowserChromeMobile || c.Browser == BrowserSafariMobile
}

// IsWebClient 判断是否为Web客户端
func (c *ClientInfo) IsWebClient() bool {
	if c == nil {
		return false
	}
	return c.ClientType == ClientTypeWeb
}

// IsMobileApp 判断是否为移动应用客户端
func (c *ClientInfo) IsMobileApp() bool {
	if c == nil {
		return false
	}
	return c.ClientType == ClientTypeMobileApp
}

// IsDesktopApp 判断是否为桌面应用客户端
func (c *ClientInfo) IsDesktopApp() bool {
	if c == nil {
		return false
	}
	return c.ClientType == ClientTypeDesktopApp
}

// IsKnownOS 判断是否为已知的操作系统
func (c *ClientInfo) IsKnownOS() bool {
	if c == nil {
		return false
	}
	return c.OS != OSUnknown
}

// IsKnownBrowser 判断是否为已知的浏览器
func (c *ClientInfo) IsKnownBrowser() bool {
	if c == nil {
		return false
	}
	return c.Browser != BrowserUnknown
}

// parseVersion 解析版本号字符串为数字数组
func parseVersion(version string) []int {
	if version == "" {
		return []int{}
	}

	parts := strings.Split(version, ".")
	result := make([]int, len(parts))

	for i, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil {
			// 如果无法转换为整数，跳过这个部分
			result[i] = 0
		} else {
			result[i] = num
		}
	}

	return result
}

// CompareVersion 比较客户端版本号，返回值含义：
// 0: 相同版本
// 1: 当前版本大于指定版本
// -1: 当前版本小于指定版本
func (c *ClientInfo) CompareVersion(targetVersion string) int {
	if c == nil {
		return -1 // 认为nil版本小于任何有效版本
	}

	currentParts := parseVersion(c.ClientVersion)
	targetParts := parseVersion(targetVersion)
	maxLen := max(len(targetParts), len(currentParts))

	for i := range maxLen {
		var cur, tar int
		if i < len(currentParts) {
			cur = currentParts[i]
		}

		if i < len(targetParts) {
			tar = targetParts[i]
		}

		if cur > tar {
			return 1
		} else if cur < tar {
			return -1
		}
	}

	return 0
}

// IsVersionEqual 检查客户端版本是否等于指定版本
func (c *ClientInfo) IsVersionEqual(targetVersion string) bool {
	if c == nil {
		return false
	}
	return c.CompareVersion(targetVersion) == 0
}

// IsVersionGreater 检查客户端版本是否大于指定版本
func (c *ClientInfo) IsVersionGreater(targetVersion string) bool {
	if c == nil {
		return false
	}
	return c.CompareVersion(targetVersion) == 1
}

// IsVersionGreaterOrEqual 检查客户端版本是否大于或等于指定版本
func (c *ClientInfo) IsVersionGreaterOrEqual(targetVersion string) bool {
	if c == nil {
		return false
	}
	return c.CompareVersion(targetVersion) >= 0
}

// IsVersionLess 检查客户端版本是否小于指定版本
func (c *ClientInfo) IsVersionLess(targetVersion string) bool {
	if c == nil {
		return true // nil版本被认为小于任何有效版本
	}
	return c.CompareVersion(targetVersion) == -1
}

// IsVersionLessOrEqual 检查客户端版本是否小于或等于指定版本
func (c *ClientInfo) IsVersionLessOrEqual(targetVersion string) bool {
	if c == nil {
		return true // nil版本被认为小于等于任何有效版本
	}
	return c.CompareVersion(targetVersion) <= 0
}

// setupClientMiddleware 客户端信息中间件，获取并设置客户端相关信息
func (s *Server) setupClientMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userAgent := ctx.Request.UserAgent()

		clientInfo := &ClientInfo{
			IP:        getClientIP(ctx),
			UserAgent: userAgent,
			Method:    ctx.Request.Method,
			Path:      ctx.Request.URL.Path,
			DeviceID:  getDeviceID(ctx),
		}

		// 解析操作系统和浏览器信息
		os, osVer, browser, browserVer := parseUserAgentInfo(userAgent)
		clientInfo.OS = os
		clientInfo.OSVersion = osVer
		clientInfo.Browser = browser
		clientInfo.BrowserVersion = browserVer

		// 解析客户端类型和版本信息
		clientInfo.ClientType, clientInfo.ClientVersion, clientInfo.AppBuild = parseClientVersionInfo(userAgent)

		// 将客户端信息设置到上下文中
		ctx.Set(ClientInfoContextKey, clientInfo)

		ctx.Next()
	}
}

// parseClientVersionInfo 从User-Agent解析客户端类型和版本信息
func parseClientVersionInfo(userAgent string) (clientType int, version, build string) {
	if userAgent == "" {
		return ClientTypeUnknown, "", ""
	}

	// 检查是否是移动APP
	if strings.Contains(userAgent, "MobileApp/") || strings.Contains(userAgent, "AppName/") {
		// 尝试匹配 "AppName/version.build" 或 "AppName/version" 格式
		appPattern := regexp.MustCompile(`AppName/([0-9]+\.[0-9]+\.[0-9]+)\.([0-9]+)`) // 版本.构建号
		matches := appPattern.FindStringSubmatch(userAgent)
		if len(matches) >= 3 {
			return ClientTypeMobileApp, matches[1], matches[2]
		}

		// 尝试匹配 "AppName/version" 格式
		appPatternSimple := regexp.MustCompile(`AppName/([0-9]+\.[0-9]+\.[0-9]+)`)
		matches = appPatternSimple.FindStringSubmatch(userAgent)
		if len(matches) >= 2 {
			return ClientTypeMobileApp, matches[1], ""
		}

		// 检查是否包含移动标识但没有版本号
		if strings.Contains(userAgent, "Mobile") || strings.Contains(userAgent, "Android") || strings.Contains(userAgent, "iPhone") {
			return ClientTypeMobileApp, "", ""
		}
	}

	// 默认为Web客户端
	return ClientTypeWeb, "", ""
}

// parseUserAgentInfo 从User-Agent字符串解析操作系统和浏览器信息
func parseUserAgentInfo(userAgent string) (os, osVersion, browser, browserVersion string) {
	if userAgent == "" {
		return
	}

	// 解析操作系统
	os, osVersion = detectOS(userAgent)

	// 解析浏览器
	browser, browserVersion = detectBrowser(userAgent)

	return
}

// detectOS 从User-Agent中检测操作系统
func detectOS(userAgent string) (os, version string) {
	lowerUA := strings.ToLower(userAgent)

	// Windows系统
	winPattern := regexp.MustCompile(`Windows NT (\d+\.\d+)`)
	if matches := winPattern.FindStringSubmatch(userAgent); len(matches) > 1 {
		return OSWindows, matches[1]
	}

	// macOS/iOS系统
	if strings.Contains(lowerUA, "mac os x") {
		macPattern := regexp.MustCompile(`Mac OS X (\d+[._]\d+(?:[._]\d+)?)`)
		if matches := macPattern.FindStringSubmatch(userAgent); len(matches) > 1 {
			version = strings.ReplaceAll(matches[1], "_", ".")
			return OSMacOS, version
		}
		return OSMacOS, ""
	}

	if strings.Contains(lowerUA, "iphone") || strings.Contains(lowerUA, "ipad") || strings.Contains(lowerUA, "ios") {
		iosPattern := regexp.MustCompile(`(iPhone|iPad|iOS) OS (\d+[._]\d+(?:[._]\d+)?)`)
		if matches := iosPattern.FindStringSubmatch(userAgent); len(matches) > 2 {
			version = strings.ReplaceAll(matches[2], "_", ".")
			return OSIOS, version
		}
		return OSIOS, ""
	}

	// Android系统
	androidPattern := regexp.MustCompile(`Android[ /](\d+(?:\.\d+(?:\.\d+)?)?)`)
	if matches := androidPattern.FindStringSubmatch(userAgent); len(matches) > 1 {
		return OSAndroid, matches[1]
	}
	if strings.Contains(lowerUA, "android") {
		return OSAndroid, ""
	}

	// Linux系统
	if strings.Contains(lowerUA, "linux") {
		return OSLinux, ""
	}

	// Unix系统
	if strings.Contains(lowerUA, "unix") || strings.Contains(lowerUA, "sunos") {
		return OSUnix, ""
	}

	return OSUnknown, ""
}

// detectBrowser 从User-Agent中检测浏览器
func detectBrowser(userAgent string) (browser, version string) {
	// 移动版Chrome浏览器
	if strings.Contains(strings.ToLower(userAgent), "chrome") && (strings.Contains(strings.ToLower(userAgent), "mobile") || strings.Contains(strings.ToLower(userAgent), "android")) {
		chromePattern := regexp.MustCompile(`Chrome[/]((\d+)(\.[\d]+)+)`)
		if matches := chromePattern.FindStringSubmatch(userAgent); len(matches) > 1 {
			return BrowserChromeMobile, matches[1]
		}
	}

	// Chrome浏览器
	chromePattern := regexp.MustCompile(`Chrome[/]((\d+)(\.[\d]+)+)`)
	if matches := chromePattern.FindStringSubmatch(userAgent); len(matches) > 1 {
		return BrowserChrome, matches[1]
	}

	// Edge浏览器
	edgePattern := regexp.MustCompile(`Edg[/]((\d+)(\.[\d]+)+)`)
	if matches := edgePattern.FindStringSubmatch(userAgent); len(matches) > 1 {
		return BrowserEdge, matches[1]
	}

	// Firefox浏览器
	firefoxPattern := regexp.MustCompile(`Firefox[/](\d+(?:\.\d+)*)`)
	if matches := firefoxPattern.FindStringSubmatch(userAgent); len(matches) > 1 {
		return BrowserFirefox, matches[1]
	}

	// Safari Mobile (iOS移动版Safari)
	if strings.Contains(userAgent, "Safari") && (strings.Contains(userAgent, "Mobile") || strings.Contains(userAgent, "iPhone") || strings.Contains(userAgent, "iPad")) {
		safariPattern := regexp.MustCompile(`Version[/](\d+(?:\.\d+)*)[^ ]* Safari`)
		if matches := safariPattern.FindStringSubmatch(userAgent); len(matches) > 1 {
			return BrowserSafariMobile, matches[1]
		}
		return BrowserSafariMobile, ""
	}

	// Safari浏览器 (注意要排除Chrome中的Safari标记)
	if strings.Contains(userAgent, "Safari") && !strings.Contains(userAgent, "Chrome") {
		safariPattern := regexp.MustCompile(`Version[/](\d+(?:\.\d+)*)[^ ]* Safari`)
		if matches := safariPattern.FindStringSubmatch(userAgent); len(matches) > 1 {
			return BrowserSafari, matches[1]
		}
		return BrowserSafari, ""
	}

	// Internet Explorer
	iePattern1 := regexp.MustCompile(`MSIE (\d+(?:\.\d+)*)`)
	if matches := iePattern1.FindStringSubmatch(userAgent); len(matches) > 1 {
		return BrowserInternetExplorer, matches[1]
	}

	iePattern2 := regexp.MustCompile(`Trident.*rv:(\d+(?:\.\d+)*)`)
	if matches := iePattern2.FindStringSubmatch(userAgent); len(matches) > 1 {
		return BrowserInternetExplorer, matches[1]
	}

	// Opera浏览器
	operaPattern := regexp.MustCompile(`OPR[/](\d+(?:\.\d+)*)`)
	if matches := operaPattern.FindStringSubmatch(userAgent); len(matches) > 1 {
		return BrowserOpera, matches[1]
	}

	// UC浏览器
	ucPattern := regexp.MustCompile(`UCBrowser[/](\d+(?:\.\d+)*)`)
	if matches := ucPattern.FindStringSubmatch(userAgent); len(matches) > 1 {
		return BrowserUCBrowser, matches[1]
	}

	// Samsung浏览器
	samsungPattern := regexp.MustCompile(`SamsungBrowser[/](\d+(?:\.\d+)*)`)
	if matches := samsungPattern.FindStringSubmatch(userAgent); len(matches) > 1 {
		return BrowserSamsung, matches[1]
	}

	return BrowserUnknown, ""
}

// GetClientInfo 获取客户端信息
func GetClientInfo(ctx *gin.Context) (res *ClientInfo) {
	clientInfo, _ := ctx.Get(ClientInfoContextKey)
	res, _ = clientInfo.(*ClientInfo)
	return
}

// getClientIP 获取客户端真实IP地址
func getClientIP(ctx *gin.Context) string {
	// 检查 X-Forwarded-For 头
	xForwardedFor := ctx.GetHeader(ForwardedForHeader)
	if xForwardedFor != "" {
		// X-Forwarded-For 可能包含多个IP地址，以逗号分隔，第一个是真实客户端IP
		ips := strings.Split(xForwardedFor, ",")
		ip := strings.TrimSpace(ips[0])
		if ip != "" {
			return ip
		}
	}

	// 检查 X-Real-Ip 头
	xRealIP := ctx.GetHeader(RealIPHeader)
	if xRealIP != "" {
		return xRealIP
	}

	// 返回远程地址的IP部分
	return ctx.ClientIP()
}

// getDeviceID 获取设备ID
func getDeviceID(ctx *gin.Context) string {
	// 首先检查标准设备ID头部
	deviceID := ctx.GetHeader(ClientDeviceIDHeaderKey)
	if deviceID != "" {
		return deviceID
	}

	// 检查常见的设备ID头部变体
	deviceID = ctx.GetHeader(DeviceIDAlt1Header)
	if deviceID != "" {
		return deviceID
	}

	// 检查可能在请求头中的其他设备标识
	deviceID = ctx.GetHeader(DeviceIDAlt2Header)
	if deviceID != "" {
		return deviceID
	}

	// 检查可能在请求头中的其他设备标识
	deviceID = ctx.GetHeader(DeviceIDAlt3Header)
	if deviceID != "" {
		return deviceID
	}

	// 检查可能在查询参数中的设备ID
	deviceID = ctx.Query("device_id")
	if deviceID != "" {
		return deviceID
	}

	// 检查表单中的设备ID
	if ctx.Request.Method == "POST" || ctx.Request.Method == "PUT" {
		deviceID = ctx.PostForm("device_id")
		if deviceID != "" {
			return deviceID
		}
	}

	// 如果都没有找到，则返回空字符串
	return ""
}
