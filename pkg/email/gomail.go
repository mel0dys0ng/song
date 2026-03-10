package email

import (
	"context"
	"fmt"

	"github.com/mel0dys0ng/song/pkg/erlogs"
	"github.com/mel0dys0ng/song/pkg/vipers"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

// GomailConfig 邮件配置
type GomailConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// GomailClient 邮件客户端
type GomailClient struct {
	*GomailConfig
}

// NewGomailClient 创建一个新的邮件客户端
func NewGomailClient(ctx context.Context, key string) *GomailClient {
	var config *GomailConfig
	err := vipers.UnmarshalKey(key, &config)
	if err != nil {
		erlogs.Convert(err).Wrap("failed to unmarshal gomail config").
			AppendFields(zap.String("key", key)).
			PanicLog(ctx)
	}
	return &GomailClient{
		GomailConfig: config,
	}
}

// SendEmail 发送邮件
func (c *GomailClient) SendEmail(to []string, subject, body string, isHTML bool) (err error) {
	// 创建邮件消息
	m := gomail.NewMessage()

	m.SetHeader("From", c.Username)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)

	// 设置邮件内容类型
	if isHTML {
		m.SetBody("text/html", body) // HTML 内容
	} else {
		m.SetBody("text/plain", body) // 纯文本内容
	}

	// 创建 SMTP 客户端
	d := gomail.NewDialer(c.Host, c.Port, c.Username, c.Password)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return
}

// SendEmailWithAttachment 发送带附件的邮件
func (c *GomailClient) SendEmailWithAttachment(to []string, subject, body string, isHTML bool, attachmentPath string) (err error) {
	// 创建邮件消息
	m := gomail.NewMessage()

	m.SetHeader("From", c.Username)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)

	// 设置邮件内容类型
	if isHTML {
		m.SetBody("text/html", body) // HTML 内容
	} else {
		m.SetBody("text/plain", body) // 纯文本内容
	}

	// 添加附件
	if attachmentPath != "" {
		m.Attach(attachmentPath)
	}

	// 创建 SMTP 客户端
	d := gomail.NewDialer(c.Host, c.Port, c.Username, c.Password)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return
}
