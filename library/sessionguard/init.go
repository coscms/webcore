package sessionguard

import (
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/config/extend"
)

func init() {
	extend.Register(`sessionGuard`, func() interface{} {
		return &Config{}
	})
}

var emptyConfig = &Config{}

func GetConfig() *Config {
	cfg, ok := common.ExtendConfig().Get(`sessionGuard`).(*Config)
	if ok {
		return cfg
	}
	return emptyConfig
}

type Config struct {
	SessionGuardConfig
	Frontend *SessionGuardConfig `json:"frontend,omitempty"`
}

type SessionGuardConfig struct {
	Expires         int64 `json:"expires"`
	IgnoreBrowserUA bool  `json:"ignoreBrowserUA"`
	IgnoreBrowserIP bool  `json:"ignoreBrowserIP"`
}
