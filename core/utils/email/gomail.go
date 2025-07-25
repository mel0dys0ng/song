package email

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

// GomailClient 邮件客户端
type GomailClient struct {
	Host     string // SMTP 服务器地址
	Port     int    // SMTP 端口
	Username string // 发件人邮箱
	Password string // 发件人邮箱密码
}

// NewGomailClient 创建一个新的邮件客户端
func NewGomailClient(host string, port int, username, password string) *GomailClient {
	return &GomailClient{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
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
