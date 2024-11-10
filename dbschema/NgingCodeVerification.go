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

type Slice_NgingCodeVerification []*NgingCodeVerification

func (s Slice_NgingCodeVerification) Range(fn func(m factory.Model) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_NgingCodeVerification) RangeRaw(fn func(m *NgingCodeVerification) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_NgingCodeVerification) GroupBy(keyField string) map[string][]*NgingCodeVerification {
	r := map[string][]*NgingCodeVerification{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		if _, y := r[vkey]; !y {
			r[vkey] = []*NgingCodeVerification{}
		}
		r[vkey] = append(r[vkey], row)
	}
	return r
}

func (s Slice_NgingCodeVerification) KeyBy(keyField string) map[string]*NgingCodeVerification {
	r := map[string]*NgingCodeVerification{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = row
	}
	return r
}

func (s Slice_NgingCodeVerification) AsKV(keyField string, valueField string) param.Store {
	r := param.Store{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = dmap[valueField]
	}
	return r
}

func (s Slice_NgingCodeVerification) Transform(transfers map[string]param.Transfer) []param.Store {
	r := make([]param.Store, len(s))
	for idx, row := range s {
		r[idx] = row.AsMap().Transform(transfers)
	}
	return r
}

func (s Slice_NgingCodeVerification) FromList(data interface{}) Slice_NgingCodeVerification {
	values, ok := data.([]*NgingCodeVerification)
	if !ok {
		for _, value := range data.([]interface{}) {
			row := &NgingCodeVerification{}
			row.FromRow(value.(map[string]interface{}))
			s = append(s, row)
		}
		return s
	}
	s = append(s, values...)

	return s
}

func NewNgingCodeVerification(ctx echo.Context) *NgingCodeVerification {
	m := &NgingCodeVerification{}
	m.SetContext(ctx)
	return m
}

// NgingCodeVerification 验证码
type NgingCodeVerification struct {
	base    factory.Base
	objects []*NgingCodeVerification

	Id         uint64 `db:"id,omitempty,pk" bson:"id,omitempty" comment:"ID" json:"id" xml:"id"`
	Code       string `db:"code" bson:"code" comment:"验证码" json:"code" xml:"code"`
	Created    uint   `db:"created" bson:"created" comment:"创建时间" json:"created" xml:"created"`
	OwnerId    uint64 `db:"owner_id" bson:"owner_id" comment:"所有者ID" json:"owner_id" xml:"owner_id"`
	OwnerType  string `db:"owner_type" bson:"owner_type" comment:"所有者类型" json:"owner_type" xml:"owner_type"`
	Used       uint   `db:"used" bson:"used" comment:"使用时间" json:"used" xml:"used"`
	Purpose    string `db:"purpose" bson:"purpose" comment:"目的" json:"purpose" xml:"purpose"`
	Start      uint   `db:"start" bson:"start" comment:"有效时间" json:"start" xml:"start"`
	End        uint   `db:"end" bson:"end" comment:"失效时间" json:"end" xml:"end"`
	Disabled   string `db:"disabled" bson:"disabled" comment:"是否禁用" json:"disabled" xml:"disabled"`
	SendMethod string `db:"send_method" bson:"send_method" comment:"发送方式(mobile-手机;email-邮箱)" json:"send_method" xml:"send_method"`
	SendTo     string `db:"send_to" bson:"send_to" comment:"发送目标" json:"send_to" xml:"send_to"`
}

// - base function

func (a *NgingCodeVerification) Trans() factory.Transactioner {
	return a.base.Trans()
}

func (a *NgingCodeVerification) Use(trans factory.Transactioner) factory.Model {
	a.base.Use(trans)
	return a
}

func (a *NgingCodeVerification) SetContext(ctx echo.Context) factory.Model {
	a.base.SetContext(ctx)
	return a
}

func (a *NgingCodeVerification) EventON(on ...bool) factory.Model {
	a.base.EventON(on...)
	return a
}

func (a *NgingCodeVerification) EventOFF(off ...bool) factory.Model {
	a.base.EventOFF(off...)
	return a
}

func (a *NgingCodeVerification) Context() echo.Context {
	return a.base.Context()
}

func (a *NgingCodeVerification) SetConnID(connID int) factory.Model {
	a.base.SetConnID(connID)
	return a
}

func (a *NgingCodeVerification) ConnID() int {
	return a.base.ConnID()
}

func (a *NgingCodeVerification) SetNamer(namer func(factory.Model) string) factory.Model {
	a.base.SetNamer(namer)
	return a
}

func (a *NgingCodeVerification) Namer() func(factory.Model) string {
	return a.base.Namer()
}

func (a *NgingCodeVerification) SetParam(param *factory.Param) factory.Model {
	a.base.SetParam(param)
	return a
}

func (a *NgingCodeVerification) Param(mw func(db.Result) db.Result, args ...interface{}) *factory.Param {
	if a.base.Param() == nil {
		return a.NewParam().SetMiddleware(mw).SetArgs(args...)
	}
	return a.base.Param().SetMiddleware(mw).SetArgs(args...)
}

func (a *NgingCodeVerification) New(structName string, connID ...int) factory.Model {
	return a.base.New(structName, connID...)
}

func (a *NgingCodeVerification) Base_() factory.Baser {
	return &a.base
}

// - current function

func (a *NgingCodeVerification) Objects() []*NgingCodeVerification {
	if a.objects == nil {
		return nil
	}
	return a.objects[:]
}

func (a *NgingCodeVerification) XObjects() Slice_NgingCodeVerification {
	return Slice_NgingCodeVerification(a.Objects())
}

func (a *NgingCodeVerification) NewObjects() factory.Ranger {
	return &Slice_NgingCodeVerification{}
}

func (a *NgingCodeVerification) InitObjects() *[]*NgingCodeVerification {
	a.objects = []*NgingCodeVerification{}
	return &a.objects
}

func (a *NgingCodeVerification) NewParam() *factory.Param {
	return factory.NewParam(factory.DefaultFactory).SetIndex(a.base.ConnID()).SetTrans(a.base.Trans()).SetCollection(a.Name_()).SetModel(a)
}

func (a *NgingCodeVerification) Short_() string {
	return "nging_code_verification"
}

func (a *NgingCodeVerification) Struct_() string {
	return "NgingCodeVerification"
}

func (a *NgingCodeVerification) Name_() string {
	b := a
	if b == nil {
		b = &NgingCodeVerification{}
	}
	if b.base.Namer() != nil {
		return WithPrefix(b.base.Namer()(b))
	}
	return WithPrefix(factory.TableNamerGet(b.Short_())(b))
}

func (a *NgingCodeVerification) CPAFrom(source factory.Model) factory.Model {
	a.SetContext(source.Context())
	a.SetConnID(source.ConnID())
	a.SetNamer(source.Namer())
	return a
}

func (a *NgingCodeVerification) Get(mw func(db.Result) db.Result, args ...interface{}) (err error) {
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

func (a *NgingCodeVerification) List(recv interface{}, mw func(db.Result) db.Result, page, size int, args ...interface{}) (func() int64, error) {
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
		case *[]*NgingCodeVerification:
			err = DBI.FireReaded(a, queryParam, Slice_NgingCodeVerification(*v))
		case []*NgingCodeVerification:
			err = DBI.FireReaded(a, queryParam, Slice_NgingCodeVerification(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *NgingCodeVerification) GroupBy(keyField string, inputRows ...[]*NgingCodeVerification) map[string][]*NgingCodeVerification {
	var rows Slice_NgingCodeVerification
	if len(inputRows) > 0 {
		rows = Slice_NgingCodeVerification(inputRows[0])
	} else {
		rows = Slice_NgingCodeVerification(a.Objects())
	}
	return rows.GroupBy(keyField)
}

func (a *NgingCodeVerification) KeyBy(keyField string, inputRows ...[]*NgingCodeVerification) map[string]*NgingCodeVerification {
	var rows Slice_NgingCodeVerification
	if len(inputRows) > 0 {
		rows = Slice_NgingCodeVerification(inputRows[0])
	} else {
		rows = Slice_NgingCodeVerification(a.Objects())
	}
	return rows.KeyBy(keyField)
}

func (a *NgingCodeVerification) AsKV(keyField string, valueField string, inputRows ...[]*NgingCodeVerification) param.Store {
	var rows Slice_NgingCodeVerification
	if len(inputRows) > 0 {
		rows = Slice_NgingCodeVerification(inputRows[0])
	} else {
		rows = Slice_NgingCodeVerification(a.Objects())
	}
	return rows.AsKV(keyField, valueField)
}

func (a *NgingCodeVerification) ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error) {
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
		case *[]*NgingCodeVerification:
			err = DBI.FireReaded(a, queryParam, Slice_NgingCodeVerification(*v))
		case []*NgingCodeVerification:
			err = DBI.FireReaded(a, queryParam, Slice_NgingCodeVerification(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *NgingCodeVerification) Insert() (pk interface{}, err error) {
	a.Created = uint(time.Now().Unix())
	a.Id = 0
	if len(a.OwnerType) == 0 {
		a.OwnerType = "user"
	}
	if len(a.Disabled) == 0 {
		a.Disabled = "N"
	}
	if len(a.SendMethod) == 0 {
		a.SendMethod = "mobile"
	}
	if a.base.Eventable() {
		err = DBI.Fire("creating", a, nil)
		if err != nil {
			return
		}
	}
	pk, err = a.Param(nil).SetSend(a).Insert()
	if err == nil && pk != nil {
		if v, y := pk.(uint64); y {
			a.Id = v
		} else if v, y := pk.(int64); y {
			a.Id = uint64(v)
		}
	}
	if err == nil && a.base.Eventable() {
		err = DBI.Fire("created", a, nil)
	}
	return
}

func (a *NgingCodeVerification) Update(mw func(db.Result) db.Result, args ...interface{}) (err error) {

	if len(a.OwnerType) == 0 {
		a.OwnerType = "user"
	}
	if len(a.Disabled) == 0 {
		a.Disabled = "N"
	}
	if len(a.SendMethod) == 0 {
		a.SendMethod = "mobile"
	}
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

func (a *NgingCodeVerification) Updatex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {

	if len(a.OwnerType) == 0 {
		a.OwnerType = "user"
	}
	if len(a.Disabled) == 0 {
		a.Disabled = "N"
	}
	if len(a.SendMethod) == 0 {
		a.SendMethod = "mobile"
	}
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

func (a *NgingCodeVerification) UpdateByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (err error) {

	if len(a.OwnerType) == 0 {
		a.OwnerType = "user"
	}
	if len(a.Disabled) == 0 {
		a.Disabled = "N"
	}
	if len(a.SendMethod) == 0 {
		a.SendMethod = "mobile"
	}
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

func (a *NgingCodeVerification) UpdatexByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (affected int64, err error) {

	if len(a.OwnerType) == 0 {
		a.OwnerType = "user"
	}
	if len(a.Disabled) == 0 {
		a.Disabled = "N"
	}
	if len(a.SendMethod) == 0 {
		a.SendMethod = "mobile"
	}
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

func (a *NgingCodeVerification) UpdateField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (err error) {
	return a.UpdateFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *NgingCodeVerification) UpdatexField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (affected int64, err error) {
	return a.UpdatexFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *NgingCodeVerification) UpdateFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (err error) {

	if val, ok := kvset["owner_type"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["owner_type"] = "user"
		}
	}
	if val, ok := kvset["disabled"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["disabled"] = "N"
		}
	}
	if val, ok := kvset["send_method"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["send_method"] = "mobile"
		}
	}
	if !a.base.Eventable() {
		return a.Param(mw, args...).SetSend(kvset).Update()
	}
	m := *a
	m.FromRow(kvset)
	editColumns := make([]string, 0, len(kvset))
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

func (a *NgingCodeVerification) UpdatexFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (affected int64, err error) {

	if val, ok := kvset["owner_type"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["owner_type"] = "user"
		}
	}
	if val, ok := kvset["disabled"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["disabled"] = "N"
		}
	}
	if val, ok := kvset["send_method"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["send_method"] = "mobile"
		}
	}
	if !a.base.Eventable() {
		return a.Param(mw, args...).SetSend(kvset).Updatex()
	}
	m := *a
	m.FromRow(kvset)
	editColumns := make([]string, 0, len(kvset))
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

func (a *NgingCodeVerification) UpdateValues(mw func(db.Result) db.Result, keysValues *db.KeysValues, args ...interface{}) (err error) {
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

func (a *NgingCodeVerification) Upsert(mw func(db.Result) db.Result, args ...interface{}) (pk interface{}, err error) {
	pk, err = a.Param(mw, args...).SetSend(a).Upsert(func() error {
		if len(a.OwnerType) == 0 {
			a.OwnerType = "user"
		}
		if len(a.Disabled) == 0 {
			a.Disabled = "N"
		}
		if len(a.SendMethod) == 0 {
			a.SendMethod = "mobile"
		}
		if !a.base.Eventable() {
			return nil
		}
		return DBI.Fire("updating", a, mw, args...)
	}, func() error {
		a.Created = uint(time.Now().Unix())
		a.Id = 0
		if len(a.OwnerType) == 0 {
			a.OwnerType = "user"
		}
		if len(a.Disabled) == 0 {
			a.Disabled = "N"
		}
		if len(a.SendMethod) == 0 {
			a.SendMethod = "mobile"
		}
		if !a.base.Eventable() {
			return nil
		}
		return DBI.Fire("creating", a, nil)
	})
	if err == nil && pk != nil {
		if v, y := pk.(uint64); y {
			a.Id = v
		} else if v, y := pk.(int64); y {
			a.Id = uint64(v)
		}
	}
	if err == nil && a.base.Eventable() {
		if pk == nil {
			err = DBI.Fire("updated", a, mw, args...)
		} else {
			err = DBI.Fire("created", a, nil)
		}
	}
	return
}

func (a *NgingCodeVerification) Delete(mw func(db.Result) db.Result, args ...interface{}) (err error) {

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

func (a *NgingCodeVerification) Deletex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {

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

func (a *NgingCodeVerification) Count(mw func(db.Result) db.Result, args ...interface{}) (int64, error) {
	return a.Param(mw, args...).Count()
}

func (a *NgingCodeVerification) Exists(mw func(db.Result) db.Result, args ...interface{}) (bool, error) {
	return a.Param(mw, args...).Exists()
}

func (a *NgingCodeVerification) Reset() *NgingCodeVerification {
	a.Id = 0
	a.Code = ``
	a.Created = 0
	a.OwnerId = 0
	a.OwnerType = ``
	a.Used = 0
	a.Purpose = ``
	a.Start = 0
	a.End = 0
	a.Disabled = ``
	a.SendMethod = ``
	a.SendTo = ``
	return a
}

func (a *NgingCodeVerification) AsMap(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["Id"] = a.Id
		r["Code"] = a.Code
		r["Created"] = a.Created
		r["OwnerId"] = a.OwnerId
		r["OwnerType"] = a.OwnerType
		r["Used"] = a.Used
		r["Purpose"] = a.Purpose
		r["Start"] = a.Start
		r["End"] = a.End
		r["Disabled"] = a.Disabled
		r["SendMethod"] = a.SendMethod
		r["SendTo"] = a.SendTo
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "Id":
			r["Id"] = a.Id
		case "Code":
			r["Code"] = a.Code
		case "Created":
			r["Created"] = a.Created
		case "OwnerId":
			r["OwnerId"] = a.OwnerId
		case "OwnerType":
			r["OwnerType"] = a.OwnerType
		case "Used":
			r["Used"] = a.Used
		case "Purpose":
			r["Purpose"] = a.Purpose
		case "Start":
			r["Start"] = a.Start
		case "End":
			r["End"] = a.End
		case "Disabled":
			r["Disabled"] = a.Disabled
		case "SendMethod":
			r["SendMethod"] = a.SendMethod
		case "SendTo":
			r["SendTo"] = a.SendTo
		}
	}
	return r
}

func (a *NgingCodeVerification) FromRow(row map[string]interface{}) {
	for key, value := range row {
		if _, ok := value.(db.RawValue); ok {
			continue
		}
		switch key {
		case "id":
			a.Id = param.AsUint64(value)
		case "code":
			a.Code = param.AsString(value)
		case "created":
			a.Created = param.AsUint(value)
		case "owner_id":
			a.OwnerId = param.AsUint64(value)
		case "owner_type":
			a.OwnerType = param.AsString(value)
		case "used":
			a.Used = param.AsUint(value)
		case "purpose":
			a.Purpose = param.AsString(value)
		case "start":
			a.Start = param.AsUint(value)
		case "end":
			a.End = param.AsUint(value)
		case "disabled":
			a.Disabled = param.AsString(value)
		case "send_method":
			a.SendMethod = param.AsString(value)
		case "send_to":
			a.SendTo = param.AsString(value)
		}
	}
}

func (a *NgingCodeVerification) GetField(field string) interface{} {
	switch field {
	case "Id":
		return a.Id
	case "Code":
		return a.Code
	case "Created":
		return a.Created
	case "OwnerId":
		return a.OwnerId
	case "OwnerType":
		return a.OwnerType
	case "Used":
		return a.Used
	case "Purpose":
		return a.Purpose
	case "Start":
		return a.Start
	case "End":
		return a.End
	case "Disabled":
		return a.Disabled
	case "SendMethod":
		return a.SendMethod
	case "SendTo":
		return a.SendTo
	default:
		return nil
	}
}

func (a *NgingCodeVerification) GetAllFieldNames() []string {
	return []string{
		"Id",
		"Code",
		"Created",
		"OwnerId",
		"OwnerType",
		"Used",
		"Purpose",
		"Start",
		"End",
		"Disabled",
		"SendMethod",
		"SendTo",
	}
}

func (a *NgingCodeVerification) HasField(field string) bool {
	switch field {
	case "Id":
		return true
	case "Code":
		return true
	case "Created":
		return true
	case "OwnerId":
		return true
	case "OwnerType":
		return true
	case "Used":
		return true
	case "Purpose":
		return true
	case "Start":
		return true
	case "End":
		return true
	case "Disabled":
		return true
	case "SendMethod":
		return true
	case "SendTo":
		return true
	default:
		return false
	}
}

func (a *NgingCodeVerification) Set(key interface{}, value ...interface{}) {
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
		case "Id":
			a.Id = param.AsUint64(vv)
		case "Code":
			a.Code = param.AsString(vv)
		case "Created":
			a.Created = param.AsUint(vv)
		case "OwnerId":
			a.OwnerId = param.AsUint64(vv)
		case "OwnerType":
			a.OwnerType = param.AsString(vv)
		case "Used":
			a.Used = param.AsUint(vv)
		case "Purpose":
			a.Purpose = param.AsString(vv)
		case "Start":
			a.Start = param.AsUint(vv)
		case "End":
			a.End = param.AsUint(vv)
		case "Disabled":
			a.Disabled = param.AsString(vv)
		case "SendMethod":
			a.SendMethod = param.AsString(vv)
		case "SendTo":
			a.SendTo = param.AsString(vv)
		}
	}
}

func (a *NgingCodeVerification) AsRow(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["id"] = a.Id
		r["code"] = a.Code
		r["created"] = a.Created
		r["owner_id"] = a.OwnerId
		r["owner_type"] = a.OwnerType
		r["used"] = a.Used
		r["purpose"] = a.Purpose
		r["start"] = a.Start
		r["end"] = a.End
		r["disabled"] = a.Disabled
		r["send_method"] = a.SendMethod
		r["send_to"] = a.SendTo
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "id":
			r["id"] = a.Id
		case "code":
			r["code"] = a.Code
		case "created":
			r["created"] = a.Created
		case "owner_id":
			r["owner_id"] = a.OwnerId
		case "owner_type":
			r["owner_type"] = a.OwnerType
		case "used":
			r["used"] = a.Used
		case "purpose":
			r["purpose"] = a.Purpose
		case "start":
			r["start"] = a.Start
		case "end":
			r["end"] = a.End
		case "disabled":
			r["disabled"] = a.Disabled
		case "send_method":
			r["send_method"] = a.SendMethod
		case "send_to":
			r["send_to"] = a.SendTo
		}
	}
	return r
}

func (a *NgingCodeVerification) ListPage(cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPage(a, cond, sorts...)
}

func (a *NgingCodeVerification) ListPageAs(recv interface{}, cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPageAs(a, recv, cond, sorts...)
}

func (a *NgingCodeVerification) ListPageByOffset(cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPageByOffset(a, cond, sorts...)
}

func (a *NgingCodeVerification) ListPageByOffsetAs(recv interface{}, cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPageByOffsetAs(a, recv, cond, sorts...)
}

func (a *NgingCodeVerification) BatchValidate(kvset map[string]interface{}) error {
	if kvset == nil {
		kvset = a.AsRow()
	}
	return DBI.Fields.BatchValidate(a.Short_(), kvset)
}

func (a *NgingCodeVerification) Validate(field string, value interface{}) error {
	return DBI.Fields.Validate(a.Short_(), field, value)
}
