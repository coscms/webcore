package notice

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/admpub/log"
	"github.com/coscms/webcore/library/msgbox"
)

type userNotices struct {
	users   IOnlineUsers //key: user
	debug   atomic.Bool
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

func (u *userNotices) SetDebug(on bool) UserMessageSystem {
	u.debug.Store(on)
	return u
}

func (u *userNotices) Debug() bool {
	return u.debug.Load()
}

func (u *userNotices) OnClose(fn ...func(user string)) UserMessageSystem {
	u.onClose = append(u.onClose, fn...)
	return u
}

func (u *userNotices) OnOpen(fn ...func(user string)) UserMessageSystem {
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
	remains := oUser.CloseClient(clientID)
	if u.Debug() {
		msgbox.Info(`[NOTICE]`, `[CloseClient][ClientID]: `+clientID)
	}
	if remains < 1 {
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
	return u.users.IsOnline(user)
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

func (u *userNotices) Count() int {
	return u.users.Count()
}

func (u *userNotices) UserList(limit int) []string {
	return u.users.UserList(limit)
}

func (u *userNotices) MakeMessageGetter(username string) (func(), <-chan *Message, error) {
	oUser, clientID := u.OpenClient(username)
	oUser.OpenMessageType(`clientID`)
	msg := NewMessage().SetMode(`-`).SetType(`clientID`).SetClientID(clientID)
	err := oUser.Send(msg)
	if err != nil {
		return nil, nil, err
	}
	msgChan := oUser.Recv(clientID)
	return func() {
		u.CloseClient(username, clientID)
	}, msgChan, err
}
