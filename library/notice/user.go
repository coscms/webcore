package notice

import (
	"fmt"
	"sync"

	"github.com/coscms/webcore/library/msgbox"
)

var _ IOnlineUser = (*OnlineUser)(nil)

type OnlineUser struct {
	User string
	*Notice
}

func (oUser *OnlineUser) GetUser() string {
	return oUser.User
}

func (oUser *OnlineUser) HasMessageType(messageTypes ...string) bool {
	return oUser.Notice.types.Has(messageTypes...)
}

func (oUser *OnlineUser) Send(message *Message, openDebug ...bool) error {
	if !oUser.Notice.types.Has(message.Type) {
		message.Failure()
		Stdout(message)
		return ErrMsgTypeNotAccept
	}
	var debug bool
	if len(openDebug) > 0 {
		debug = openDebug[0]
	}
	if debug {
		msgbox.Debug(`[NOTICE]`, `[Send][MessageTo]: `+oUser.User)
	}
	err := oUser.Notice.messages.Send(message)
	if err != nil {
		message.Failure()
		if debug {
			msgbox.Debug(`[NOTICE]`, `[Send][MessageTo]: `+oUser.User+` [NotFoundClientID]: `+fmt.Sprint(message.ClientID))
		}
	}
	return err
}

func (oUser *OnlineUser) Recv(clientID string) <-chan *Message {
	return oUser.Notice.messages.Recv(clientID)
}

func (oUser *OnlineUser) ClearMessage() {
	oUser.Notice.messages.Clear()
}

func (oUser *OnlineUser) ClearMessageType(types ...string) {
	oUser.Notice.types.Clear(types...)
}

func (oUser *OnlineUser) OpenMessageType(types ...string) {
	oUser.Notice.types.Open(types...)
}

func (oUser *OnlineUser) CountType() int {
	return oUser.Notice.types.Size()
}

func NewOnlineUser(user string) *OnlineUser {
	return &OnlineUser{
		User:   user,
		Notice: NewNotice(),
	}
}

func NewOnlineUsers() *OnlineUsers {
	return &OnlineUsers{
		user: map[string]IOnlineUser{},
	}
}

var _ IOnlineUsers = (*OnlineUsers)(nil)

type OnlineUsers struct {
	lock sync.RWMutex
	user map[string]IOnlineUser //key: user
}

func (o *OnlineUsers) GetOk(user string, noLock ...bool) (IOnlineUser, bool) {
	if len(noLock) > 0 && noLock[0] {
		oUser, exists := o.user[user]
		return oUser, exists
	}
	o.lock.RLock()
	oUser, exists := o.user[user]
	o.lock.RUnlock()
	return oUser, exists
}

func (o *OnlineUsers) OnlineStatus(users ...string) map[string]bool {
	r := map[string]bool{}
	o.lock.RLock()
	for _, user := range users {
		_, r[user] = o.user[user]
	}
	o.lock.RUnlock()
	return r
}

func (o *OnlineUsers) IsOnline(user string) bool {
	o.lock.RLock()
	_, ok := o.user[user]
	o.lock.RUnlock()
	return ok
}

func (o *OnlineUsers) Set(user string, oUser IOnlineUser) {
	o.lock.Lock()
	o.user[user] = oUser
	o.lock.Unlock()
}

func (o *OnlineUsers) Delete(user string) {
	o.lock.Lock()
	delete(o.user, user)
	o.lock.Unlock()
}

func (o *OnlineUsers) Clear() {
	o.lock.Lock()
	o.user = map[string]IOnlineUser{}
	o.lock.Unlock()
}

// Count 在线用户人数
func (o *OnlineUsers) Count() int {
	o.lock.RLock()
	n := len(o.user)
	o.lock.RUnlock()
	return n
}

// UserList 在线用户列表 limit=-1时为获取全部用户
func (o *OnlineUsers) UserList(limit int) []string {
	if limit == 0 {
		return []string{}
	}
	o.lock.RLock()
	l := len(o.user)
	if limit > 0 {
		l = limit
	}
	r := make([]string, 0, l)
	for user := range o.user {
		r = append(r, user)
		if limit > 0 && len(r) >= limit {
			break
		}
	}
	o.lock.RUnlock()
	return r
}
