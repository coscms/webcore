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

package file

import (
	"fmt"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"

	"github.com/coscms/webcore/dbschema"
	"github.com/coscms/webcore/library/fileupdater"
	uploadLibrary "github.com/coscms/webcore/library/upload"
)

func NewEmbedded(ctx echo.Context, fileMdls ...*File) *Embedded {
	if ctx == nil {
		panic(`ctx is nil`)
	}
	var fileM *File
	if len(fileMdls) > 0 {
		fileM = fileMdls[0]
	} else {
		fileM = NewFile(ctx)
	}
	m := &Embedded{
		NgingFileEmbedded: dbschema.NewNgingFileEmbedded(ctx),
		File:              fileM,
		Moved:             NewMoved(ctx),
	}
	return m
}

type Embedded struct {
	*dbschema.NgingFileEmbedded
	File             *File
	Moved            *Moved
	replacedViewURLs map[string]string // viewURL => newViewURL
	updater          *fileupdater.FileUpdater
}

func (f *Embedded) Updater(table string, field string, tableID string) *fileupdater.FileUpdater {
	if f.updater == nil {
		f.updater = fileupdater.New(f)
	}
	f.updater.Set(table, field, tableID)
	return f.updater
}

func (f *Embedded) FileIDs() []uint64 {
	var fileIDs []uint64
	if len(f.FileIds) == 0 {
		return fileIDs
	}
	fileStrIDs := strings.Split(f.FileIds, `,`)
	fileIDs = make([]uint64, 0, len(fileStrIDs))
	for _, fileID := range fileStrIDs {
		fileIDs = append(fileIDs, param.AsUint64(fileID))
	}
	return fileIDs
}

func (f *Embedded) ReplacedViewURLs() map[string]string {
	return f.replacedViewURLs
}

// DeleteByTableID 删除嵌入文件
func (f *Embedded) DeleteByTableID(project string, table string, tableID string) error {
	_, err := f.ListByOffset(nil, nil, 0, -1, db.And(
		db.Cond{`table_id`: tableID},
		db.Cond{`table_name`: table},
		db.Cond{`project`: project},
	))
	if err != nil {
		return err
	}
	var ids []uint64
	for _, row := range f.Objects() {
		err = f.File.UpdateField(nil, `used_times`, db.Raw(`used_times-1`), db.And(
			db.Cond{`used_times`: db.Gt(0)},
			db.Cond{`id`: db.In(strings.Split(row.FileIds, `,`))},
		))
		if err != nil {
			return err
		}
		ids = append(ids, row.Id)
	}
	if len(ids) > 0 {
		err = f.Delete(nil, db.Cond{`id`: db.In(ids)})
	}
	return err
}

func (f *Embedded) UpdateByFileID(project string, table string, field string, tableID string, fileID uint64) (uint64, error) {
	err := f.File.Incr(fileID)
	if err != nil {
		return 0, err
	}
	m := dbschema.NewNgingFileEmbedded(f.Context())
	err = m.Get(nil, db.And(
		db.Cond{`table_id`: tableID},
		db.Cond{`field_name`: field},
		db.Cond{`table_name`: table},
	))
	var newID uint64
	if err != nil {
		if err != db.ErrNoMoreRows {
			return newID, err
		}
		m.Reset()
		m.FieldName = field
		m.TableName = table
		m.Project = project
		m.TableId = tableID
		m.FileIds = fmt.Sprint(fileID)
		_, err = m.Insert()
		newID = m.Id
	}
	return newID, err
}

func (f *Embedded) UpdateEmbedded(embedded bool, project string, table string, field string, tableID string, fileIds ...uint64) (err error) {
	m := dbschema.NewNgingFileEmbedded(f.Context())
	err = m.Get(nil, db.And(
		db.Cond{`table_id`: tableID},
		db.Cond{`field_name`: field},
		db.Cond{`table_name`: table},
	))
	if err != nil {
		if err != db.ErrNoMoreRows {
			return
		}
		if len(fileIds) < 1 {
			err = nil
			return
		}
		// 不存在时，添加
		m.Reset()
		m.FieldName = field
		m.TableName = table
		m.Project = project
		m.TableId = tableID
		if embedded {
			m.Embedded = `Y`
		} else {
			m.Embedded = `N`
		}
		m.FileIds = ""
		err = f.File.Incr(fileIds...)
		if err != nil {
			return
		}
		for _, v := range fileIds {
			m.FileIds += com.Str(v) + `,`
		}
		m.FileIds = strings.TrimSuffix(m.FileIds, ",")
		f.FileIds = m.FileIds // 供FileIDs()使用
		_, err = m.Insert()
		return
	}
	isEmpty := len(fileIds) < 1
	if isEmpty { // 删除关联记录
		err = f.DeleteByInstance(m)
		return
	}
	postFidsString := com.JoinNumbers(fileIds, `,`)
	if m.FileIds == postFidsString {
		return
	}
	oldFids := param.Split(m.FileIds, ",").Uint64()
	//新增引用
	err = f.AddFileByIds(fileIds, oldFids...)
	if err != nil {
		return
	}
	//已删除引用
	err = f.DeleteFileByIds(oldFids, fileIds...)
	if err != nil {
		return
	}
	m.FileIds = postFidsString
	f.FileIds = m.FileIds // 供FileIDs()使用
	err = f.UpdateField(nil, `file_ids`, m.FileIds, db.Cond{`id`: m.Id})
	return
}

// FilterNotExistsFileIDs 仅仅返回数据库中有记录的文件ID
func (f *Embedded) FilterNotExistsFileIDs(fids []uint64, files []string) []uint64 {
	if len(fids) > 0 {
		fids = f.File.GetIDByIDs(fids)
	}
	if len(files) == 0 {
		return fids
	}
	ids := f.File.GetIDByViewURLs(files)
	if len(fids) == 0 {
		return ids
	}
	for _, id := range ids {
		if !com.InUint64Slice(id, fids) {
			fids = append(fids, id)
		}
	}
	return fids
}

// RelationEmbeddedFiles 关联嵌入的文件
// @param project 项目名称
// @param table 表名称
// @param field 被嵌入的字段名
// @param tableID 表中行主键ID
// @param v 内容
// @return
// @author AdamShen <swh@admpub.com>
func (f *Embedded) RelationEmbeddedFiles(project string, table string, field string, tableID string, v string) error {
	var (
		files      []string
		oldFileIDs []uint64 //旧文件ID
	)
	uploadLibrary.ParseEmbeddedRes(v, func(file string, fileID uint64) {
		var exists bool
		if fileID > 0 {
			exists = com.InUint64Slice(fileID, oldFileIDs)
		} else {
			exists = com.InSlice(file, files)
		}
		if exists {
			return
		}
		if fileID > 0 {
			oldFileIDs = append(oldFileIDs, fileID)
		} else {
			files = append(files, file)
		}
	})

	// 仅仅提取数据库中有记录的数据
	oldFileIDs = f.FilterNotExistsFileIDs(oldFileIDs, files)

	err := f.UpdateEmbedded(true, project, table, field, tableID, oldFileIDs...)
	return err
}

func (f *Embedded) RelationFiles(project string, table string, field string, tableID string, v string, seperator ...string) error {
	var (
		files      []string
		oldFileIDs []uint64 //旧文件ID
	)
	uploadLibrary.ParseRelatedRes(v, func(file string, fileID uint64) {
		var exists bool
		if fileID > 0 {
			exists = com.InUint64Slice(fileID, oldFileIDs)
		} else {
			if len(file) == 0 {
				return
			}
			exists = com.InSlice(file, files)
		}
		if exists {
			return
		}
		if fileID > 0 {
			oldFileIDs = append(oldFileIDs, fileID)
		} else {
			files = append(files, file)
		}
	}, seperator...)

	// 仅仅提取数据库中有记录的数据
	oldFileIDs = f.FilterNotExistsFileIDs(oldFileIDs, files)

	err := f.UpdateEmbedded(false, project, table, field, tableID, oldFileIDs...)
	return err
}
