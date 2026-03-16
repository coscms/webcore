package email

import (
	"crypto/tls"
	"time"

	"github.com/admpub/email"
	"github.com/admpub/log"
	"github.com/admpub/mail"
)

var (
	//QueueSize 每批容量
	QueueSize = 50
	//QueueWait 队列等待时间
	QueueWait = time.Second * 10
)

type queueItem struct {
	Email  *email.Email
	Config Config
}

func (q *queueItem) send1() error {
	log.Debug(`<SendMail> Using: send1`)
	if q.Config.SMTP.Secure == "SSL" || q.Config.SMTP.Secure == "TLS" {
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         q.Config.SMTP.Host,
		}
		return q.Email.SendWithTLS(q.Config.SMTP.Address(), q.Config.Auth, tlsconfig)
	}
	return q.Email.Send(q.Config.SMTP.Address(), q.Config.Auth)
}

func (q *queueItem) send2() error {
	log.Debug(`<SendMail> Using: send2`)
	return mail.SendMail(
		q.Config.Subject,
		string(q.Config.Content),
		q.Config.ToAddress,
		q.Config.From,
		q.Config.CcAddress,
		q.Config.SMTP,
		nil,
	)
}

func (q *queueItem) Send() (err error) {
	if q.Config.Timeout <= 0 || q.Email == nil {
		if q.Email == nil {
			return q.send2()
		}
		return q.send1()
	}
	done := make(chan bool)
	go func() {
		err = q.send1()
		done <- true
	}()
	t := time.NewTicker(time.Second * time.Duration(q.Config.Timeout))
	defer t.Stop()
	select {
	case <-done:
		return
	case <-t.C:
		log.Error("发送邮件超时，采用备用方案发送")
		close(done)
	}
	return q.send2()
}
