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
	"io"
	"log"
	"time"

	"github.com/coscms/webcore/dbschema"
	"github.com/webx-top/echo"
)

func newStorageMock(ctx echo.Context, cfg dbschema.NgingCloudBackup) (Storager, error) {
	return NewStorageMock(), nil
}

func NewStorageMock() Storager {
	return &StorageMock{}
}

type StorageMock struct {
}

func (s *StorageMock) Connect() (err error) {
	log.Println(`StorageMock: Connect`)
	return
}

func (s *StorageMock) Put(ctx context.Context, reader io.Reader, ppath string, size int64) (err error) {
	log.Println(`StorageMock: Put --->`, ppath)
	time.Sleep(time.Second * 2)
	return
}

func (s *StorageMock) Download(ctx context.Context, ppath string, w io.Writer) error {
	log.Println(`StorageMock: Download --->`, ppath)
	return nil
}

func (s *StorageMock) Restore(ctx context.Context, ppath string, destpath string, callback func(from, to string)) error {
	log.Println(`StorageMock: Restore --->`, ppath)
	return nil
}

func (s *StorageMock) RemoveDir(ctx context.Context, ppath string) error {
	log.Println(`StorageMock: RemoveDir --->`, ppath)
	return nil
}

func (s *StorageMock) Remove(ctx context.Context, ppath string) error {
	log.Println(`StorageMock: Remove --->`, ppath)
	return nil
}

func (s *StorageMock) Close() (err error) {
	log.Println(`StorageMock: Close`)
	return
}
