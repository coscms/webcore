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
	"sync"

	"github.com/admpub/events"
	"github.com/webx-top/echo"
)

var (
	defaultUserNotices UserMessageSystem
	once               sync.Once
)

func Initialize() {
	defaultUserNotices = NewUserNotices(false, nil)
	echo.OnCallback(`nging.user.logout.success`, onLogout, `notice.closeMessage`)
}

func onLogout(e events.Event) error {
	username := e.Context.String(`username`)
	CloseMessage(username)
	return nil
}

func Default() UserMessageSystem {
	once.Do(Initialize)
	return defaultUserNotices
}

func SetDebug(on bool) {
	Default().SetDebug(on)
}

func OnClose(fn ...func(user string)) UserMessageSystem {
	return Default().OnClose(fn...)
}

func OnOpen(fn ...func(user string)) UserMessageSystem {
	return Default().OnOpen(fn...)
}

func Send(user string, message *Message) error {
	return Default().Send(user, message)
}

func Recv(user string, clientID string) <-chan *Message {
	return Default().Recv(user, clientID)
}

func CloseClient(user string, clientID string) bool {
	return Default().CloseClient(user, clientID)
}

func IsOnline(user string) bool {
	return Default().IsOnline(user)
}

func OnlineStatus(users ...string) map[string]bool {
	return Default().OnlineStatus(users...)
}

func OnlineUserCount() int {
	return Default().Count()
}

func OnlineUserList(limit int) []string {
	return Default().UserList(limit)
}

func OpenClient(user string) (oUser IOnlineUser, clientID string) {
	return Default().OpenClient(user)
}

func CloseMessage(user string, types ...string) {
	Default().CloseMessage(user, types...)
}

func OpenMessage(user string, types ...string) {
	Default().OpenMessage(user, types...)
}

func Clear() {
	Default().Clear()
}
