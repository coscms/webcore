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
	"database/sql/driver"
	"errors"
	"fmt"
	stdLog "log"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/admpub/confl"
	"github.com/admpub/log"
	"github.com/admpub/mysql-schema-sync/sync"
	"github.com/coscms/webcore/cmd/bootconfig"
	"github.com/coscms/webcore/library/config/subconfig/sdb"
	"github.com/coscms/webcore/library/config/subconfig/ssystem"
	"github.com/coscms/webcore/library/nretry"
	"github.com/coscms/webcore/library/nsql"
	"github.com/webx-top/com"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/lib/sqlbuilder"
	"github.com/webx-top/db/mysql"
	"github.com/webx-top/echo"
)

func MustGetConfig() *Config {
	if FromFile() == nil {
		FromCLI().ParseConfig()
	}
	return FromFile()
}

func InitConfig() (*Config, error) {
	configFiles := []string{FromCLI().Conf}
	if !strings.HasSuffix(configFiles[0], `.conf`) {
		configFiles = append(
			configFiles,
			strings.TrimSuffix(configFiles[0], filepath.Ext(configFiles[0]))+`.conf`,
		)
	}
	configFiles = append(configFiles,
		filepath.Join(FromCLI().Confd, `config.sample.conf`),
		filepath.Join(FromCLI().Confd, `config.yaml.sample`),
	)
	var (
		configFile      string
		err             error
		temporaryConfig = NewConfig()
	)
	temporaryConfig.settings.Debug = bootconfig.Develop
	for key, conf := range configFiles {
		if !filepath.IsAbs(conf) {
			conf = filepath.Join(echo.Wd(), conf)
			configFiles[key] = conf
			if key == 0 {
				FromCLI().Conf = conf
			}
		}
		_, err = os.Stat(conf)
		if err == nil {
			configFile = conf
			break
		}
		if !os.IsNotExist(err) {
			return temporaryConfig, err
		}
	}
	if err != nil {
		return temporaryConfig, err
	}
	_, err = confl.DecodeFile(configFile, temporaryConfig)
	if err != nil {
		return temporaryConfig, err
	}
	temporaryConfig.SetDefaults()

	return temporaryConfig, nil
}

func GenerateDefaultConfig() error {
	b, err := confl.Marshal(FromFile())
	if err != nil {
		return err
	}
	err = os.WriteFile(FromCLI().Conf, b, os.ModePerm)
	return err
}

func ParseConfig() error {
	conf, err := InitConfig()
	if err != nil {
		return err
	}
	InitSessionOptions(conf)
	if IsInstalled() {
		if FromFile() != nil {
			if !FromFile().connectedDB || !reflect.DeepEqual(conf.DB, FromFile().DB) {
				if err = conf.connectDB(); err != nil {
					return err
				}
			}
			err = FromFile().Reload(conf)
		} else {
			err = conf.connectDB()
		}
		if err != nil {
			return err
		}
	}
	conf.AsDefault()
	return err
}

var (
	DBConnecters = map[string]func(sdb.DB) (sqlbuilder.Database, error){
		`mysql`: ConnectMySQL,
	}
	DBInstallers = map[string]func(string) error{
		`mysql`: ExecMySQL,
	}
	DBCreaters = map[string]func(error, sdb.DB) error{
		`mysql`: CreaterMySQL,
	}
	DBUpgraders = map[string]func(string, *sync.Config, sdb.DB) (DBOperators, error){
		`mysql`: UpgradeMySQL,
	}
	DBEngines         = echo.NewKVData().Add(`mysql`, `MySQL`)
	ParseTimeDuration = ssystem.ParseTimeDuration
	ParseBytes        = ssystem.ParseBytes
)

type DBOperators struct {
	Source      sync.DBOperator
	Destination sync.DBOperator
}

func CreaterMySQL(err error, c sdb.DB) error {
	if strings.Contains(err.Error(), `Unknown database`) {
		dbName := c.Database
		c.Database = ``
		err2 := ConnectDB(c, 0, `default`)
		if err2 != nil {
			return err
		}
		charset := c.Charset()
		if len(charset) == 0 {
			charset = sdb.MySQLDefaultCharset
		}
		sqlStr := "CREATE DATABASE `" + dbName + "` COLLATE '" + charset + "_general_ci'"
		_, err = factory.NewParam().SetCollection(sqlStr).Exec()
		if err != nil {
			return err
		}
		c.Database = dbName
		err = ConnectDB(c, 0, `default`)
	}
	return err
}

func UpgradeMySQL(schema string, syncConfig *sync.Config, cfg sdb.DB) (DBOperators, error) {
	syncConfig.DestDSN = cfg.User + `:` + cfg.Password + `@(` + cfg.Host + `)/` + cfg.Database
	t := `?`
	for key, value := range cfg.Options {
		syncConfig.DestDSN += t + fmt.Sprintf("%s=%s", key, url.QueryEscape(value))
		t = `&`
	}
	syncConfig.SQLPreprocessor = func() func(string) string {
		charset := cfg.Charset()
		if len(charset) == 0 {
			charset = sdb.MySQLDefaultCharset
		}
		return func(sqlStr string) string {
			return nsql.ReplaceCharset(sqlStr, charset)
		}
	}()
	return DBOperators{Source: sync.NewMySchemaData(schema, `source`)}, nil
}

func ConnectMySQL(c sdb.DB) (sqlbuilder.Database, error) {
	settings := c.ToMySQL()
	return mysql.Open(settings)
}

func ExecMySQL(sqlStr string) error {
	_, err := factory.NewParam().SetCollection(sqlStr).Exec()
	if err != nil {
		stdLog.Println(err.Error(), `->SQL:`, sqlStr)
	}
	return err
}

func QueryTo(sqlStr string, result interface{}) (sqlbuilder.Iterator, error) {
	return factory.NewParam().SetRecv(result).SetCollection(sqlStr).QueryTo()
}

func ConnectDB(c sdb.DB, index int, name string) error {
	if index == 0 {
		factory.CloseAll()
		factory.SetDebug(c.Debug)
	}
	fn, ok := DBConnecters[c.Type]
	if !ok {
		return ErrUnknowDatabaseType
	}
	err := nretry.Retry(10, func() error {
		database, err := fn(c)
		if err != nil {
			if !errors.Is(err, driver.ErrBadConn) && !com.IsNetworkOrHostDown(err, false) {
				return nretry.NoRetry(fmt.Errorf(`failed to connect %v: %w`, c.Type, err))
			}
			return fmt.Errorf(`failed to connect %v: %w`, c.Type, err)
		}

		log.Debugf(`successfully connected to the database: %s`, c.Description())

		c.SetConn(database)
		cluster := factory.NewCluster().AddMaster(database)
		factory.SetCluster(index, cluster)
		if len(name) > 0 {
			factory.SetIndexName(index, name)
		}
		return err
	})
	return err
}

func MustOK(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

var CmdIsRunning = com.CmdIsRunning

func Table(table string) string {
	return FromFile().DB.Table(table)
}

func ToTable(m sqlbuilder.Name_) string {
	return FromFile().DB.ToTable(m)
}

func FixWd() error {
	executableFile := filepath.Base(os.Args[0])
	if strings.HasSuffix(executableFile, `.test.exe`) ||
		strings.HasSuffix(executableFile, `.test`) ||
		strings.HasPrefix(os.Args[0], os.TempDir()) {
		_, file, _, ok := runtime.Caller(0)
		if ok {
			parts := strings.SplitN(file, echo.FilePathSeparator+`vendor`+echo.FilePathSeparator, 2)
			if len(parts) == 2 {
				echo.SetWorkDir(parts[0])
			}
		}
		return nil
	}

	// from os.Getwd()
	if filepath.IsAbs(echo.Wd()) {
		executable := filepath.Join(echo.Wd(), executableFile)
		if com.IsFile(executable) {
			return nil
		}
	}

	// from os.Args[0]
	dir := filepath.Dir(os.Args[0])
	absDir, err := filepath.Abs(dir)
	if err != nil {
		log.Fatalf(`failed to filepath.Abs(%q): %v`, dir, err.Error())
	}
	echo.SetWorkDir(absDir)
	executable := filepath.Join(echo.Wd(), executableFile)
	if com.IsFile(executable) {
		return nil
	}

	// from os.Executable()
	executable, err = os.Executable()
	if err != nil {
		log.Fatalf(`failed to os.Executable(): %v`, err.Error())
	}
	_executable := executable
	executable, err = filepath.EvalSymlinks(_executable)
	if err != nil {
		log.Fatalf(`failed to filepath.EvalSymlinks(%q): %v`, _executable, err.Error())
	}
	dir = filepath.Dir(executable)
	absDir, err = filepath.Abs(dir)
	if err != nil {
		log.Fatalf(`failed to filepath.Abs(%q): %v`, dir, err.Error())
	}
	echo.SetWorkDir(absDir)
	return nil
}
