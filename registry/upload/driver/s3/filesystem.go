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

package s3

import (
	"context"
	"io"
	"os"
	"path"

	"github.com/webx-top/echo"

	"github.com/admpub/errors"
	"github.com/coscms/webcore/library/s3manager"
	"github.com/coscms/webcore/library/s3manager/fileinfo"
	"github.com/coscms/webcore/library/s3manager/s3client"
	uploadLibrary "github.com/coscms/webcore/library/upload"
	"github.com/coscms/webcore/model"
	"github.com/coscms/webcore/registry/upload"
	"github.com/coscms/webcore/registry/upload/driver"
	"github.com/coscms/webcore/registry/upload/driver/local"
)

const Name = `s3`

var _ upload.Storer = &Filesystem{}

func init() {
	upload.StorerRegister(Name, func(ctx context.Context, subdir string, options ...driver.Option) (upload.Storer, error) {
		return NewFilesystem(ctx, subdir, options...)
	})
}

func NewFilesystem(ctx context.Context, subdir string, options ...driver.Option) (*Filesystem, error) {
	base := local.NewFilesystem(ctx, subdir, options...)
	m, err := model.GetCloudStorage(ctx, base.Config.StorerID)
	if err != nil {
		return nil, errors.WithMessage(err, Name)
	}
	mgr := s3client.New(m.NgingCloudStorage, 0)
	if _, err := mgr.Connect(); err != nil {
		return nil, errors.WithMessage(err, Name)
	}
	base.Config.BaseURL = m.Baseurl
	return &Filesystem{
		Filesystem: base,
		model:      m,
		mgr:        mgr,
	}, nil
}

// Filesystem 文件系统存储引擎
type Filesystem struct {
	*local.Filesystem
	model *model.CloudStorage
	mgr   *s3manager.S3Manager
}

// Name 引擎名
func (f *Filesystem) Name() string {
	return Name
}

func (f *Filesystem) ErrIsNotExist(err error) bool {
	return f.mgr.ErrIsNotExist(err)
}

// Exists 判断文件是否存在
func (f *Filesystem) Exists(ctx context.Context, file string) (bool, error) {
	return f.mgr.Exists(ctx, file)
}

// FileInfo 获取文件信息
func (f *Filesystem) FileInfo(ctx context.Context, file string) (os.FileInfo, error) {
	objectInfo, err := f.mgr.Stat(ctx, file)
	if err != nil {
		return nil, errors.WithMessage(err, Name)
	}
	return fileinfo.New(objectInfo), nil
}

// SendFile 下载文件
func (f *Filesystem) SendFile(ctx echo.Context, file string) error {
	ctx.Request().Form().Set(`inline`, `1`)
	err := f.mgr.Download(ctx, file)
	if err != nil {
		err = errors.WithMessage(err, Name)
	}
	return err
}

// FileDir 物理路径文件夹
func (f *Filesystem) FileDir(subpath string) string {
	return path.Join(uploadLibrary.UploadURLPath, f.Subdir, subpath)
}

// Put 上传文件
func (f *Filesystem) Put(ctx context.Context, dstFile string, src io.Reader, size int64) (savePath string, viewURL string, err error) {
	savePath = f.FileDir(dstFile)
	//viewURL = `[storage:`+param.AsString(f.model.Id)+`]`+f.URLDir(dstFile)
	viewURL = f.PublicURL(dstFile)
	err = f.mgr.Put(ctx, src, savePath, size)
	if err != nil {
		err = errors.WithMessage(err, Name)
	}
	return
}

// Get 获取文件读取接口
func (f *Filesystem) Get(ctx context.Context, dstFile string) (io.ReadCloser, error) {
	object, err := f.mgr.Get(ctx, dstFile)
	if err != nil {
		return nil, errors.WithMessage(err, Name)
	}
	exists, err := f.mgr.StatIsExists(object.Stat())
	if !exists {
		err = os.ErrNotExist
	}
	return object, err
}

// Delete 删除文件
func (f *Filesystem) Delete(ctx context.Context, dstFile string) error {
	err := f.mgr.Remove(ctx, dstFile)
	if err != nil {
		err = errors.WithMessage(err, Name)
	}
	return err
}

// DeleteDir 删除文件夹及其内部文件
func (f *Filesystem) DeleteDir(ctx context.Context, dstDir string) error {
	err := f.mgr.RemoveDir(ctx, dstDir)
	if err != nil {
		err = errors.WithMessage(err, Name)
	}
	return err
}

// Move 移动文件
func (f *Filesystem) Move(ctx context.Context, src, dst string) error {
	err := f.mgr.Rename(ctx, src, dst)
	if err != nil {
		err = errors.WithMessage(err, Name)
	}
	return err
}

// Close 关闭连接
func (f *Filesystem) Close() error {
	return nil
}

// FixURL 改写文件网址
func (f *Filesystem) FixURL(content string, embedded ...bool) string {
	rowsByID := f.model.CachedList()
	return uploadLibrary.ReplacePlaceholder(content, func(id string) string {
		r, y := rowsByID[id]
		if !y {
			return ``
		}
		return r.Baseurl
	})
}
