package sqlite

import (
	"path/filepath"
	"reflect"
	"strings"
	"time"

	sqlitestore "github.com/coscms/session-sqlitestore"
	_ "github.com/coscms/session-sqlitestore/driver"
	sqlstore "github.com/coscms/session-sqlstore"
	"github.com/coscms/webcore/library/config"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/session/engine"
	"github.com/webx-top/echo/middleware/session/engine/cookie"
	"github.com/webx-top/echo/param"
)

func init() {
	config.RegisterSessionStore(`sqlite`, `SQLite3存储`, initSessionStoreMySQL)
}

var sessionStoreSQLiteOptions *sqlitestore.Options

func initSessionStoreMySQL(c *config.Config, cookieOptions *cookie.CookieOptions, sessionConfig param.Store) (changed bool, err error) {
	sqliteOptions := &sqlitestore.Options{
		Path: sessionConfig.String(`path`, sessionConfig.String(`savePath`)),
		Options: sqlstore.Options{
			Table:         sessionConfig.String(`table`),
			KeyPairs:      cookieOptions.KeyPairs,
			MaxAge:        sessionConfig.Int(`maxAge`),
			EmptyDataAge:  sessionConfig.Int(`emptyDataAge`),
			MaxLength:     sessionConfig.Int(`maxLength`),
			CheckInterval: time.Duration(sessionConfig.Int64(`checkInterval`)) * time.Second,
			MaxReconnect:  sessionConfig.Int(`maxReconnect`),
		},
	}
	if len(sqliteOptions.Path) == 0 {
		sqliteOptions.Path = filepath.Join(echo.Wd(), `data`, `temp`, `sessions.db`)
	} else if com.IsDir(sqliteOptions.Path) {
		sqliteOptions.Path = filepath.Join(sqliteOptions.Path, `sessions.db`)
	} else if strings.HasSuffix(sqliteOptions.Path, echo.FilePathSeparator) {
		sqliteOptions.Path = filepath.Join(sqliteOptions.Path, `sessions.db`)
	}
	if len(sqliteOptions.Table) == 0 {
		sqliteOptions.Table = `sessions`
	}
	if sqliteOptions.MaxReconnect <= 0 {
		sqliteOptions.MaxReconnect = 30
	}
	if sessionStoreSQLiteOptions == nil || !engine.Exists(`sqlite`) || !reflect.DeepEqual(sqliteOptions, sessionStoreSQLiteOptions) {
		sqlitestore.RegWithOptions(sqliteOptions)
		sessionStoreSQLiteOptions = sqliteOptions
		changed = true
	}
	return
}
