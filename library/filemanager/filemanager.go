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

package filemanager

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/coscms/webcore/library/charset"
	uploadClient "github.com/webx-top/client/upload"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"
)

var (
	EncodedSep   = com.URLEncode(`/`)
	EncodedSlash = com.URLEncode(`\`)
	EncodedSepa  = com.URLEncode(echo.FilePathSeparator)
)

func New(root string, editableMaxSize int, ctx echo.Context) (*fileManager, error) {
	rootFS, err := os.OpenRoot(root)
	if err != nil {
		return nil, err
	}
	return &fileManager{
		Context:         ctx,
		Root:            rootFS,
		EditableMaxSize: editableMaxSize,
	}, err
}

type fileManager struct {
	echo.Context
	Root            *os.Root
	EditableMaxSize int
}

func (f *fileManager) Close() error {
	return f.Root.Close()
}

func (f *fileManager) RealPath(filePath string) string {
	absPath := f.Root.Name()
	if len(filePath) > 0 {
		filePath = filepath.Clean(`/` + filePath)
		absPath = filepath.Join(f.Root.Name(), filePath)
	}
	return absPath
}

func (f *fileManager) Edit(file string, content string, encoding string) (interface{}, error) {
	fi, err := f.Root.Stat(file)
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return nil, errors.New(f.T(`不能编辑文件夹`))
	}
	if f.EditableMaxSize > 0 && fi.Size() > int64(f.EditableMaxSize) {
		return nil, errors.New(f.T(`很抱歉，不支持编辑超过%v的文件`, com.FormatByte(f.EditableMaxSize, 2, true)))
	}
	encoding = strings.ToLower(encoding)
	isUTF8 := encoding == `` || encoding == `utf-8`
	if f.IsPost() {
		b := []byte(content)
		if !isUTF8 {
			b, err = charset.Convert(`utf-8`, encoding, b)
			if err != nil {
				return ``, err
			}
		}
		err = f.Root.WriteFile(file, b, fi.Mode())
		return nil, err
	}
	b, err := f.Root.ReadFile(file)
	if err == nil && !isUTF8 {
		b, err = charset.Convert(encoding, `utf-8`, b)
	}
	return string(b), err
}

func (f *fileManager) CreateFile(file string, content string, encoding string) error {
	_, err := f.Root.Stat(file)
	if err == nil {
		return f.NewError(code.DataAlreadyExists, `新建文件失败，文件已经存在`)
	}
	if !os.IsNotExist(err) {
		return err
	}
	encoding = strings.ToLower(encoding)
	isUTF8 := encoding == `` || encoding == `utf-8`
	b := []byte(content)
	if !isUTF8 {
		b, err = charset.Convert(`utf-8`, encoding, b)
		if err != nil {
			return err
		}
	}
	err = f.Root.WriteFile(file, b, 0664)
	return err
}

func (f *fileManager) Remove(file string) error {
	fi, err := f.Root.Stat(file)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return f.Root.RemoveAll(file)
	}
	return f.Root.Remove(file)
}

func (f *fileManager) Mkdir(dir string, mode os.FileMode) error {
	return f.Root.MkdirAll(dir, mode)
}

func (f *fileManager) Rename(oldName string, newName string) (err error) {
	if len(newName) > 0 {
		err = f.Root.Rename(oldName, newName)
	} else {
		err = errors.New(f.T(`请输入有效的文件名称`))
	}
	return
}

func (f *fileManager) Chmod(file string, owner Perm, group Perm, otherUser Perm) error {
	return f.ChmodByCodes(file, owner.ToCodes(), group.ToCodes(), otherUser.ToCodes())
}

func (f *fileManager) ChmodByCodes(file string, owner [3]uint32, group [3]uint32, otherUser [3]uint32) (err error) {
	ownerP := owner[0] + owner[1] + owner[2]
	if !ValidatePermNumber(ownerP) || !ValidatePermCodes(owner) {
		err = fmt.Errorf(`invalid owner permission number: %+v`, owner)
		return
	}
	groupP := group[0] + group[1] + group[2]
	if !ValidatePermNumber(groupP) || !ValidatePermCodes(group) {
		err = fmt.Errorf(`invalid group permission number: %+v`, group)
		return
	}
	otherUserP := otherUser[0] + otherUser[1] + otherUser[2]
	if !ValidatePermNumber(otherUserP) || !ValidatePermCodes(otherUser) {
		err = fmt.Errorf(`invalid otherUser permission number: %+v`, otherUser)
		return
	}
	v := fmt.Sprintf(`%d%d%d`, ownerP, groupP, otherUserP)
	var n uint64
	n, err = strconv.ParseUint(v, 8, 32)
	if err != nil {
		return
	}
	err = f.Root.Chmod(file, os.FileMode(uint32(n)))
	return
}

func (f *fileManager) Chown(file string, username string, groupName ...string) (err error) {
	var u *user.User
	u, err = user.Lookup(username)
	if err != nil {
		return
	}
	uid := param.AsInt(u.Uid)
	if uid <= 0 {
		err = fmt.Errorf(`failed to parse uid of %s: %s`, username, u.Uid)
		return
	}
	var gid int
	if len(groupName) == 0 || len(groupName[0]) == 0 {
		gid = param.AsInt(u.Gid)
		if gid <= 0 {
			err = fmt.Errorf(`failed to parse gid of %s: %s`, username, u.Gid)
			return
		}
	} else {
		var g *user.Group
		g, err = user.LookupGroup(groupName[0])
		if err != nil {
			return
		}
		gid = param.AsInt(g.Gid)
		if gid <= 0 {
			err = fmt.Errorf(`failed to parse gid of %s: %s`, groupName, u.Gid)
			return
		}
	}
	err = f.Root.Chown(file, uid, gid)
	return
}

func (f *fileManager) ChownByID(file string, uid int, gid int) (err error) {
	err = f.Root.Chown(file, uid, gid)
	return
}

func (f *fileManager) enterPath(file string) (d http.File, fi os.FileInfo, err error) {
	d, err = f.Root.Open(file)
	if err != nil {
		return
	}
	//defer d.Close()
	fi, err = d.Stat()
	return
}

func (f *fileManager) Upload(fpath string,
	chunkUpload *uploadClient.ChunkUpload,
	chunkOpts ...uploadClient.ChunkInfoOpter) (err error) {
	var (
		d  http.File
		fi os.FileInfo
	)
	d, fi, err = f.enterPath(fpath)
	if d != nil {
		defer d.Close()
	}
	if err != nil {
		return
	}
	if !fi.IsDir() {
		return errors.New(f.T(`路径不正确: %s`, fpath))
	}
	var filePath string
	var chunked bool // 是否支持分片
	if chunkUpload != nil {
		_, err := chunkUpload.Upload(f.Request().StdRequest(), chunkOpts...)
		if err != nil {
			if !errors.Is(err, uploadClient.ErrChunkUnsupported) {
				if errors.Is(err, uploadClient.ErrChunkUploadCompleted) ||
					errors.Is(err, uploadClient.ErrFileUploadCompleted) {
					return nil
				}
				return err
			}
		} else {
			if !chunkUpload.Merged() {
				return nil
			}
			chunked = true
			filePath = chunkUpload.GetSavePath()
		}
	}
	absPath := filepath.Join(f.Root.Name(), filepath.Clean(`/`+fpath))
	if !chunked {
		fileHdr, err := f.SaveUploadedFile(`file`, absPath)
		if err != nil {
			return err
		}
		filePath = filepath.Join(absPath, filepath.Base(fileHdr.Filename))
	}
	pipe := f.Form(`pipe`)
	switch pipe {
	case `unzip`:
		err = com.Unzip(filePath, absPath)
		if err == nil {
			err = os.Remove(filePath)
			if err != nil {
				if !os.IsNotExist(err) {
					err = errors.New(f.T(`压缩包已经成功解压，但是删除压缩包失败：`) + err.Error())
				} else {
					err = nil
				}
			}
		}
		return err
	default:
		if chunked {
			newfile := filepath.Join(absPath, filepath.Base(filePath))
			err = com.Rename(filePath, newfile)
			if err != nil {
				if !os.IsNotExist(err) {
					err = fmt.Errorf(`move %s to %s: %w`, filePath, newfile, err)
				} else {
					err = nil
				}
			}
		}
		return err
	}
}

func (f *fileManager) List(absPath string, sortBy ...string) (err error, exit bool, dirs []os.FileInfo) {
	var (
		d  http.File
		fi os.FileInfo
	)
	d, fi, err = f.enterPath(absPath)
	if d != nil {
		defer d.Close()
	}
	if err != nil {
		return
	}
	if !fi.IsDir() {
		fileName := filepath.Base(absPath)
		inline := f.Formx(`inline`).Bool()
		return f.Attachment(d, fileName, fi.ModTime(), inline), true, nil
	}

	dirs, err = d.Readdir(-1)
	if len(sortBy) > 0 {
		switch sortBy[0] {
		case `time`:
			sort.Sort(SortByModTime(dirs))
		case `-time`:
			sort.Sort(SortByModTimeDesc(dirs))
		case `name`:
		case `-name`:
			sort.Sort(SortByNameDesc(dirs))
		case `type`:
			fallthrough
		default:
			sort.Sort(SortByFileType(dirs))
		}
	} else {
		sort.Sort(SortByFileType(dirs))
	}
	if f.Format() == "json" {
		dirList, fileList := f.ListTransfer(dirs)
		data := f.Data()
		data.SetData(echo.H{
			`dirList`:  dirList,
			`fileList`: fileList,
		})
		return f.JSON(data), true, nil
	}
	return
}

func (f *fileManager) ListTransfer(dirs []os.FileInfo) (dirList []echo.H, fileList []echo.H) {
	dirList = []echo.H{}
	fileList = []echo.H{}
	for _, d := range dirs {
		item := echo.H{
			`name`:  d.Name(),
			`size`:  d.Size(),
			`mode`:  d.Mode().String(),
			`mtime`: d.ModTime().Format(`2006-01-02 15:04:05`),
			//`sys`:   d.Sys(),
		}
		if d.IsDir() {
			dirList = append(dirList, item)
			continue
		}
		fileList = append(fileList, item)
	}
	return
}
