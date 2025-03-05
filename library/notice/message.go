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

type Callback struct {
	Success func(*Message)
	Failure func(*Message)
}

type Message struct {
	ClientID string        `json:"client_id" xml:"client_id"`
	ID       interface{}   `json:"id" xml:"id"`
	Type     string        `json:"type" xml:"type"`
	Title    string        `json:"title" xml:"title"`
	Status   int           `json:"status" xml:"status"`
	Content  interface{}   `json:"content" xml:"content"`
	Mode     string        `json:"mode" xml:"mode"` //显示模式：notify/element/modal
	Progress *ProgressInfo `json:"progress" xml:"progress"`
	callback *Callback
}

func (m *Message) Success() {
	if m.callback != nil && m.callback.Success != nil {
		m.callback.Success(m)
	}
}

func (m *Message) Failure() {
	if m.callback != nil && m.callback.Failure != nil {
		m.callback.Failure(m)
	}
}

func (m *Message) Release() {
	m.ClientID = ``
	m.ID = nil
	m.Type = ``
	m.Title = ``
	m.Status = 0
	m.Content = nil
	m.Mode = ``
	m.callback = nil
	if m.Progress != nil {
		releaseProgressInfo(m.Progress)
		m.Progress = nil
	}
	releaseMessage(m)
}

func (m *Message) SetCallback(cb *Callback) *Message {
	m.callback = cb
	return m
}

func (m *Message) SetType(t string) *Message {
	m.Type = t
	return m
}

func (m *Message) SetTitle(title string) *Message {
	m.Title = title
	return m
}

func (m *Message) SetID(id interface{}) *Message {
	m.ID = id
	return m
}

func (m *Message) SetClientID(clientID string) *Message {
	m.ClientID = clientID
	return m
}

func (m *Message) SetStatus(status int) *Message {
	m.Status = status
	return m
}

func (m *Message) SetContent(content interface{}) *Message {
	m.Content = content
	return m
}

func (m *Message) SetMode(mode string) *Message {
	m.Mode = mode
	return m
}

func (m *Message) SetProgress(progress *Progress) *Message {
	clonedProg := acquireProgressInfo()
	progress.CopyToInfo(clonedProg)
	m.Progress = clonedProg
	return m
}
