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

package config

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	stdSync "sync"
	"time"

	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/coscms/webcore/library/setup"
)

const (
	LockFileName = `installed.lock`
	ConfigName   = `ConfigFromFile`
	SettingName  = `ConfigFromDB`
)

var _ = FixWd()

var (
	Installed             sql.NullBool
	installedSchemaVer    float64
	installedPkgSchemaVer = map[string]float64{}
	installedTime         time.Time
	defaultConfig         *Config
	defaultConfigMu       stdSync.RWMutex
	defaultCLIConfig      *CLIConfig
	onceCLIConfig         stdSync.Once
	onceUpgrade           stdSync.Once
	sqlCollection         = NewSQLCollection().RegisterInstall(`nging`, setup.InstallSQL)

	// Errors
	ErrUnknowDatabaseType = errors.New(`unkown database type`)
)

func GetInstalledPkgSchemaVer(pkg string) float64 {
	v, ok := installedPkgSchemaVer[pkg]
	if !ok {
		return -1
	}
	return v
}

func initCLIConfig() {
	defaultCLIConfig = NewCLIConfig()
}

func FromCLI() *CLIConfig {
	onceCLIConfig.Do(initCLIConfig)
	return defaultCLIConfig
}

func FromFile() *Config {
	defaultConfigMu.RLock()
	v := defaultConfig
	defaultConfigMu.RUnlock()
	return v
}

func FromDB(group ...string) echo.H {
	return echo.GetStoreByKeys(SettingName, group...)
}

func GetSQLCollection() *SQLCollection {
	return sqlCollection
}

func RegisterInstallSQL(project string, installSQL string) {
	sqlCollection.RegisterInstall(project, installSQL)
}

func RegisterInsertSQL(project string, insertSQL string) {
	sqlCollection.RegisterInsert(project, insertSQL)
}

func RegisterPreupgradeSQL(project string, version, preupgradeSQL string) {
	sqlCollection.RegisterPreupgrade(project, version, preupgradeSQL)
}

func GetInsertSQLs() map[string][]string {
	return sqlCollection.Insert
}

func GetInstallSQLs() map[string][]string {
	return sqlCollection.Install
}

func GetPreupgradeSQLs() map[string]map[string][]string {
	return sqlCollection.Preupgrade
}

func genInstalledLockFileContent(now time.Time, verInfo *VersionInfo) (string, error) {
	content := now.Format(`2006-01-02 15:04:05`) + "\n" + fmt.Sprint(verInfo.DBSchema)
	jsonV, err := com.JSONEncodeToString(verInfo.PkgDBSchemas)
	if err != nil {
		return content, err
	}
	content += "\n" + jsonV
	return content, nil
}

func parseInstalledLockFileContent(content string) (time.Time, float64, map[string]float64, error) {
	content = strings.TrimSpace(content)
	var installedTime, installedSchemaVer, installedPkgSchemaVer string
	lines := strings.Split(content, "\n")
	com.SliceExtract(lines, &installedTime, &installedSchemaVer, &installedPkgSchemaVer)
	var t time.Time
	var schemaV float64
	var pkgSchemaV = map[string]float64{}
	var err error
	if len(installedSchemaVer) > 0 {
		schemaV, _ = strconv.ParseFloat(strings.TrimSpace(installedSchemaVer), 64)
	}
	if len(installedPkgSchemaVer) > 0 {
		err = com.JSONDecodeString(installedPkgSchemaVer, &pkgSchemaV)
	}
	if len(installedTime) > 0 {
		t, _ = time.Parse(`2006-01-02 15:04:05`, strings.TrimSpace(installedTime))
	}
	return t, schemaV, pkgSchemaV, err
}

func SetInstalled(lockFile string) error {
	now := time.Now()
	content, err := genInstalledLockFileContent(now, Version)
	if err != nil {
		return err
	}
	err = os.WriteFile(lockFile, []byte(content), os.ModePerm)
	if err != nil {
		return err
	}
	installedTime = now
	installedSchemaVer = Version.DBSchema
	Installed.Valid = true
	Installed.Bool = true
	return nil
}

func InstalledLockFile() string {
	for _, lockFile := range []string{
		filepath.Join(FromCLI().ConfDir(), LockFileName),
		filepath.Join(echo.Wd(), LockFileName),
	} {
		if info, err := os.Stat(lockFile); err == nil && !info.IsDir() {
			return lockFile
		}
	}
	return ``
}

func IsInstalled() bool {
	if Installed.Valid {
		return Installed.Bool
	}
	lockFile := InstalledLockFile()
	if len(lockFile) > 0 {
		if b, e := os.ReadFile(lockFile); e == nil {
			content := string(b)
			installedTime, installedSchemaVer, installedPkgSchemaVer, e = parseInstalledLockFileContent(content)
			if e != nil {
				log.Error(e)
			}
		}
		Installed.Valid = true
		Installed.Bool = true
	}
	return Installed.Bool
}

func GetSQLInstallFiles() ([]string, error) {
	confDIR := FromCLI().Confd
	sqlFile := filepath.Join(confDIR, `install.sql`)
	var sqlFiles []string
	if com.FileExists(sqlFile) {
		sqlFiles = append(sqlFiles, sqlFile)
	}
	matches, err := filepath.Glob(confDIR + echo.FilePathSeparator + `install.*.sql`)
	if len(matches) > 0 {
		sqlFiles = append(sqlFiles, matches...)
	}
	return sqlFiles, err
}

func GetPreupgradeSQLFiles() []string {
	confDIR := FromCLI().Confd
	sqlFiles := []string{}
	matches, _ := filepath.Glob(confDIR + echo.FilePathSeparator + `preupgrade.*.sql`)
	if len(matches) > 0 {
		sqlFiles = append(sqlFiles, matches...)
	}
	return sqlFiles
}

func GetSQLInsertFiles() []string {
	confDIR := FromCLI().Confd
	sqlFile := filepath.Join(confDIR, `insert.sql`)
	sqlFiles := []string{}
	if com.FileExists(sqlFile) {
		sqlFiles = append(sqlFiles, sqlFile)
	}
	matches, _ := filepath.Glob(confDIR + echo.FilePathSeparator + `insert.*.sql`)
	if len(matches) > 0 {
		sqlFiles = append(sqlFiles, matches...)
	}
	return sqlFiles
}
