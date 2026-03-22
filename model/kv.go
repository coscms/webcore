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

package model

import (
	"strings"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"golang.org/x/sync/singleflight"

	"github.com/coscms/webcore/dbschema"
)

func NewKv(ctx echo.Context) *Kv {
	m := &Kv{
		NgingKv: dbschema.NewNgingKv(ctx),
	}
	return m
}

type Kv struct {
	*dbschema.NgingKv
}

func (s *Kv) check() error {
	ctx := s.Context()
	s.Key = strings.TrimSpace(s.Key)
	if len(s.Key) == 0 {
		return ctx.NewError(code.InvalidParameter, `键不能为空`).SetZone(`key`)
	}
	s.Type = strings.TrimSpace(s.Type)
	if len(s.Type) == 0 {
		return ctx.NewError(code.InvalidParameter, `类型不能为空`).SetZone(`type`)
	}
	var (
		exists bool
		err    error
	)
	if s.Id > 0 { // edit
		exists, err = s.Exists(nil, db.And(
			db.Cond{`key`: s.Key},
			db.Cond{`type`: s.Type},
			db.Cond{`id`: db.NotEq(s.Id)},
		))
	} else {
		exists, err = s.Exists(nil, db.And(
			db.Cond{`key`: s.Key},
			db.Cond{`type`: s.Type},
		))
	}
	if err != nil {
		return err
	}
	if exists {
		return ctx.NewError(code.DataAlreadyExists, `键"%v"已经存在`, s.Key).SetZone(`key`)
	}
	return nil
}

func (s *Kv) Get(mw func(db.Result) db.Result, args ...interface{}) error {
	err := s.NgingKv.Get(mw, args...)
	if err != nil {
		return err
	}
	return nil
}

var sg singleflight.Group

// AutoCreateKey 自动创建 key
// value: 0. 值; 1. 说明; 2. 帮助说明
func (s *Kv) AutoCreateKey(typ string, key string, value ...string) error {
	if len(typ) == 0 {
		typ = AutoCreatedType
	}
	_, err, _ := sg.Do(typ+`&`+key, func() (interface{}, error) {
		err := s.autoCreateKey(typ, key, value...)
		return nil, err
	})
	return err
}

func (s *Kv) autoCreateKey(typ string, key string, value ...string) error {
	typParts := strings.SplitN(typ, `|`, 2)
	typ = typParts[0]
	m := dbschema.NewNgingKv(s.Context())
	m.Key = key
	m.Type = typ
	m.ChildKeyType = KvDefaultDataType
	com.SliceExtract(value, &m.Value, &m.Description, &m.Help)
	m.Updated = uint(time.Now().Unix())
	_, err := m.Insert()
	if err != nil {
		return err
	}

	var exists bool
	exists, err = m.Exists(nil, `key`, typ)
	if err != nil || exists {
		return err
	}

	m.Reset()
	m.Key = typ
	m.Type = KvRootType
	m.ChildKeyType = KvDefaultDataType
	if len(typParts) >= 2 {
		m.Value = typParts[1]
	}
	if len(m.Value) == 0 && m.Key == AutoCreatedType {
		m.Value = `自动创建`
	}
	_, err = m.Insert()
	return err
}

// GetValue 获取 key 的值
// defaultValue: 0. 默认值; 1. 说明; 2. 帮助说明 (1 和 2 仅在自动创建时有用)
func (s *Kv) GetValue(key string, defaultValue ...string) (string, error) {
	err := s.NgingKv.Get(func(r db.Result) db.Result {
		return r.Select(`value`)
	}, db.And(
		db.Cond{`key`: key},
		db.Cond{`type`: db.NotEq(KvRootType)},
	))
	if err != nil {
		if err == db.ErrNoMoreRows {
			if err = s.AutoCreateKey(AutoCreatedType, key, defaultValue...); err != nil {
				s.Context().Logger().Error(err)
			}
		}
		if len(defaultValue) > 0 {
			return defaultValue[0], err
		}
		return ``, err
	}
	if len(defaultValue) > 0 && len(s.Value) == 0 {
		return defaultValue[0], err
	}
	return s.Value, err
}

// GetTypeValues
// typ: "type|typeName"
// defaultValue: {"key":"Value|Description|Help"}
func (s *Kv) GetTypeValues(typ string, defaultValue ...map[string]string) (map[string]string, error) {
	_, err := s.NgingKv.ListByOffset(nil, func(r db.Result) db.Result {
		return r.Select(`value`)
	}, 0, -1, db.Cond{`type`: strings.SplitN(typ, `|`, 2)[0]})
	if err != nil {
		if len(defaultValue) > 0 {
			result := make(map[string]string, len(defaultValue[0]))
			for key, value := range defaultValue[0] {
				values := strings.SplitN(value, `|`, 2)
				result[key] = values[0]
			}
			return result, err
		}
		return nil, err
	}
	rows := s.Objects()
	if len(rows) == 0 {
		if len(defaultValue) > 0 {
			result := make(map[string]string, len(defaultValue[0]))
			for key, value := range defaultValue[0] {
				values := strings.SplitN(value, `|`, 3)
				if err = s.AutoCreateKey(typ, key, values...); err != nil {
					s.Context().Logger().Error(err)
				}
				result[key] = values[0]
			}
			return result, err
		}
		return nil, err
	}
	values := make(map[string]string, len(rows))
	for _, row := range rows {
		values[row.Key] = row.Value
	}
	return values, err
}

func (s *Kv) Add() (pk interface{}, err error) {
	if err = s.check(); err != nil {
		return nil, err
	}
	s.NgingKv.Updated = uint(time.Now().Unix())
	return s.NgingKv.Insert()
}

func (s *Kv) Edit(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	if err = s.check(); err != nil {
		return err
	}
	return s.NgingKv.Update(mw, args...)
}

func (s *Kv) Delete(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	var rows []*dbschema.NgingKv
	s.NgingKv.ListByOffset(&rows, nil, 0, -1, args...)
	var types []string
	for _, row := range rows {
		if row.Type != KvRootType {
			continue
		}
		if com.InSlice(row.Key, types) {
			continue
		}
		types = append(types, row.Key)
	}
	if len(types) > 0 {
		err = s.NgingKv.Delete(nil, db.Cond{`type`: db.In(types)})
		if err != nil {
			return
		}
	}
	return s.NgingKv.Delete(mw, args...)
}

func (s *Kv) IsRootType(typ string) bool {
	return typ == KvRootType
}

func (s *Kv) SetSingleField(id int, field string, value string) error {
	set := echo.H{}
	switch field {
	case "value", "key", "sort", "child_key_type":
		set[field] = value
	default:
		return s.Context().E(`不支持修改字段: %v`, field)
	}
	return s.UpdateFields(nil, set, `id`, id)
}

func (s *Kv) KvTypeList(excludeIDs ...uint) []*dbschema.NgingKv {
	cond := db.NewCompounds()
	cond.AddKV(`type`, KvRootType)
	if len(excludeIDs) > 0 && excludeIDs[0] > 0 {
		cond.AddKV(`id`, db.NotEq(excludeIDs[0]))
	}
	_, err := s.ListByOffset(nil, func(r db.Result) db.Result {
		return r.OrderBy(`sort`)
	}, 0, -1, cond.And())
	if err == nil {
		return s.Objects()
	}
	return nil
}

func (s *Kv) ListByType(typ string, excludeIDs ...uint) []*dbschema.NgingKv {
	cond := db.NewCompounds()
	cond.AddKV(`type`, typ)
	if len(excludeIDs) > 0 && excludeIDs[0] > 0 {
		cond.AddKV(`id`, db.NotEq(excludeIDs[0]))
	}
	_, err := s.ListByOffset(nil, func(r db.Result) db.Result {
		return r.OrderBy(`sort`)
	}, 0, -1, cond.And())
	if err == nil {
		return s.Objects()
	}
	return nil
}

func (s *Kv) GetFromTypeList(typeList []*dbschema.NgingKv, key string) string {
	if key == KvRootType {
		return KvRootType
	}
	for _, row := range typeList {
		if row.Key == key {
			return row.Value
		}
	}
	return key
}

func (s *Kv) ListToMap(typeList []*dbschema.NgingKv) map[string]*dbschema.NgingKv {
	r := map[string]*dbschema.NgingKv{}
	for _, row := range typeList {
		r[row.Key] = row
	}
	return r
}
