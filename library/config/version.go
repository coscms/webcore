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
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

// Version 版本信息
var Version = &VersionInfo{Name: `Nging`}

type VersionInfo struct {
	Name      string    //软件名称
	Number    string    //版本号 1.0.1
	Package   string    //套餐
	Label     string    //版本标签 beta/alpha/stable
	DBSchema  float64   //数据库表版本 例如：1.2
	BuildTime string    //构建时间
	BuildOS   string    //构建目标系统
	BuildArch string    //构建目标架构
	CommitID  string    //GIT提交ID
	Licensed  bool      //是否已授权
	Expired   time.Time //过期时间
}

func (v *VersionInfo) IsExpired() bool {
	if v.Expired.IsZero() {
		return false
	}
	return v.Expired.Before(time.Now())
}

func (v *VersionInfo) String() string {
	return v.Name + ` ` + v.VString()
}

func (v *VersionInfo) VString() string {
	var licenseTag string
	if v.Licensed {
		licenseTag = `licensed`
	} else {
		licenseTag = `unlicensed`
	}
	if len(v.Package) > 0 {
		licenseTag += `(` + v.Package + `)`
	}
	return `v` + v.VNumberString() + ` ` + licenseTag
}

func (v *VersionInfo) VNumberString() string {
	return JoinVersionLabel(v.Number, v.Label)
}

func JoinVersionLabel(number string, label string) string {
	version := number
	if len(label) > 0 && label != `stable` {
		version += `-` + label
	}
	return version
}

func (v *VersionInfo) IsNew(number string, label string, dump ...bool) bool {
	return v.IsNewWithBuildTime(number, label, ``, dump...)
}

func (v *VersionInfo) IsNewWithBuildTime(number string, label string, buildTime string, dump ...bool) bool {
	var hasNew bool
	newVersion := JoinVersionLabel(number, label)
	oldVersion := v.VNumberString()
	if len(dump) > 0 && dump[0] {
		d := echo.H{`newVersion`: newVersion, `oldVersion`: oldVersion}
		if len(buildTime) > 0 {
			d[`newBuildTime`] = buildTime
			d[`oldBuildTime`] = v.BuildTime
		}
		echo.Dump(d)
	}
	compared := com.VersionCompare(newVersion, oldVersion)
	switch compared {
	case com.VersionCompareGt:
		hasNew = true
	case com.VersionCompareEq:
		if len(buildTime) > 0 {
			hasNew = com.Uint64(v.BuildTime) < com.Uint64(buildTime)
		}
	}
	return hasNew
}
