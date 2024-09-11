// @generated Do not edit this file, which is automatically generated by the generator.

package dbschema

import (
	"fmt"

	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/lib/factory/pagination"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

type Slice_NgingOauthAgree []*NgingOauthAgree

func (s Slice_NgingOauthAgree) Range(fn func(m factory.Model) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_NgingOauthAgree) RangeRaw(fn func(m *NgingOauthAgree) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_NgingOauthAgree) GroupBy(keyField string) map[string][]*NgingOauthAgree {
	r := map[string][]*NgingOauthAgree{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		if _, y := r[vkey]; !y {
			r[vkey] = []*NgingOauthAgree{}
		}
		r[vkey] = append(r[vkey], row)
	}
	return r
}

func (s Slice_NgingOauthAgree) KeyBy(keyField string) map[string]*NgingOauthAgree {
	r := map[string]*NgingOauthAgree{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = row
	}
	return r
}

func (s Slice_NgingOauthAgree) AsKV(keyField string, valueField string) param.Store {
	r := param.Store{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = dmap[valueField]
	}
	return r
}

func (s Slice_NgingOauthAgree) Transform(transfers map[string]param.Transfer) []param.Store {
	r := make([]param.Store, len(s))
	for idx, row := range s {
		r[idx] = row.AsMap().Transform(transfers)
	}
	return r
}

func (s Slice_NgingOauthAgree) FromList(data interface{}) Slice_NgingOauthAgree {
	values, ok := data.([]*NgingOauthAgree)
	if !ok {
		for _, value := range data.([]interface{}) {
			row := &NgingOauthAgree{}
			row.FromRow(value.(map[string]interface{}))
			s = append(s, row)
		}
		return s
	}
	s = append(s, values...)

	return s
}

func NewNgingOauthAgree(ctx echo.Context) *NgingOauthAgree {
	m := &NgingOauthAgree{}
	m.SetContext(ctx)
	return m
}

// NgingOauthAgree oauth2服务端用户授权表
type NgingOauthAgree struct {
	base    factory.Base
	objects []*NgingOauthAgree

	Uid     uint   `db:"uid,pk" bson:"uid" comment:"用户ID" json:"uid" xml:"uid"`
	AppId   string `db:"app_id,pk" bson:"app_id" comment:"AppID" json:"app_id" xml:"app_id"`
	Scopes  string `db:"scopes" bson:"scopes" comment:"授权信息" json:"scopes" xml:"scopes"`
	Created uint   `db:"created" bson:"created" comment:"创建时间" json:"created" xml:"created"`
	Updated uint   `db:"updated" bson:"updated" comment:"更新时间" json:"updated" xml:"updated"`
}

// - base function

func (a *NgingOauthAgree) Trans() factory.Transactioner {
	return a.base.Trans()
}

func (a *NgingOauthAgree) Use(trans factory.Transactioner) factory.Model {
	a.base.Use(trans)
	return a
}

func (a *NgingOauthAgree) SetContext(ctx echo.Context) factory.Model {
	a.base.SetContext(ctx)
	return a
}

func (a *NgingOauthAgree) EventON(on ...bool) factory.Model {
	a.base.EventON(on...)
	return a
}

func (a *NgingOauthAgree) EventOFF(off ...bool) factory.Model {
	a.base.EventOFF(off...)
	return a
}

func (a *NgingOauthAgree) Context() echo.Context {
	return a.base.Context()
}

func (a *NgingOauthAgree) SetConnID(connID int) factory.Model {
	a.base.SetConnID(connID)
	return a
}

func (a *NgingOauthAgree) ConnID() int {
	return a.base.ConnID()
}

func (a *NgingOauthAgree) SetNamer(namer func(factory.Model) string) factory.Model {
	a.base.SetNamer(namer)
	return a
}

func (a *NgingOauthAgree) Namer() func(factory.Model) string {
	return a.base.Namer()
}

func (a *NgingOauthAgree) SetParam(param *factory.Param) factory.Model {
	a.base.SetParam(param)
	return a
}

func (a *NgingOauthAgree) Param(mw func(db.Result) db.Result, args ...interface{}) *factory.Param {
	if a.base.Param() == nil {
		return a.NewParam().SetMiddleware(mw).SetArgs(args...)
	}
	return a.base.Param().SetMiddleware(mw).SetArgs(args...)
}

func (a *NgingOauthAgree) New(structName string, connID ...int) factory.Model {
	return a.base.New(structName, connID...)
}

func (a *NgingOauthAgree) Base_() factory.Baser {
	return &a.base
}

// - current function

func (a *NgingOauthAgree) Objects() []*NgingOauthAgree {
	if a.objects == nil {
		return nil
	}
	return a.objects[:]
}

func (a *NgingOauthAgree) XObjects() Slice_NgingOauthAgree {
	return Slice_NgingOauthAgree(a.Objects())
}

func (a *NgingOauthAgree) NewObjects() factory.Ranger {
	return &Slice_NgingOauthAgree{}
}

func (a *NgingOauthAgree) InitObjects() *[]*NgingOauthAgree {
	a.objects = []*NgingOauthAgree{}
	return &a.objects
}

func (a *NgingOauthAgree) NewParam() *factory.Param {
	return factory.NewParam(factory.DefaultFactory).SetIndex(a.base.ConnID()).SetTrans(a.base.Trans()).SetCollection(a.Name_()).SetModel(a)
}

func (a *NgingOauthAgree) Short_() string {
	return "nging_oauth_agree"
}

func (a *NgingOauthAgree) Struct_() string {
	return "NgingOauthAgree"
}

func (a *NgingOauthAgree) Name_() string {
	b := a
	if b == nil {
		b = &NgingOauthAgree{}
	}
	if b.base.Namer() != nil {
		return WithPrefix(b.base.Namer()(b))
	}
	return WithPrefix(factory.TableNamerGet(b.Short_())(b))
}

func (a *NgingOauthAgree) CPAFrom(source factory.Model) factory.Model {
	a.SetContext(source.Context())
	a.SetConnID(source.ConnID())
	a.SetNamer(source.Namer())
	return a
}

func (a *NgingOauthAgree) Get(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	base := a.base
	if !a.base.Eventable() {
		err = a.Param(mw, args...).SetRecv(a).One()
		a.base = base
		return
	}
	queryParam := a.Param(mw, args...).SetRecv(a)
	if err = DBI.FireReading(a, queryParam); err != nil {
		return
	}
	err = queryParam.One()
	a.base = base
	if err == nil {
		err = DBI.FireReaded(a, queryParam)
	}
	return
}

func (a *NgingOauthAgree) List(recv interface{}, mw func(db.Result) db.Result, page, size int, args ...interface{}) (func() int64, error) {
	if recv == nil {
		recv = a.InitObjects()
	}
	if !a.base.Eventable() {
		return a.Param(mw, args...).SetPage(page).SetSize(size).SetRecv(recv).List()
	}
	queryParam := a.Param(mw, args...).SetPage(page).SetSize(size).SetRecv(recv)
	if err := DBI.FireReading(a, queryParam); err != nil {
		return nil, err
	}
	cnt, err := queryParam.List()
	if err == nil {
		switch v := recv.(type) {
		case *[]*NgingOauthAgree:
			err = DBI.FireReaded(a, queryParam, Slice_NgingOauthAgree(*v))
		case []*NgingOauthAgree:
			err = DBI.FireReaded(a, queryParam, Slice_NgingOauthAgree(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *NgingOauthAgree) GroupBy(keyField string, inputRows ...[]*NgingOauthAgree) map[string][]*NgingOauthAgree {
	var rows Slice_NgingOauthAgree
	if len(inputRows) > 0 {
		rows = Slice_NgingOauthAgree(inputRows[0])
	} else {
		rows = Slice_NgingOauthAgree(a.Objects())
	}
	return rows.GroupBy(keyField)
}

func (a *NgingOauthAgree) KeyBy(keyField string, inputRows ...[]*NgingOauthAgree) map[string]*NgingOauthAgree {
	var rows Slice_NgingOauthAgree
	if len(inputRows) > 0 {
		rows = Slice_NgingOauthAgree(inputRows[0])
	} else {
		rows = Slice_NgingOauthAgree(a.Objects())
	}
	return rows.KeyBy(keyField)
}

func (a *NgingOauthAgree) AsKV(keyField string, valueField string, inputRows ...[]*NgingOauthAgree) param.Store {
	var rows Slice_NgingOauthAgree
	if len(inputRows) > 0 {
		rows = Slice_NgingOauthAgree(inputRows[0])
	} else {
		rows = Slice_NgingOauthAgree(a.Objects())
	}
	return rows.AsKV(keyField, valueField)
}

func (a *NgingOauthAgree) ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error) {
	if recv == nil {
		recv = a.InitObjects()
	}
	if !a.base.Eventable() {
		return a.Param(mw, args...).SetOffset(offset).SetSize(size).SetRecv(recv).List()
	}
	queryParam := a.Param(mw, args...).SetOffset(offset).SetSize(size).SetRecv(recv)
	if err := DBI.FireReading(a, queryParam); err != nil {
		return nil, err
	}
	cnt, err := queryParam.List()
	if err == nil {
		switch v := recv.(type) {
		case *[]*NgingOauthAgree:
			err = DBI.FireReaded(a, queryParam, Slice_NgingOauthAgree(*v))
		case []*NgingOauthAgree:
			err = DBI.FireReaded(a, queryParam, Slice_NgingOauthAgree(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *NgingOauthAgree) Insert() (pk interface{}, err error) {
	a.Created = uint(time.Now().Unix())
	if a.base.Eventable() {
		err = DBI.Fire("creating", a, nil)
		if err != nil {
			return
		}
	}
	pk, err = a.Param(nil).SetSend(a).Insert()

	if err == nil && a.base.Eventable() {
		err = DBI.Fire("created", a, nil)
	}
	return
}

func (a *NgingOauthAgree) Update(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	a.Updated = uint(time.Now().Unix())
	if !a.base.Eventable() {
		return a.Param(mw, args...).SetSend(a).Update()
	}
	if err = DBI.Fire("updating", a, mw, args...); err != nil {
		return
	}
	if err = a.Param(mw, args...).SetSend(a).Update(); err != nil {
		return
	}
	return DBI.Fire("updated", a, mw, args...)
}

func (a *NgingOauthAgree) Updatex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {
	a.Updated = uint(time.Now().Unix())
	if !a.base.Eventable() {
		return a.Param(mw, args...).SetSend(a).Updatex()
	}
	if err = DBI.Fire("updating", a, mw, args...); err != nil {
		return
	}
	if affected, err = a.Param(mw, args...).SetSend(a).Updatex(); err != nil {
		return
	}
	err = DBI.Fire("updated", a, mw, args...)
	return
}

func (a *NgingOauthAgree) UpdateByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (err error) {
	a.Updated = uint(time.Now().Unix())
	if !a.base.Eventable() {
		return a.Param(mw, args...).UpdateByStruct(a, fields...)
	}
	editColumns := make([]string, len(fields))
	for index, field := range fields {
		editColumns[index] = com.SnakeCase(field)
	}
	if err = DBI.FireUpdate("updating", a, editColumns, mw, args...); err != nil {
		return
	}
	if err = a.Param(mw, args...).UpdateByStruct(a, fields...); err != nil {
		return
	}
	err = DBI.FireUpdate("updated", a, editColumns, mw, args...)
	return
}

func (a *NgingOauthAgree) UpdatexByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (affected int64, err error) {
	a.Updated = uint(time.Now().Unix())
	if !a.base.Eventable() {
		return a.Param(mw, args...).UpdatexByStruct(a, fields...)
	}
	editColumns := make([]string, len(fields))
	for index, field := range fields {
		editColumns[index] = com.SnakeCase(field)
	}
	if err = DBI.FireUpdate("updating", a, editColumns, mw, args...); err != nil {
		return
	}
	if affected, err = a.Param(mw, args...).UpdatexByStruct(a, fields...); err != nil {
		return
	}
	err = DBI.FireUpdate("updated", a, editColumns, mw, args...)
	return
}

func (a *NgingOauthAgree) UpdateField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (err error) {
	return a.UpdateFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *NgingOauthAgree) UpdatexField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (affected int64, err error) {
	return a.UpdatexFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *NgingOauthAgree) UpdateFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (err error) {
	kvset["updated"] = uint(time.Now().Unix())
	if !a.base.Eventable() {
		return a.Param(mw, args...).SetSend(kvset).Update()
	}
	m := *a
	m.FromRow(kvset)
	var editColumns []string
	for column := range kvset {
		editColumns = append(editColumns, column)
	}
	if err = DBI.FireUpdate("updating", &m, editColumns, mw, args...); err != nil {
		return
	}
	if err = a.Param(mw, args...).SetSend(kvset).Update(); err != nil {
		return
	}
	return DBI.FireUpdate("updated", &m, editColumns, mw, args...)
}

func (a *NgingOauthAgree) UpdatexFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (affected int64, err error) {
	kvset["updated"] = uint(time.Now().Unix())
	if !a.base.Eventable() {
		return a.Param(mw, args...).SetSend(kvset).Updatex()
	}
	m := *a
	m.FromRow(kvset)
	var editColumns []string
	for column := range kvset {
		editColumns = append(editColumns, column)
	}
	if err = DBI.FireUpdate("updating", &m, editColumns, mw, args...); err != nil {
		return
	}
	if affected, err = a.Param(mw, args...).SetSend(kvset).Updatex(); err != nil {
		return
	}
	err = DBI.FireUpdate("updated", &m, editColumns, mw, args...)
	return
}

func (a *NgingOauthAgree) UpdateValues(mw func(db.Result) db.Result, keysValues *db.KeysValues, args ...interface{}) (err error) {
	if !a.base.Eventable() {
		return a.Param(mw, args...).SetSend(keysValues).Update()
	}
	m := *a
	m.FromRow(keysValues.Map())
	if err = DBI.FireUpdate("updating", &m, keysValues.Keys(), mw, args...); err != nil {
		return
	}
	if err = a.Param(mw, args...).SetSend(keysValues).Update(); err != nil {
		return
	}
	return DBI.FireUpdate("updated", &m, keysValues.Keys(), mw, args...)
}

func (a *NgingOauthAgree) Upsert(mw func(db.Result) db.Result, args ...interface{}) (pk interface{}, err error) {
	pk, err = a.Param(mw, args...).SetSend(a).Upsert(func() error {
		a.Updated = uint(time.Now().Unix())
		if !a.base.Eventable() {
			return nil
		}
		return DBI.Fire("updating", a, mw, args...)
	}, func() error {
		a.Created = uint(time.Now().Unix())
		if !a.base.Eventable() {
			return nil
		}
		return DBI.Fire("creating", a, nil)
	})

	if err == nil && a.base.Eventable() {
		if pk == nil {
			err = DBI.Fire("updated", a, mw, args...)
		} else {
			err = DBI.Fire("created", a, nil)
		}
	}
	return
}

func (a *NgingOauthAgree) Delete(mw func(db.Result) db.Result, args ...interface{}) (err error) {

	if !a.base.Eventable() {
		return a.Param(mw, args...).Delete()
	}
	if err = DBI.Fire("deleting", a, mw, args...); err != nil {
		return
	}
	if err = a.Param(mw, args...).Delete(); err != nil {
		return
	}
	return DBI.Fire("deleted", a, mw, args...)
}

func (a *NgingOauthAgree) Deletex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {

	if !a.base.Eventable() {
		return a.Param(mw, args...).Deletex()
	}
	if err = DBI.Fire("deleting", a, mw, args...); err != nil {
		return
	}
	if affected, err = a.Param(mw, args...).Deletex(); err != nil {
		return
	}
	err = DBI.Fire("deleted", a, mw, args...)
	return
}

func (a *NgingOauthAgree) Count(mw func(db.Result) db.Result, args ...interface{}) (int64, error) {
	return a.Param(mw, args...).Count()
}

func (a *NgingOauthAgree) Exists(mw func(db.Result) db.Result, args ...interface{}) (bool, error) {
	return a.Param(mw, args...).Exists()
}

func (a *NgingOauthAgree) Reset() *NgingOauthAgree {
	a.Uid = 0
	a.AppId = ``
	a.Scopes = ``
	a.Created = 0
	a.Updated = 0
	return a
}

func (a *NgingOauthAgree) AsMap(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["Uid"] = a.Uid
		r["AppId"] = a.AppId
		r["Scopes"] = a.Scopes
		r["Created"] = a.Created
		r["Updated"] = a.Updated
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "Uid":
			r["Uid"] = a.Uid
		case "AppId":
			r["AppId"] = a.AppId
		case "Scopes":
			r["Scopes"] = a.Scopes
		case "Created":
			r["Created"] = a.Created
		case "Updated":
			r["Updated"] = a.Updated
		}
	}
	return r
}

func (a *NgingOauthAgree) FromRow(row map[string]interface{}) {
	for key, value := range row {
		switch key {
		case "uid":
			a.Uid = param.AsUint(value)
		case "app_id":
			a.AppId = param.AsString(value)
		case "scopes":
			a.Scopes = param.AsString(value)
		case "created":
			a.Created = param.AsUint(value)
		case "updated":
			a.Updated = param.AsUint(value)
		}
	}
}

func (a *NgingOauthAgree) Set(key interface{}, value ...interface{}) {
	switch k := key.(type) {
	case map[string]interface{}:
		for kk, vv := range k {
			a.Set(kk, vv)
		}
	default:
		var (
			kk string
			vv interface{}
		)
		if k, y := key.(string); y {
			kk = k
		} else {
			kk = fmt.Sprint(key)
		}
		if len(value) > 0 {
			vv = value[0]
		}
		switch kk {
		case "Uid":
			a.Uid = param.AsUint(vv)
		case "AppId":
			a.AppId = param.AsString(vv)
		case "Scopes":
			a.Scopes = param.AsString(vv)
		case "Created":
			a.Created = param.AsUint(vv)
		case "Updated":
			a.Updated = param.AsUint(vv)
		}
	}
}

func (a *NgingOauthAgree) AsRow(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["uid"] = a.Uid
		r["app_id"] = a.AppId
		r["scopes"] = a.Scopes
		r["created"] = a.Created
		r["updated"] = a.Updated
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "uid":
			r["uid"] = a.Uid
		case "app_id":
			r["app_id"] = a.AppId
		case "scopes":
			r["scopes"] = a.Scopes
		case "created":
			r["created"] = a.Created
		case "updated":
			r["updated"] = a.Updated
		}
	}
	return r
}

func (a *NgingOauthAgree) ListPage(cond *db.Compounds, sorts ...interface{}) error {
	_, err := pagination.NewLister(a, nil, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, cond.And()).Paging(a.Context())
	return err
}

func (a *NgingOauthAgree) ListPageAs(recv interface{}, cond *db.Compounds, sorts ...interface{}) error {
	_, err := pagination.NewLister(a, recv, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, cond.And()).Paging(a.Context())
	return err
}

func (a *NgingOauthAgree) BatchValidate(kvset map[string]interface{}) error {
	if kvset == nil {
		kvset = a.AsRow()
	}
	return DBI.Fields.BatchValidate(a.Short_(), kvset)
}

func (a *NgingOauthAgree) Validate(field string, value interface{}) error {
	return DBI.Fields.Validate(a.Short_(), field, value)
}
