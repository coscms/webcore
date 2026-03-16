/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package email

import (
	"context"
	"strings"
	"time"

	"github.com/admpub/email"
	"github.com/admpub/log"
	"github.com/admpub/mail"
	"github.com/admpub/once"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/notice"
)

var (
	sendCh                 chan *queueItem
	smtpClient             *mail.SMTPClient
	smtpClientOnce         once.Once
	ctxGlobal              context.Context
	ctxCancel              context.CancelFunc
	onceInitial, onceReset = once.OnceValue(func() chan *queueItem {
		return initial(context.Background())
	})
)

func SMTPClient(conf *mail.SMTPConfig) *mail.SMTPClient {
	smtpClientOnce.Do(func() {
		c := mail.NewSMTPClient(conf)
		smtpClient = &c
	})
	return smtpClient
}

func Restart() {
	if ctxCancel != nil {
		ctxCancel()
	}
	onceReset()
}

func initial(parent context.Context, queueSizes ...int) chan *queueItem {
	var queueSize int
	if len(queueSizes) > 0 {
		queueSize = queueSizes[0]
	}
	if sendCh != nil {
		close(sendCh)
	}
	if queueSize <= 0 {
		queueSize = config.FromFile().Settings().Email.QueueSize
		if queueSize <= 0 {
			queueSize = QueueSize
		}
	}
	ctxGlobal, ctxCancel = context.WithCancel(parent)
	sendCh = make(chan *queueItem, queueSize)
	go func() {
		defer ctxCancel()
		for {
			select {
			case m, ok := <-sendCh:
				if !ok {
					return
				}
				noticer := m.Config.Noticer
				if noticer == nil {
					noticer = notice.DefaultNoticer
				}
				noticer("<SendMail> Sending: "+m.Config.ToAddress, 1)
				err := m.Send()
				if err != nil {
					noticer("<SendMail> Result: "+m.Config.ToAddress+" Error: "+err.Error(), 0)
				} else {
					noticer("<SendMail> Result: "+m.Config.ToAddress+" [OK]", 1)
				}
				for _, callback := range Callbacks {
					callback(&m.Config, err)
				}
			case <-ctxGlobal.Done():
				log.Debugf(`<SendMail> %v`, context.Canceled)
			}
		}
	}()
	return sendCh
}

func SendMail(conf *Config) error {
	sendCh := onceInitial()
	if conf.SMTP == nil {
		return ErrSMTPNoSet
	}
	if len(conf.SMTP.Host) == 0 || len(conf.SMTP.Username) == 0 {
		return ErrSMTPNoSet
	}
	if len(conf.From) == 0 {
		return ErrSenderNoSet
	}
	if len(conf.ToAddress) == 0 {
		return ErrRecipientNoSet
	}
	if conf.Auth == nil {
		conf.Auth = conf.SMTP.Auth()
	}
	var mail *email.Email
	if conf.Engine == `email` || conf.Engine == `send1` {
		mail = email.NewEmail()
		mail.From = conf.From
		if len(mail.From) == 0 {
			mail.From = conf.SMTP.Username
			if !strings.Contains(mail.From, `@`) {
				mail.From += `@` + conf.SMTP.Host
			}
		}
		mail.To = []string{conf.ToAddress}
		mail.Subject = conf.Subject
		mail.HTML = conf.Content
		if len(conf.CcAddress) > 0 {
			mail.Cc = conf.CcAddress
		}
	}
	item := &queueItem{Email: mail, Config: *conf}
	t := time.NewTicker(QueueWait)
	defer t.Stop()
	select {
	case sendCh <- item:
		return nil
	case <-t.C:
		return ErrSendChannelTimeout
	}
}
