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

package notice

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/admpub/log"
	"github.com/coscms/webcore/library/msgbox"
)

type Notice struct {
	types    NoticeTyper
	messages NoticeMessager
}

func (a *Notice) CountClient() int {
	return a.messages.Size()
}

func (a *Notice) CloseClient(clientID string) {
	a.messages.Delete(clientID)
}

func (a *Notice) OpenClient(clientID string) {
	a.messages.Add(clientID)
}

func NewMessageWithValue(typ string, title string, content interface{}, status ...int) *Message {
	st := Succeed
	if len(status) > 0 {
		st = status[0]
	}
	msg := acquireMessage()
	msg.Type = typ
	msg.Title = title
	msg.Status = st
	msg.Content = content
	return msg
}

func NewMessage() *Message {
	return acquireMessage()
}

func NewNotice() *Notice {
	return &Notice{
		types:    newNoticeTypes(),
		messages: newNoticeMessages(),
	}
}

type userNotices struct {
	users   IOnlineUsers //key: user
	_debug  atomic.Bool
	onClose []func(user string)
	onOpen  []func(user string)
}

func NewUserNotices(debug bool, users IOnlineUsers) *userNotices {
	if users == nil {
		users = NewOnlineUsers()
	}
	u := &userNotices{
		users:   users,
		onClose: []func(user string){},
		onOpen:  []func(user string){},
	}
	u.SetDebug(debug)
	return u
}

func Stdout(message *Message) {
	if message.Status == Succeed {
		log.Okay(message.Content)
		//os.Stdout.WriteString(fmt.Sprint(message.Content))
	} else {
		log.Error(message.Content)
		//os.Stderr.WriteString(fmt.Sprint(message.Content))
	}
	message.Release()
}

func (u *userNotices) SetDebug(on bool) *userNotices {
	u._debug.Store(on)
	return u
}

func (u *userNotices) Debug() bool {
	return u._debug.Load()
}

func (u *userNotices) OnClose(fn ...func(user string)) *userNotices {
	u.onClose = append(u.onClose, fn...)
	return u
}

func (u *userNotices) OnOpen(fn ...func(user string)) *userNotices {
	u.onOpen = append(u.onOpen, fn...)
	return u
}

func (u *userNotices) Sendable(user string, types ...string) bool {
	oUser, exists := u.users.GetOk(user)
	if !exists {
		return false
	}
	return oUser.HasMessageType(types...)
}

func (u *userNotices) Send(user string, message *Message) error {
	debug := u.Debug()
	if debug {
		msgbox.Debug(`[NOTICE]`, `[Send][FindUser]: `+user)
	}
	oUser, exists := u.users.GetOk(user)
	if !exists {
		if debug {
			msgbox.Debug(`[NOTICE]`, `[Send][NotFoundUser]: `+user)
		}
		Stdout(message)
		return ErrUserNotOnline
	}
	if debug {
		msgbox.Debug(`[NOTICE]`, `[Send][CheckRecvType]: `+message.Type+` (for user: `+user+`)`)
	}
	return oUser.Send(message, debug)
}

func (u *userNotices) Recv(user string, clientID string) <-chan *Message {
	oUser, exists := u.users.GetOk(user)
	if !exists {
		oUser = NewOnlineUser(user)
		u.users.Set(user, oUser)
	}
	return oUser.Recv(clientID)
}

func (u *userNotices) CloseClient(user string, clientID string) bool {
	oUser, exists := u.users.GetOk(user)
	if !exists {
		return true
	}
	oUser.CloseClient(clientID)
	if u.Debug() {
		msgbox.Info(`[NOTICE]`, `[CloseClient][ClientID]: `+clientID)
	}
	if oUser.CountClient() < 1 {
		oUser.ClearMessage()
		u.users.Delete(user)
		for _, fn := range u.onClose {
			fn(user)
		}
		return true
	}
	return false
}

func (u *userNotices) IsOnline(user string) bool {
	return u.users.OnlineStatus(user)[user]
}

func (u *userNotices) OnlineStatus(users ...string) map[string]bool {
	return u.users.OnlineStatus(users...)
}

func (u *userNotices) OpenClient(user string) (oUser IOnlineUser, clientID string) {
	var exists bool
	oUser, exists = u.users.GetOk(user)
	if !exists {
		oUser = NewOnlineUser(user)
		u.users.Set(user, oUser)
		for _, fn := range u.onOpen {
			fn(user)
		}
	}
	clientID = fmt.Sprint(time.Now().UnixMilli())
	oUser.OpenClient(clientID)
	if u.Debug() {
		msgbox.Info(`[NOTICE]`, `[OpenClient][ClientID]: `+clientID)
	}
	return
}

func (u *userNotices) CloseMessage(user string, types ...string) {
	oUser, exists := u.users.GetOk(user)
	if !exists {
		return
	}
	oUser.ClearMessageType(types...)
}

func (u *userNotices) OpenMessage(user string, types ...string) {
	oUser, exists := u.users.GetOk(user)
	if !exists {
		oUser = NewOnlineUser(user)
		u.users.Set(user, oUser)
	}
	oUser.OpenMessageType(types...)
}

func (u *userNotices) Clear() {
	u.users.Clear()
}
