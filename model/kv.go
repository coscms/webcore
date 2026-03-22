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

const (
	HKeyDescription = `description`
	HKeyHelp        = `help`
)

// GetTypeValues
// typ: "type|typeName"
// defaultValue: echo.NewKVData().Add(`key`, `value`, echo.KVOptHKV(`description`, `说明`), echo.KVOptHKV(`help`, `帮助`))
func (s *Kv) GetTypeValues(typ string, defaultValue ...*echo.KVData) (echo.KVList, error) {
	_, err := s.NgingKv.ListByOffset(nil, func(r db.Result) db.Result {
		return r.Select(`key`, `value`, `description`, `help`).OrderBy(`sort`, `id`)
	}, 0, -1, db.Cond{`type`: strings.SplitN(typ, `|`, 2)[0]})
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0].Clone().Slice(), err
		}
		return nil, err
	}
	rows := s.Objects()
	if len(rows) == 0 {
		if len(defaultValue) > 0 {
			for _, item := range defaultValue[0].Slice() {
				var description string
				var help string
				if item.H != nil {
					description = item.H.String(HKeyDescription)
					help = item.H.String(HKeyHelp)
				}
				if err = s.AutoCreateKey(typ, item.K, item.V, description, help); err != nil {
					s.Context().Logger().Error(err)
				}
			}
			return defaultValue[0].Clone().Slice(), err
		}
		return nil, err
	}
	values := make(echo.KVList, 0, len(rows))
	for _, row := range rows {
		var h echo.H
		if len(defaultValue) > 0 {
			item := defaultValue[0].GetItem(row.Key)
			if item != nil && item.H != nil {
				h = item.H.Clone()
			}
		}
		if h == nil {
			h = echo.H{}
		}
		h.Set(HKeyDescription, row.Description)
		h.Set(HKeyHelp, row.Help)
		values.Add(row.Key, row.Value, echo.KVOptH(h))
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
