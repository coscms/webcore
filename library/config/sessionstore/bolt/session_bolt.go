package bolt

import (
	"path/filepath"
	"reflect"
	"time"

	boltstore "github.com/coscms/session-boltstore"
	"github.com/coscms/webcore/library/config"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/session/engine"
	"github.com/webx-top/echo/middleware/session/engine/cookie"
	"github.com/webx-top/echo/param"
)

func init() {
	config.RegisterSessionStore(`bolt`, `BoltDB存储`, initSessionStoreBolt)
}

var sessionStoreBoltOptions *boltstore.BoltOptions

func initSessionStoreBolt(_ *config.Config, cookieOptions *cookie.CookieOptions, sessionConfig param.Store) (changed bool, err error) {
	boltOptions := &boltstore.BoltOptions{
		File:          sessionConfig.String(`savePath`),
		KeyPairs:      cookieOptions.KeyPairs,
		BucketName:    sessionConfig.String(`bucketName`),
		MaxLength:     sessionConfig.Int(`maxLength`),
		EmptyDataAge:  sessionConfig.Int(`emptyDataAge`),
		CheckInterval: time.Duration(sessionConfig.Int64(`checkInterval`)) * time.Second,
	}
	if len(boltOptions.BucketName) == 0 {
		boltOptions.BucketName = `sessions`
	}
	if len(boltOptions.File) == 0 {
		boltOptions.File = filepath.Join(echo.Wd(), `data`, `cache`, `sessions`, `bolt`)
	}
	if com.IsDir(boltOptions.File) {
		boltOptions.File = filepath.Join(boltOptions.File, `bolt`)
	}
	if sessionStoreBoltOptions == nil || !engine.Exists(`bolt`) || !reflect.DeepEqual(boltOptions, sessionStoreBoltOptions) {
		boltstore.RegWithOptions(boltOptions)
		sessionStoreBoltOptions = boltOptions
		changed = true
	}
	return
}
