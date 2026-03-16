package email

import (
	"net/smtp"

	"github.com/admpub/mail"
	"github.com/coscms/webcore/library/notice"
)

type Config struct {
	ID         uint64 //RequestID
	Engine     string
	SMTP       *mail.SMTPConfig
	From       string
	ToAddress  string
	ToUsername string
	Subject    string
	Content    []byte
	CcAddress  []string
	Auth       smtp.Auth
	Timeout    int64
	Noticer    notice.Noticer
}

func (c *Config) Send() error {
	return SendMail(c)
}
