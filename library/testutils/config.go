package testutils

import (
	"os"
	"path/filepath"

	"github.com/admpub/log"
	"github.com/coscms/webcore/library/config"
	"github.com/webx-top/com"
)

func InitConfig() {
	log.Sync()
	config.FromCLI().Conf = os.Getenv(`COSCMS_CONF`)
	if len(config.FromCLI().Conf) == 0 || !com.FileExists(config.FromCLI().Conf) {
		config.FromCLI().Conf = filepath.Join(os.Getenv("GOPATH"), `src`, `github.com/admpub/nging/config/config.yaml`)
	}
	if err := config.ParseConfig(); err != nil {
		panic(err)
	}
	config.FromFile().SetDebug(true)
	config.FromFile().ConnectedDB()
}
