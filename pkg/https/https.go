package https

import (
	"github.com/gin-gonic/gin"
	"github.com/mel0dys0ng/song/internal/core/https"
)

type (
	Option               = https.Option
	Route                = https.Route
	PriorityMiddleware   = https.Middleware
	MiddlewareHandleFunc = https.MiddlewareHandleFunc
	RequestResponseData  = https.RequestResponseData
	RecoveredData        = https.RecoveredData
	ClientInfo           = https.ClientInfo
	ResponseData         = https.ResponseData
	ResponseOption       = https.ResponseOption
)

// Constants
const (
	// OS constants
	OSWindows = https.OSWindows
	OSMacOS   = https.OSMacOS
	OSIOS     = https.OSIOS
	OSAndroid = https.OSAndroid
	OSLinux   = https.OSLinux
	OSUnix    = https.OSUnix
	OSUnknown = https.OSUnknown

	// Browser constants
	BrowserChrome           = https.BrowserChrome
	BrowserEdge             = https.BrowserEdge
	BrowserFirefox          = https.BrowserFirefox
	BrowserSafari           = https.BrowserSafari
	BrowserInternetExplorer = https.BrowserInternetExplorer
	BrowserOpera            = https.BrowserOpera
	BrowserSafariMobile     = https.BrowserSafariMobile
	BrowserChromeMobile     = https.BrowserChromeMobile
	BrowserUCBrowser        = https.BrowserUCBrowser
	BrowserSamsung          = https.BrowserSamsung
	BrowserUnknown          = https.BrowserUnknown

	// Client type constants
	ClientTypeUnknown    = https.ClientTypeUnknown
	ClientTypeWeb        = https.ClientTypeWeb
	ClientTypeMobileApp  = https.ClientTypeMobileApp
	ClientTypeDesktopApp = https.ClientTypeDesktopApp

	// CSRF constants
	CSRFDefaultEnable         = https.CSRFDefaultEnable
	CSRFDefaultLookupType     = https.CSRFDefaultLookupType
	CSRFDefaultLookupName     = https.CSRFDefaultLookupName
	CSRFDefaultCookieName     = https.CSRFDefaultCookieName
	CSRFDefaultCookieDomain   = https.CSRFDefaultCookieDomain
	CSRFDefaultCookiePath     = https.CSRFDefaultCookiePath
	CSRFDefaultCookieMaxAge   = https.CSRFDefaultCookieMaxAge
	CSRFDefaultCookieSecure   = https.CSRFDefaultCookieSecure
	CSRFDefaultCookieHttpOnly = https.CSRFDefaultCookieHttpOnly

	LookupTypeHeader = https.LookupTypeHeader
	LookupTypeForm   = https.LookupTypeForm
	LookupTypeQuery  = https.LookupTypeQuery

	LookupNameHeader = https.LookupNameHeader
	LookupNameForm   = https.LookupNameForm
	LookupNameQuery  = https.LookupNameQuery

	// Response constants
	ResponseSuccessCode = https.ResponseSuccessCode
	ResponseUnknownCode = https.ResponseUnknownCode

	// Sign constants
	DefaultSignSecret = https.DefaultSignSecret
	DefaultSignTTL    = https.DefaultSignTTL

	// Config constants
	DefaultPort              = https.DefaultPort
	DefaultHost              = https.DefaultHost
	DefaultHttpsOpen         = https.DefaultHttpsOpen
	DefaultHttpsKeyFile      = https.DefaultHttpsKeyFile
	DefaultHttpsKeyCert      = https.DefaultHttpsKeyCert
	DefaultKeepAlive         = https.DefaultKeepAlive
	DefaultReadTimeout       = https.DefaultReadTimeout
	DefaultReadHeaderTimeout = https.DefaultReadHeaderTimeout
	DefaultWriteTimeout      = https.DefaultWriteTimeout
	DefaultIdleTimeout       = https.DefaultIdleTimeout
	DefaultHammerTime        = https.DefaultHammerTime
	DefaultMaxHeaderBytes    = https.DefaultMaxHeaderBytes
	DefaultTmpDir            = https.DefaultTmpDir
	DefaultCorsMaxAge        = https.DefaultCorsMaxAge
)

// New return a new HTTP Server
// @Param opts []Option the option of http server
func New(opts []Option) *https.Server {
	return https.New(opts)
}

// GetClientInfo 获取客户端信息
func GetClientInfo(ctx *gin.Context) (res *https.ClientInfo) {
	return https.GetClientInfo(ctx)
}

// Response 请求响应。
// @Param data any 响应数据
// @param err error 请求错误。成功时，nil或者级别低于warning的错误；失败时，不为nil且级别高于warning的错误
// @Param opts []inetrnal.ResponseOption 自定义响应选项，可覆盖默认根据err和data设定的选项
func Response(ctx *gin.Context, data any, err error, opts ...https.ResponseOption) {
	https.Response(ctx, data, err, opts...)
}

// ResponseSuccess 请求成功时的响应
// @param data any 响应数据
// @Param opts []inetrnal.ResponseOption 自定义响应选项，可覆盖默认根据data设定的选项
func ResponseSuccess(ctx *gin.Context, data any, opts ...https.ResponseOption) {
	https.ResponseSuccess(ctx, data, opts...)
}

// ResponseError 请求失败（发生错误或异常）时的响应
// @param err error 请求错误。成功时，nil或者级别低于warning的错误；失败时，不为nil且级别高于warning的错误
// @Param opts []inetrnal.ResponseOption 自定义响应选项，可覆盖默认根据data设定的选项
func ResponseError(ctx *gin.Context, err error, opts ...https.ResponseOption) {
	https.ResponseError(ctx, err, opts...)
}

// ResponseWithStatus 自定义http status响应，默认JSON格式
func ResponseWithStatus(ctx *gin.Context, status int, opts ...https.ResponseOption) {
	https.ResponseWithStatus(ctx, status, opts...)
}

// ResponseFromContext 从上下文中获取响应数据
func ResponseFromContext(ctx *gin.Context) *ResponseData {
	return https.ResponseFromContext(ctx)
}

// GenerateSign 为给定的数据生成签名，可用于客户端生成签名
// 签名规则:
// 1. 将数据按key字母顺序排序
// 2. 按照 secret@ttl#key1=value1&key2=value2... 的格式拼接字符串
// 3. 对拼接的字符串进行MD5加密
// 4. 如果secret为空则使用默认密钥(DefaultSignSecret)
// 5. 如果ttl为0则使用默认TTL(DefaultSignTTL)
// 注意: 数据中的HeaderKeySign字段会被排除
func GenerateSign(data map[string]string, secret string, ttl int) string {
	return https.GenerateSign(data, secret, ttl)
}
