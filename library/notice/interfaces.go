package notice

import "io"

type NProgressor interface {
	Notifier
	Progressor
}

type Notifier interface {
	Send(message interface{}, statusCode int) error
	Success(message interface{}) error
	Failure(message interface{}) error
}

type Progressor interface {
	Add(n int64) NProgressor
	Done(n int64) NProgressor
	AutoComplete(on bool) NProgressor
	Complete() NProgressor
	Reset()
	ProxyReader(r io.Reader) io.ReadCloser
	ProxyWriter(w io.Writer) io.WriteCloser
	Callback(total int64, exec func(callback func(strLen int)) error) error
	SetControl(IsExited) NProgressor
	Progress() *Progress
}

type UserMessageSystem interface {
	SetDebug(on bool) UserMessageSystem
	Debug() bool
	OnClose(fn ...func(user string)) UserMessageSystem
	OnOpen(fn ...func(user string)) UserMessageSystem
	Sendable(user string, types ...string) bool
	Send(user string, message *Message) error
	Recv(user string, clientID string) <-chan *Message
	CloseClient(user string, clientID string) bool
	IsOnline(user string) bool
	OnlineStatus(users ...string) map[string]bool
	OpenClient(user string) (oUser IOnlineUser, clientID string)
	OpenClientWithID(user string, clientID string) (oUser IOnlineUser)
	CloseMessage(user string, types ...string)
	OpenMessage(user string, types ...string)
	Clear()
	Count() int
	UserList(limit int) []string
	MakeMessageGetter(username string, messageTypes ...string) (func(), <-chan *Message, error)
	MakeMessageGetterWithClientID(username string, clientID string, messageTypes ...string) (func(), <-chan *Message)
}

type IOnlineUser interface {
	GetUser() string
	HasMessageType(messageTypes ...string) bool
	Send(message *Message, openDebug ...bool) error
	Recv(clientID string) <-chan *Message
	ClearMessage()
	ClearMessageType(types ...string)
	OpenMessageType(types ...string)
	CountType() int
	CountClient() int
	CloseClient(clientID string) int
	OpenClient(clientID string)
}

type IOnlineUsers interface {
	GetOk(user string, noLock ...bool) (IOnlineUser, bool)
	OnlineStatus(users ...string) map[string]bool
	IsOnline(user string) bool
	Set(user string, oUser IOnlineUser)
	Delete(user string)
	Clear()
	Count() int
	UserList(limit int) []string
}

type NoticeMessager interface {
	Size() int
	Delete(clientID string) int
	Clear()
	Add(clientID string)
	Send(message *Message) error
	Recv(clientID string) <-chan *Message
}

type NoticeTyper interface {
	Has(types ...string) bool
	Size() int
	Clear(types ...string)
	Open(types ...string)
}
