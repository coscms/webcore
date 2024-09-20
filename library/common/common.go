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

package common

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/admpub/decimal"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/tplfunc"
)

// Ok 操作成功
func Ok(v string) Successor {
	return NewOk(v)
}

// Err 获取错误信息
func Err(ctx echo.Context, err error) (ret interface{}) {
	if err != nil {
		return ProcessError(ctx, err)
	}
	flash := ctx.Flash()
	if flash != nil {
		if errMsg, ok := flash.(string); ok {
			ret = errors.New(errMsg)
		} else {
			ret = flash
		}
	} else {
		ret = GetEncMsg(ctx)
	}
	return
}

func WithEncErrMsg(msg string) string {
	return `encMsg=` + EncErrMsg(msg)
}

func WithEncOkMsg(msg string) string {
	return `encMsg=` + EncOkMsg(msg)
}

func EncErrMsg(msg string) string {
	return CryptorConfig().Encode(`err:` + msg)
}

func EncOkMsg(msg string) string {
	return CryptorConfig().Encode(`ok:` + msg)
}

func GetEncMsg(ctx echo.Context) (ret interface{}) {
	encMsg := ctx.Query(`encMsg`)
	if len(encMsg) == 0 {
		encMsg = ctx.GetCookie(`EncMsg`)
		if len(encMsg) > 0 {
			ctx.Cookie().Set(`EncMsg`, ``, -1)
		}
	}
	if len(encMsg) == 0 {
		return
	}
	msg := CryptorConfig().Decode(encMsg)
	if len(msg) == 0 {
		return
	}
	parts := strings.SplitN(msg, `:`, 2)
	if len(parts) == 2 {
		switch parts[0] {
		case `ok`:
			ret = Ok(parts[1])
		default:
			ret = errors.New(parts[1])
		}
	}
	return
}

// SendOk 记录成功信息
func SendOk(ctx echo.Context, msg string, storeInCookie ...bool) {
	if ctx.IsAjax() || ctx.Format() != echo.ContentTypeHTML {
		ctx.Data().SetInfo(msg, 1)
		return
	}
	if len(storeInCookie) > 0 && storeInCookie[0] {
		ctx.Cookie().Set(`EncMsg`, CryptorConfig().Encode(`ok:`+msg))
		return
	}
	ctx.Session().AddFlash(Ok(msg))
}

// SendFail 记录失败信息
func SendFail(ctx echo.Context, msg string, storeInCookie ...bool) {
	if ctx.IsAjax() || ctx.Format() != echo.ContentTypeHTML {
		ctx.Data().SetInfo(msg, 0)
		return
	}
	if len(storeInCookie) > 0 && storeInCookie[0] {
		ctx.Cookie().Set(`EncMsg`, CryptorConfig().Encode(`err:`+msg))
		return
	}
	ctx.Session().AddFlash(msg)
}

// SendErr 记录错误信息 (SendFail的别名)
func SendErr(ctx echo.Context, err error) {
	err = ProcessError(ctx, err)
	SendFail(ctx, err.Error())
}

func AddLinks(ctx echo.Context, title string, url string, icon string, color ...string) *echo.KVList {
	links, ok := ctx.Get(`links`).(echo.KVList)
	if !ok {
		links = echo.KVList{}
	}
	options := []echo.KVOption{echo.KVOptHKV(`icon`, icon)}
	if len(color) > 0 {
		options = append(options, echo.KVOptHKV(`color`, color))
	}
	links.Add(title, url, options...)
	ctx.Set(`links`, links)
	return &links
}

type ConfigFromDB interface {
	ConfigFromDB() echo.H
}

var notWordRegexp = regexp.MustCompile(`[^\w]+`)

// LookPath 获取二进制可执行文件路径
func LookPath(bin string, otherPaths ...string) (string, error) {
	envVarName := `NGING_` + notWordRegexp.ReplaceAllString(strings.ToUpper(bin), `_`) + `_PATH`
	envVarValue := os.Getenv(envVarName)
	if len(envVarValue) > 0 {
		if com.IsFile(envVarValue) {
			return envVarValue, nil
		}
		envVarValue = filepath.Join(envVarValue, bin)
		if com.IsFile(envVarValue) {
			return envVarValue, nil
		}
	}
	findPath, err := exec.LookPath(bin)
	if err == nil {
		return findPath, err
	}
	if !errors.Is(err, exec.ErrNotFound) {
		return findPath, err
	}
	for _, binPath := range otherPaths {
		binPath = filepath.Join(binPath, bin)
		if com.IsFile(binPath) {
			return binPath, nil
		}
	}
	return findPath, err
}

func SeekLinesWithoutComments(r io.Reader) (string, error) {
	var content string
	err := com.SeekLines(r, WithoutCommentsLineParser(func(line string) error {
		content += line + "\n"
		return nil
	}))
	return content, err
}

func WithoutCommentsLineParser(exec func(string) error) func(string) error {
	var commentStarted bool
	return func(line string) error {
		lineClean := strings.TrimSpace(line)
		if len(lineClean) == 0 {
			return nil
		}
		if commentStarted {
			if strings.HasSuffix(lineClean, `*/`) {
				commentStarted = false
			}
			return nil
		}
		switch lineClean[0] {
		case '#':
			return nil
		case '/':
			if len(lineClean) > 1 {
				switch lineClean[1] {
				case '/':
					return nil
				case '*':
					commentStarted = true
					return nil
				}
			}
		}

		//content += line + "\n"
		return exec(line)
	}
}

func Float64Sum(numbers ...float64) float64 {
	if len(numbers) == 0 {
		return 0
	}
	d := decimal.NewFromFloat(numbers[0])
	if len(numbers) > 1 {
		var d2 decimal.Decimal
		for _, number := range numbers[1:] {
			d2 = decimal.NewFromFloat(number)
			d = d.Add(d2)
		}
	}
	number, _ := d.Float64()
	return number
}

func Float32Sum(numbers ...float32) float32 {
	if len(numbers) == 0 {
		return 0
	}
	d := decimal.NewFromFloat32(numbers[0])
	if len(numbers) > 1 {
		var d2 decimal.Decimal
		for _, number := range numbers[1:] {
			d2 = decimal.NewFromFloat32(number)
			d = d.Add(d2)
		}
	}
	number, _ := d.Float64()
	return float32(number)
}

func TemplateTags(keys ...string) echo.H {
	r := echo.H{}
	for _, key := range keys {
		r[key] = tplfunc.TemplateTag(key)
	}
	return r
}

func OSAbsPath(ppath string) string {
	if len(ppath) == 0 {
		return ppath
	}
	if !filepath.IsAbs(ppath) && !com.FileExists(ppath) {
		ppath = filepath.Join(echo.Wd(), ppath)
	}
	return ppath
}

func CopyFormDataFrom(ctx echo.Context, formData map[string][]string) {
	for key, vals := range formData {
		for idx, val := range vals {
			if idx == 0 {
				ctx.Request().Form().Set(key, val)
			} else {
				ctx.Request().Form().Add(key, val)
			}
		}
	}
}
