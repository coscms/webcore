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
	"net/smtp"

	"github.com/admpub/mail"
	"github.com/coscms/webcore/library/notice"
)

type Config struct {
	ID         uint64 //RequestID
	Engine     string
	SMTP       *mail.SMTPConfig `json:"-" xml:"-"`
	From       string
	ToAddress  string
	ToUsername string
	Subject    string
	Content    []byte
	CcAddress  []string
	Auth       smtp.Auth `json:"-" xml:"-"`
	Timeout    int64
	Noticer    notice.Noticer       `json:"-" xml:"-"`
	Callback   func(*Config, error) `json:"-" xml:"-"`
}

func (c *Config) Send() error {
	return SendMail(c)
}

func (c *Config) SetCallback(f func(*Config, error)) *Config {
	c.Callback = f
	return c
}

func (c *Config) FireCallback(cfg *Config, err error) {
	if c.Callback == nil {
		return
	}
	c.Callback(cfg, err)
}
