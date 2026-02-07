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

type Notice struct {
	types    NoticeTyper
	messages NoticeMessager
}

func (a *Notice) CountClient() int {
	return a.messages.Size()
}

func (a *Notice) CloseClient(clientID string) int {
	return a.messages.Delete(clientID)
}

func (a *Notice) CloseAllClient() {
	a.messages.Clear()
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
	msg.ClientID = AnyClientID
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
