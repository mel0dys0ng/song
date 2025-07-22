package cljent

import (
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/song/utils/crypto"
	"github.com/ua-parser/uap-go/uaparser"
)

const (
	ClientTypeUnknown           = 0  // 未知
	ClientTypePCWeb             = 1  // PC Web
	ClentTypeHarmonyWeb         = 2  // Harmony Web
	ClentTypeIOSWeb             = 3  // IOS Web
	ClentTypeAndroidWeb         = 4  // Android Web
	ClentTypeHarmonyApp         = 5  // Harmony App
	ClentTypeIOSApp             = 6  // IOS App
	ClentTypeAndroidApp         = 7  // Android App
	ClentTypeHarmonyAppH5       = 8  // Harmony App Inner H5
	ClentTypeIOSAppH5           = 9  // IOS App Inner H5
	ClentTypeAndroidAppH5       = 10 // Android App Inner H5
	ClientTypeWeixinMiniProgram = 11 // 微信小程序
	ClientTypeDouyinMiniProgram = 12 // 抖音小程序
	ClientTypeQQMiniProgram     = 13 // QQ小程序
	ClientTypeMax               = ClientTypeQQMiniProgram

	HeaderNameDTV = "X-Song-Dtv"
)

type ClientInfo struct {
	ClientDID              string `gorm:"column:client_did" json:"client_did"`                               // 客户端设备ID
	ClientType             uint8  `gorm:"column:client_type" json:"client_type"`                             // 客户端类型，0 : 未知, 1 : PC Web, 2 : Harmony Web, 3 : IOS Web, 4 : Android Web, 5 : Harmony App, 6 : IOS App, 7 : Android App, 8 : Harmony App Webview, 9 : IOS App Webview, 10: Android App Webview, 11: 微信小程序, 12: 抖音小程序, 13: QQ小程序
	ClientVersion          string `gorm:"column:client_version" json:"client_version"`                       // 客户端版本
	ClientDevice           string `gorm:"column:client_device" json:"client_device"`                         // 客户端设备
	ClientDeviceBrand      string `gorm:"column:client_device_brand" json:"client_device_brand"`             // 客户端设备brand
	ClientDeviceModel      string `gorm:"column:client_device_model" json:"client_device_model"`             // 客户端设备model
	ClientOS               string `gorm:"column:client_os" json:"client_os"`                                 // 客户端系统
	ClientOSVersion        string `gorm:"column:client_os_version" json:"client_os_version"`                 // 客户端系统版本
	ClientUserAgent        string `gorm:"column:client_user_agent" json:"client_user_agent"`                 // 客户端UA
	ClientUserAgentName    string `gorm:"column:client_user_agent_name" json:"client_user_agent_name"`       // 客户端UA名称
	ClientUserAgentVersion string `gorm:"column:client_user_agent_version" json:"client_user_agent_version"` // 客户端UA版本
	ClientIP               string `gorm:"column:client_ip" json:"client_ip"`                                 // 客户端IP
}

// NewClientInfo 返回请求客户端信息
func NewClientInfo(ctx *gin.Context, dtv *DTV) (res *ClientInfo, err error) {
	if dtv == nil {
		err = errors.New("dtv is nil")
		return
	}

	client := &ClientInfo{
		ClientDID:       dtv.Did,
		ClientType:      dtv.ClientType,
		ClientVersion:   dtv.ClientVersion,
		ClientUserAgent: ctx.GetHeader("User-Agent"),
		ClientIP:        ctx.ClientIP(),
	}

	ua, err := ParseClientUserAgent(client.ClientUserAgent)
	if err != nil || ua == nil {
		return
	}

	client.ClientUserAgentName = ua.UserAgent.Family
	client.ClientUserAgentVersion = ua.UserAgent.ToVersionString()
	client.ClientOS = ua.Os.Family
	client.ClientOSVersion = ua.Os.ToVersionString()
	client.ClientDevice = ua.Device.Family
	client.ClientDeviceBrand = ua.Device.Brand
	client.ClientDeviceModel = ua.Device.Model

	res = client
	return
}

func IsValidClientType(ct uint8) bool {
	return ct > 0 && ct < ClientTypeMax
}

func IsValidClientVersion(cv string) bool {
	return len(cv) > 0
}

func ParseClientUserAgent(ua string) (res *uaparser.Client, err error) {
	parser, err := uaparser.NewFromBytes(uaparser.DefinitionYaml)
	if err == nil {
		res = parser.Parse(ua)
	}
	return
}

// GetClientDID 返回 ClientDID 字段的值，如果 receiver 为 nil，则返回零值。
func (c *ClientInfo) GetClientDID() string {
	if c == nil {
		return ""
	}
	return c.ClientDID
}

// GetClientType 返回 ClientType 字段的值，如果 receiver 为 nil，则返回零值。
func (c *ClientInfo) GetClientType() uint8 {
	if c == nil {
		return 0
	}
	return c.ClientType
}

// GetClientVersion 返回 ClientVersion 字段的值，如果 receiver 为 nil，则返回零值。
func (c *ClientInfo) GetClientVersion() string {
	if c == nil {
		return ""
	}
	return c.ClientVersion
}

// GetClientDevice 返回 ClientDevice 字段的值，如果 receiver 为 nil，则返回零值。
func (c *ClientInfo) GetClientDevice() string {
	if c == nil {
		return ""
	}
	return c.ClientDevice
}

// GetClientDeviceBrand 返回 ClientDeviceBrand 字段的值，如果 receiver 为 nil，则返回零值。
func (c *ClientInfo) GetClientDeviceBrand() string {
	if c == nil {
		return ""
	}
	return c.ClientDeviceBrand
}

// GetClientDeviceModel 返回 ClientDeviceModel 字段的值，如果 receiver 为 nil，则返回零值。
func (c *ClientInfo) GetClientDeviceModel() string {
	if c == nil {
		return ""
	}
	return c.ClientDeviceModel
}

// GetClientOS 返回 ClientOS 字段的值，如果 receiver 为 nil，则返回零值。
func (c *ClientInfo) GetClientOS() string {
	if c == nil {
		return ""
	}
	return c.ClientOS
}

// GetClientOSVersion 返回 ClientOSVersion 字段的值，如果 receiver 为 nil，则返回零值。
func (c *ClientInfo) GetClientOSVersion() string {
	if c == nil {
		return ""
	}
	return c.ClientOSVersion
}

// GetClientUserAgent 返回 ClientUserAgent 字段的值，如果 receiver 为 nil，则返回零值。
func (c *ClientInfo) GetClientUserAgent() string {
	if c == nil {
		return ""
	}
	return c.ClientUserAgent
}

// GetClientUserAgentName 返回 ClientUserAgentName 字段的值，如果 receiver 为 nil，则返回零值。
func (c *ClientInfo) GetClientUserAgentName() string {
	if c == nil {
		return ""
	}
	return c.ClientUserAgentName
}

// GetClientUserAgentVersion 返回 ClientUserAgentVersion 字段的值，如果 receiver 为 nil，则返回零值。
func (c *ClientInfo) GetClientUserAgentVersion() string {
	if c == nil {
		return ""
	}
	return c.ClientUserAgentVersion
}

// GetClientIP 返回 ClientIP 字段的值，如果 receiver 为 nil，则返回零值。
func (c *ClientInfo) GetClientIP() string {
	if c == nil {
		return ""
	}
	return c.ClientIP
}

func (c *ClientInfo) String() string {
	if c != nil {
		if bytes, err := json.Marshal(c); err == nil {
			return string(bytes)
		}
	}
	return ""
}

func (c *ClientInfo) Md5() string {
	if c != nil {
		return crypto.Md5(c.String())
	}
	return ""
}
