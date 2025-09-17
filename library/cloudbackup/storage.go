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

package cloudbackup

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/coscms/webcore/dbschema"
	"github.com/coscms/webcore/library/notice"
	"github.com/coscms/webcore/model"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

type Storager interface {
	Connect() (err error)
	Put(ctx context.Context, reader io.Reader, ppath string, size int64) (err error)
	Download(ctx context.Context, ppath string, w io.Writer) error
	RemoveDir(ctx context.Context, ppath string) error
	Remove(ctx context.Context, ppath string) error
	Restore(ctx context.Context, ppath string, destpath string, callback func(from, to string)) error
	Close() (err error)
}

type ProgressorSetter interface {
	SetProgressor(notice.Progressor)
}

type Form struct {
	Type        string
	Label       string
	Name        string
	Required    bool
	Pattern     string
	Placeholder string
}

var Forms = map[string][]Form{}

func HasForm(engineName string, formName string) bool {
	for _, f := range Forms[engineName] {
		if f.Name == formName {
			return true
		}
	}
	return false
}

type Constructor func(ctx echo.Context, cfg dbschema.NgingCloudBackup) (Storager, error)

var storages = map[string]Constructor{
	`mock`: newStorageMock,
}

func Register(name string, constructor Constructor, forms []Form, label string) {
	storages[name] = constructor
	Forms[name] = forms
	model.CloudBackupStorageEngines.Add(name, label)
}

var ErrUnsupported = errors.New(`unsupported storage engine`)
var ErrEmptyConfig = errors.New(`empty config`)

func NewStorage(ctx echo.Context, cfg dbschema.NgingCloudBackup) (Storager, error) {
	cr, ok := storages[cfg.StorageEngine]
	if !ok {
		return nil, fmt.Errorf(`%w: %s`, ErrUnsupported, cfg.StorageEngine)
	}
	return cr(ctx, cfg)
}

func DownloadFile(s Storager, ctx context.Context, ppath string, dest string) error {
	dir := filepath.Dir(dest)
	com.MkdirAll(dir, os.ModePerm)
	fi, err := os.Create(dest)
	if err != nil {
		return err
	}
	err = s.Download(ctx, ppath, fi)
	fi.Close()
	return err
}
