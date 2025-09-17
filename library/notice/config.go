package notice

import (
	"context"
	"time"

	"github.com/webx-top/echo"
)

func SetUser(user string) func(*Config) {
	return func(cfg *Config) {
		cfg.User = user
	}
}

func SetClientID(clientID string) func(*Config) {
	return func(cfg *Config) {
		cfg.ClientID = clientID
	}
}

func SetType(typ string) func(*Config) {
	return func(cfg *Config) {
		cfg.Type = typ
	}
}

func SetID(id interface{}) func(*Config) {
	return func(cfg *Config) {
		cfg.ID = id
	}
}

func SetMode(mode string) func(*Config) {
	return func(cfg *Config) {
		cfg.Mode = mode
	}
}

func SetTimeout(timeout time.Duration) func(*Config) {
	return func(cfg *Config) {
		cfg.Timeout = timeout
	}
}

func SetIsExited(isExited IsExited) func(*Config) {
	return func(cfg *Config) {
		cfg.IsExited = isExited
	}
}

func SetAutoComplete(autoComplete bool) func(*Config) {
	return func(cfg *Config) {
		cfg.AutoComplete = autoComplete
	}
}

type Config struct {
	User         string
	Type         string // Topic
	ClientID     string
	ID           interface{}
	IsExited     IsExited
	Timeout      time.Duration
	Mode         string // element / notify
	AutoComplete bool   // 是否自动完成
	progress     *Progress
}

func NewConfig() *Config {
	return &Config{}
}

func NewConfigByContext(eCtx echo.Context, noticeType string, user string, opts ...func(*Config)) *Config {
	noticerConfig := &Config{
		User: user,
		Type: noticeType,
	}
	noticerConfig.ClientID = GetClientID(eCtx)
	if len(user) > 0 && len(noticerConfig.ClientID) > 0 {
		noticerConfig.ID = GetNoticeID(eCtx)
		noticerConfig.Mode = GetNoticeMode(eCtx)
	}
	for _, opt := range opts {
		opt(noticerConfig)
	}
	return noticerConfig
}

func (c *Config) SetUser(user string) *Config {
	c.User = user
	return c
}

func (c *Config) SetType(typ string) *Config {
	c.Type = typ
	return c
}

func (c *Config) SetClientID(clientID string) *Config {
	c.ClientID = clientID
	return c
}

func (c *Config) SetID(id interface{}) *Config {
	c.ID = id
	return c
}

func (c *Config) SetTimeout(t time.Duration) *Config {
	c.Timeout = t
	return c
}

func (c *Config) SetIsExited(isExited IsExited) *Config {
	c.IsExited = isExited
	return c
}

func (c *Config) SetMode(mode string) *Config {
	c.Mode = mode
	return c
}

func (c *Config) SetAutoComplete(autoComplete bool) *Config {
	c.AutoComplete = autoComplete
	return c
}

func (c *Config) Progress() *Progress {
	if c.progress != nil {
		return c.progress
	}
	c.progress = NewProgress(c.AutoComplete).SetControl(c.IsExited)
	return c.progress
}

func (c *Config) Noticer(ctx context.Context) Noticer {
	return NewNoticer(ctx, c)
}
