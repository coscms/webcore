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

package license

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/coscms/webcore/cmd/bootconfig"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/notice"
	"github.com/coscms/webcore/library/selfupdate"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

type ProductVersion struct {
	Version          string `comment:"版本号(格式1.0.1)" json:"version" xml:"version"`
	Type             string `comment:"版本类型(stable-稳定版;beta-公测版;alpha-内测版)" json:"type" xml:"type"`
	OsArch           string `comment:"支持的操作系统(多个用逗号分隔)，留空表示不限制" json:"os_arch" xml:"os_arch"`
	ReleasedAt       uint   `comment:"发布时间" json:"released_at" xml:"released_at"`
	ForceUpgrade     string `comment:"是否强行升级为此版本" json:"force_upgrade" xml:"force_upgrade"`
	Description      string `comment:"发布说明" json:"description" xml:"description"`
	Remark           string `comment:"备注" json:"remark" xml:"remark"`
	DownloadURL      string `comment:"下载网址" json:"download_url" xml:"download_url"`
	Sign             string `comment:"下载后验证签名(多个签名之间用逗号分隔)" json:"sign" xml:"sign"`
	DownloadURLOther string `comment:"备用下载网址" json:"download_url_other" xml:"download_url_other"`

	// local
	downloadedPath string
	extractedDir   string
	backupDir      string
	executable     string
	isNew          bool
	prog           notice.NProgressor
}

func (v *ProductVersion) SetDownloadedPath(downloadedPath string) {
	v.downloadedPath = downloadedPath
}

func (v *ProductVersion) SetExecutable(executable string) {
	v.executable = executable
}

func (v *ProductVersion) DownloadedPath() string {
	return v.downloadedPath
}

func (v *ProductVersion) Executable() string {
	return v.executable
}

func (v *ProductVersion) SetProgressor(prog notice.NProgressor) {
	prog.Reset()
	v.prog = prog
}

func (v *ProductVersion) IsNew() bool {
	return v.isNew
}

func (v *ProductVersion) clean(downloadDir string, newVersionDir string) {
	if len(newVersionDir) == 0 {
		return
	}
	dirEntries, _ := os.ReadDir(newVersionDir)
	downloadFolder := filepath.Base(downloadDir)
	for _, dirEntry := range dirEntries {
		if dirEntry.Name() == downloadFolder {
			continue
		}
		ppath := filepath.Join(newVersionDir, dirEntry.Name())
		os.RemoveAll(ppath)
		v.prog.Send(fmt.Sprintf(`clean up old files %q`, ppath), notice.StateSuccess)
	}
}

func (v *ProductVersion) Extract() error {
	if len(v.downloadedPath) == 0 {
		return fmt.Errorf(`failed to download: %s`, v.DownloadURL)
	}
	downloadDir := filepath.Dir(v.downloadedPath)
	v.backupDir = filepath.Join(downloadDir, `backup`)
	newVersionDir := filepath.Join(echo.Wd(), `data/cache/nging-new-version`)
	v.extractedDir = filepath.Join(newVersionDir, `latest`)
	v.clean(downloadDir, newVersionDir)
	com.MkdirAll(v.extractedDir, os.ModePerm)
	ddp := filepath.Join(v.extractedDir, `download_dir.txt`)
	err := os.WriteFile(ddp, com.Str2bytes(downloadDir), 0666)
	if err != nil {
		return fmt.Errorf(`%w: %s`, err, ddp)
	}
	v.prog.Send(fmt.Sprintf(`extract the file %q to %q`, v.downloadedPath, v.extractedDir), notice.StateSuccess)
	if strings.EqualFold(filepath.Ext(v.downloadedPath), `.zip`) {
		err = com.Unzip(v.downloadedPath, v.extractedDir)
	} else {
		_, err = com.UnTarGz(v.downloadedPath, v.extractedDir)
	}
	if err != nil {
		v.prog.Send(fmt.Sprintf(`failed to extract %q to %q`, v.downloadedPath, v.extractedDir), notice.StateFailure)
		return fmt.Errorf(`%w: %s`, err, v.downloadedPath)
	}
	subDir := strings.SplitN(filepath.Base(v.downloadedPath), `.`, 2)[0]
	_extractedDir := filepath.Join(v.extractedDir, subDir)
	if com.FileExists(_extractedDir) {
		v.extractedDir = _extractedDir
	}
	if len(v.executable) == 0 {
		if err = v.findExecutable(); err != nil {
			return err
		}
	}
	os.Chmod(v.executable, 0755)
	if len(v.Version) == 0 {
		err = v.getVersionInfo()
	}
	return err
}

func (v *ProductVersion) findExecutable() error {
	files, err := filepath.Glob(v.extractedDir + echo.FilePathSeparator + `*.sha256`)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		fileName := filepath.Base(v.downloadedPath)
		parts := strings.SplitN(fileName, `_`, 2)
		if len(parts) != 2 {
			parts = strings.SplitN(fileName, `.`, 2)
		}
		executable := filepath.Join(v.extractedDir, parts[0])
		if com.IsWindows {
			executable += `.exe`
		}
		if com.FileExists(executable) {
			v.executable = executable
			return err
		}
	} else {
		for _, sha256file := range files {
			executable := strings.TrimSuffix(sha256file, `.sha256`)
			if com.FileExists(executable) {
				v.executable = executable
			}
		}
	}
	if len(v.executable) == 0 {
		err = fmt.Errorf(`no executable found in folder %s`, v.extractedDir)
	}
	return err
}

func (v *ProductVersion) getVersionInfo() error {
	b, err := exec.Command(v.executable, `version`).CombinedOutput()
	if err != nil {
		return err
	}
	s := string(b)
	var ver *config.VersionInfo
	ver, err = parseVersionInfo(s)
	if err != nil {
		return err
	}
	if ver.Name != config.Version.Name {
		err = fmt.Errorf(`software name mismatch (new != old): %s != %s`, ver.Name, config.Version.Name)
		return err
	}
	if ver.Package != Package() {
		err = fmt.Errorf(`software package mismatch (new != old): %s != %s`, ver.Package, Package())
		return err
	}
	v.Version = ver.Number
	v.Type = ver.Label
	return err
}

func parseVersionInfo(s string) (*config.VersionInfo, error) {
	var err error
	// $ ./nging version
	// Nging v5.3.3 licensed(vanguard)
	// Schema: v7.7001
	// Build: 20250424162337
	s = strings.TrimSpace(s)
	rows := strings.Split(s, "\n")
	if len(rows) < 3 {
		err = fmt.Errorf(`failed to query version information from command line output: %s`, "\n"+s)
		return nil, err
	}
	rows = rows[len(rows)-3:]
	ngingVer := strings.TrimSpace(rows[0])
	parts := strings.Split(ngingVer, ` `)
	if len(parts) != 3 {
		err = fmt.Errorf(`failed to query version information from command line output: %s`, "\n"+ngingVer)
		return nil, err
	}
	ver := &config.VersionInfo{
		Name: parts[0],
	}
	if !strings.HasPrefix(parts[1], `v`) {
		err = fmt.Errorf(`failed to query version information from command line output: %s`, "\n"+ngingVer)
		return ver, err
	}
	ver.Number = strings.TrimPrefix(parts[1], `v`)
	vParts := strings.SplitN(ver.Number, `-`, 2)
	if len(vParts) == 2 {
		ver.Number = vParts[0]
		ver.Label = vParts[1]
	}
	parts2 := strings.SplitN(parts[2], `(`, 2)
	if len(parts2) == 2 {
		packageName = strings.TrimSuffix(parts2[1], `)`)
		ver.Package = packageName
	}

	parts = strings.SplitN(rows[1], ` `, 2)
	if len(parts) == 2 && parts[0] == `Schema:` {
		ver.DBSchema = com.Float64(strings.TrimPrefix(parts[1], `v`))
	}

	parts = strings.SplitN(rows[2], ` `, 2)
	if len(parts) == 2 && parts[0] == `Build:` {
		ver.BuildTime = parts[1]
	}
	return ver, err
}

func (v *ProductVersion) Upgrade(ctx echo.Context, ngingDir string, restartMode ...string) error {
	if len(v.extractedDir) == 0 {
		if err := v.Extract(); err != nil {
			return err
		}
	}
	executable := filepath.Base(v.executable)
	backupDir := v.backupDir
	com.MkdirAll(backupDir, 0755)
	v.prog.Send(fmt.Sprintf(`copy the files from %q to %q`, v.extractedDir, ngingDir), notice.StateSuccess)
	var backupFiles []string
	var extension string
	if com.IsWindows {
		extension = `.exe`
	}

	// 复制文件到安装位置
	err := com.CopyDir(v.extractedDir, ngingDir, func(filePath string) bool {
		// 复制前的过滤操作：返回 true 代表是需要被过滤(跳过)的文件
		// 此处顺带进行文件备份
		// fmt.Println(filePath)
		if filePath == `download_dir.txt` {
			return true // 跳过标记文件
		}
		oldFile := filepath.Join(ngingDir, filePath)
		if fi, err := os.Stat(oldFile); err == nil {
			if fi.IsDir() {
				dir := filepath.Join(backupDir, filePath)
				if !com.FileExists(dir) {
					com.MkdirAll(dir, fi.Mode())
				}
			} else {
				if filePath == `startup`+extension {
					return true // 跳过此处复制。如果需要升级 startup，需要手动升级
				}
				backupFile := filepath.Join(backupDir, filePath)
				err = com.Copy(oldFile, backupFile) // 备份文件
				if err != nil {
					v.prog.Send(fmt.Sprintf(`failed to back up file %q: %v`, backupFile, err), notice.StateFailure)
				} else {
					v.prog.Send(fmt.Sprintf(`back up file %q to %q`, oldFile, backupFile), notice.StateSuccess)
				}
				backupFiles = append(backupFiles, backupFile) // 记录下来备份的文件
			}
		}
		if executable == filePath {
			return true // 跳过此处复制，采用单独的替换逻辑来处理
		}
		return false // 需要复制(不跳过)
	})
	if err != nil {
		return err
	}

	// 备份失败之后的还原处理
	restore := func() {
		for _, backupFile := range backupFiles {
			targetFile := strings.TrimPrefix(backupFile, backupDir)
			targetFile = filepath.Join(ngingDir, targetFile)
			err := com.Copy(backupFile, targetFile)
			if err != nil {
				v.prog.Send(fmt.Sprintf(`failed to restore file %q: %v`, targetFile, err), notice.StateFailure)
			} else {
				v.prog.Send(fmt.Sprintf(`restore file %q to %q`, backupFile, targetFile), notice.StateSuccess)
			}
		}
	}

	if len(v.executable) > 0 {
		var fp *os.File
		fp, err = os.Open(v.executable)
		if err != nil {
			return fmt.Errorf(`%w: %v`, err, v.executable)
		}
		defer fp.Close()
		targetExecutable := filepath.Join(ngingDir, executable)
		v.prog.Send(fmt.Sprintf(`update file %q`, targetExecutable), notice.StateSuccess)
		err = selfupdate.Update(fp, targetExecutable)
		if err != nil {
			v.prog.Send(fmt.Sprintf(`failed to update file %q: %v`, targetExecutable, err), notice.StateFailure)
			err = fmt.Errorf(`%w: %v`, err, targetExecutable)
			return err
		}
		v.prog.Send(fmt.Sprintf(`restart file %q`, targetExecutable), notice.StateSuccess)
		err = selfupdate.Restart(func(err error) {
			if err == nil {
				// v.prog.Send(`exit the current process`, notice.StateSuccess)
				// os.Exit(0)
				return
			}
			v.prog.Send(fmt.Sprintf(`failed to restart file %q: %v`, targetExecutable, err), notice.StateFailure)
			v.prog.Send(`start restoring files`, notice.StateSuccess)
			restore()
		}, targetExecutable, restartMode...)
		if err == nil {
			v.prog.Send(`successfully upgrade `+bootconfig.SoftwareName+` to version `+v.Version, notice.StateSuccess)
		}
	}
	v.prog.Complete()
	return err
}
