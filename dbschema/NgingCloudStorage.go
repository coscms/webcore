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

type Slice_NgingCloudStorage []*NgingCloudStorage

func (s Slice_NgingCloudStorage) Range(fn func(m factory.Model) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_NgingCloudStorage) RangeRaw(fn func(m *NgingCloudStorage) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_NgingCloudStorage) GroupBy(keyField string) map[string][]*NgingCloudStorage {
	r := map[string][]*NgingCloudStorage{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		if _, y := r[vkey]; !y {
			r[vkey] = []*NgingCloudStorage{}
		}
		r[vkey] = append(r[vkey], row)
	}
	return r
}

func (s Slice_NgingCloudStorage) KeyBy(keyField string) map[string]*NgingCloudStorage {
	r := map[string]*NgingCloudStorage{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = row
	}
	return r
}

func (s Slice_NgingCloudStorage) AsKV(keyField string, valueField string) param.Store {
	r := param.Store{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = dmap[valueField]
	}
	return r
}

func (s Slice_NgingCloudStorage) Transform(transfers map[string]param.Transfer) []param.Store {
	r := make([]param.Store, len(s))
	for idx, row := range s {
		r[idx] = row.AsMap().Transform(transfers)
	}
	return r
}

func (s Slice_NgingCloudStorage) FromList(data interface{}) Slice_NgingCloudStorage {
	values, ok := data.([]*NgingCloudStorage)
	if !ok {
		for _, value := range data.([]interface{}) {
			row := &NgingCloudStorage{}
			row.FromRow(value.(map[string]interface{}))
			s = append(s, row)
		}
		return s
	}
	s = append(s, values...)

	return s
}

func NewNgingCloudStorage(ctx echo.Context) *NgingCloudStorage {
	m := &NgingCloudStorage{}
	m.SetContext(ctx)
	return m
}

// NgingCloudStorage 云存储账号
type NgingCloudStorage struct {
	base    factory.Base
	objects []*NgingCloudStorage

	Id       uint   `db:"id,omitempty,pk" bson:"id,omitempty" comment:"ID" json:"id" xml:"id"`
	Name     string `db:"name" bson:"name" comment:"名称" json:"name" xml:"name"`
	Type     string `db:"type" bson:"type" comment:"存储类型(aws,oss,cos)" json:"type" xml:"type"`
	Key      string `db:"key" bson:"key" comment:"key" json:"key" xml:"key"`
	Secret   string `db:"secret" bson:"secret" comment:"密钥(加密处理)" json:"secret" xml:"secret"`
	Bucket   string `db:"bucket" bson:"bucket" comment:"存储桶" json:"bucket" xml:"bucket"`
	Endpoint string `db:"endpoint" bson:"endpoint" comment:"地域节点" json:"endpoint" xml:"endpoint"`
	Region   string `db:"region" bson:"region" comment:"地区" json:"region" xml:"region"`
	Secure   string `db:"secure" bson:"secure" comment:"是否(Y/N)HTTPS" json:"secure" xml:"secure"`
	Baseurl  string `db:"baseurl" bson:"baseurl" comment:"资源基础网址" json:"baseurl" xml:"baseurl"`
	Created  uint   `db:"created" bson:"created" comment:"创建时间" json:"created" xml:"created"`
	Updated  uint   `db:"updated" bson:"updated" comment:"修改时间" json:"updated" xml:"updated"`
}

// - base function

func (a *NgingCloudStorage) Trans() factory.Transactioner {
	return a.base.Trans()
}

func (a *NgingCloudStorage) Use(trans factory.Transactioner) factory.Model {
	a.base.Use(trans)
	return a
}

func (a *NgingCloudStorage) SetContext(ctx echo.Context) factory.Model {
	a.base.SetContext(ctx)
	return a
}

func (a *NgingCloudStorage) EventON(on ...bool) factory.Model {
	a.base.EventON(on...)
	return a
}

func (a *NgingCloudStorage) EventOFF(off ...bool) factory.Model {
	a.base.EventOFF(off...)
	return a
}

func (a *NgingCloudStorage) Context() echo.Context {
	return a.base.Context()
}

func (a *NgingCloudStorage) SetConnID(connID int) factory.Model {
	a.base.SetConnID(connID)
	return a
}

func (a *NgingCloudStorage) ConnID() int {
	return a.base.ConnID()
}

func (a *NgingCloudStorage) SetNamer(namer func(factory.Model) string) factory.Model {
	a.base.SetNamer(namer)
	return a
}

func (a *NgingCloudStorage) Namer() func(factory.Model) string {
	return a.base.Namer()
}

func (a *NgingCloudStorage) SetParam(param *factory.Param) factory.Model {
	a.base.SetParam(param)
	return a
}

func (a *NgingCloudStorage) Param(mw func(db.Result) db.Result, args ...interface{}) *factory.Param {
	if a.base.Param() == nil {
		return a.NewParam().SetMiddleware(mw).SetArgs(args...)
	}
	return a.base.Param().SetMiddleware(mw).SetArgs(args...)
}

func (a *NgingCloudStorage) New(structName string, connID ...int) factory.Model {
	return a.base.New(structName, connID...)
}

func (a *NgingCloudStorage) Base_() factory.Baser {
	return &a.base
}

// - current function

func (a *NgingCloudStorage) Objects() []*NgingCloudStorage {
	if a.objects == nil {
		return nil
	}
	return a.objects[:]
}

func (a *NgingCloudStorage) XObjects() Slice_NgingCloudStorage {
	return Slice_NgingCloudStorage(a.Objects())
}

func (a *NgingCloudStorage) NewObjects() factory.Ranger {
	return &Slice_NgingCloudStorage{}
}

func (a *NgingCloudStorage) InitObjects() *[]*NgingCloudStorage {
	a.objects = []*NgingCloudStorage{}
	return &a.objects
}

func (a *NgingCloudStorage) NewParam() *factory.Param {
	return factory.NewParam(factory.DefaultFactory).SetIndex(a.base.ConnID()).SetTrans(a.base.Trans()).SetCollection(a.Name_()).SetModel(a)
}

func (a *NgingCloudStorage) Short_() string {
	return "nging_cloud_storage"
}

func (a *NgingCloudStorage) Struct_() string {
	return "NgingCloudStorage"
}

func (a *NgingCloudStorage) Name_() string {
	b := a
	if b == nil {
		b = &NgingCloudStorage{}
	}
	if b.base.Namer() != nil {
		return WithPrefix(b.base.Namer()(b))
	}
	return WithPrefix(factory.TableNamerGet(b.Short_())(b))
}

func (a *NgingCloudStorage) CPAFrom(source factory.Model) factory.Model {
	a.SetContext(source.Context())
	a.SetConnID(source.ConnID())
	a.SetNamer(source.Namer())
	return a
}

func (a *NgingCloudStorage) Get(mw func(db.Result) db.Result, args ...interface{}) (err error) {
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

func (a *NgingCloudStorage) List(recv interface{}, mw func(db.Result) db.Result, page, size int, args ...interface{}) (func() int64, error) {
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
		case *[]*NgingCloudStorage:
			err = DBI.FireReaded(a, queryParam, Slice_NgingCloudStorage(*v))
		case []*NgingCloudStorage:
			err = DBI.FireReaded(a, queryParam, Slice_NgingCloudStorage(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *NgingCloudStorage) GroupBy(keyField string, inputRows ...[]*NgingCloudStorage) map[string][]*NgingCloudStorage {
	var rows Slice_NgingCloudStorage
	if len(inputRows) > 0 {
		rows = Slice_NgingCloudStorage(inputRows[0])
	} else {
		rows = Slice_NgingCloudStorage(a.Objects())
	}
	return rows.GroupBy(keyField)
}

func (a *NgingCloudStorage) KeyBy(keyField string, inputRows ...[]*NgingCloudStorage) map[string]*NgingCloudStorage {
	var rows Slice_NgingCloudStorage
	if len(inputRows) > 0 {
		rows = Slice_NgingCloudStorage(inputRows[0])
	} else {
		rows = Slice_NgingCloudStorage(a.Objects())
	}
	return rows.KeyBy(keyField)
}

func (a *NgingCloudStorage) AsKV(keyField string, valueField string, inputRows ...[]*NgingCloudStorage) param.Store {
	var rows Slice_NgingCloudStorage
	if len(inputRows) > 0 {
		rows = Slice_NgingCloudStorage(inputRows[0])
	} else {
		rows = Slice_NgingCloudStorage(a.Objects())
	}
	return rows.AsKV(keyField, valueField)
}

func (a *NgingCloudStorage) ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error) {
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
		case *[]*NgingCloudStorage:
			err = DBI.FireReaded(a, queryParam, Slice_NgingCloudStorage(*v))
		case []*NgingCloudStorage:
			err = DBI.FireReaded(a, queryParam, Slice_NgingCloudStorage(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *NgingCloudStorage) Insert() (pk interface{}, err error) {
	a.Created = uint(time.Now().Unix())
	a.Id = 0
	if len(a.Type) == 0 {
		a.Type = "aws"
	}
	if len(a.Secure) == 0 {
		a.Secure = "Y"
	}
	if a.base.Eventable() {
		err = DBI.Fire("creating", a, nil)
		if err != nil {
			return
		}
	}
	pk, err = a.Param(nil).SetSend(a).Insert()
	if err == nil && pk != nil {
		if v, y := pk.(uint); y {
			a.Id = v
		} else if v, y := pk.(int64); y {
			a.Id = uint(v)
		}
	}
	if err == nil && a.base.Eventable() {
		err = DBI.Fire("created", a, nil)
	}
	return
}

func (a *NgingCloudStorage) Update(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	a.Updated = uint(time.Now().Unix())
	if len(a.Type) == 0 {
		a.Type = "aws"
	}
	if len(a.Secure) == 0 {
		a.Secure = "Y"
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

func (a *NgingCloudStorage) Updatex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {
	a.Updated = uint(time.Now().Unix())
	if len(a.Type) == 0 {
		a.Type = "aws"
	}
	if len(a.Secure) == 0 {
		a.Secure = "Y"
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

func (a *NgingCloudStorage) UpdateByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (err error) {
	a.Updated = uint(time.Now().Unix())
	if len(a.Type) == 0 {
		a.Type = "aws"
	}
	if len(a.Secure) == 0 {
		a.Secure = "Y"
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

func (a *NgingCloudStorage) UpdatexByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (affected int64, err error) {
	a.Updated = uint(time.Now().Unix())
	if len(a.Type) == 0 {
		a.Type = "aws"
	}
	if len(a.Secure) == 0 {
		a.Secure = "Y"
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

func (a *NgingCloudStorage) UpdateField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (err error) {
	return a.UpdateFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *NgingCloudStorage) UpdatexField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (affected int64, err error) {
	return a.UpdatexFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *NgingCloudStorage) UpdateFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (err error) {

	if val, ok := kvset["type"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["type"] = "aws"
		}
	}
	if val, ok := kvset["secure"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["secure"] = "Y"
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

func (a *NgingCloudStorage) UpdatexFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (affected int64, err error) {

	if val, ok := kvset["type"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["type"] = "aws"
		}
	}
	if val, ok := kvset["secure"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["secure"] = "Y"
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

func (a *NgingCloudStorage) UpdateValues(mw func(db.Result) db.Result, keysValues *db.KeysValues, args ...interface{}) (err error) {
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

func (a *NgingCloudStorage) Upsert(mw func(db.Result) db.Result, args ...interface{}) (pk interface{}, err error) {
	pk, err = a.Param(mw, args...).SetSend(a).Upsert(func() error {
		a.Updated = uint(time.Now().Unix())
		if len(a.Type) == 0 {
			a.Type = "aws"
		}
		if len(a.Secure) == 0 {
			a.Secure = "Y"
		}
		if !a.base.Eventable() {
			return nil
		}
		return DBI.Fire("updating", a, mw, args...)
	}, func() error {
		a.Created = uint(time.Now().Unix())
		a.Id = 0
		if len(a.Type) == 0 {
			a.Type = "aws"
		}
		if len(a.Secure) == 0 {
			a.Secure = "Y"
		}
		if !a.base.Eventable() {
			return nil
		}
		return DBI.Fire("creating", a, nil)
	})
	if err == nil && pk != nil {
		if v, y := pk.(uint); y {
			a.Id = v
		} else if v, y := pk.(int64); y {
			a.Id = uint(v)
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

func (a *NgingCloudStorage) Delete(mw func(db.Result) db.Result, args ...interface{}) (err error) {

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

func (a *NgingCloudStorage) Deletex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {

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

func (a *NgingCloudStorage) Count(mw func(db.Result) db.Result, args ...interface{}) (int64, error) {
	return a.Param(mw, args...).Count()
}

func (a *NgingCloudStorage) Exists(mw func(db.Result) db.Result, args ...interface{}) (bool, error) {
	return a.Param(mw, args...).Exists()
}

func (a *NgingCloudStorage) Reset() *NgingCloudStorage {
	a.Id = 0
	a.Name = ``
	a.Type = ``
	a.Key = ``
	a.Secret = ``
	a.Bucket = ``
	a.Endpoint = ``
	a.Region = ``
	a.Secure = ``
	a.Baseurl = ``
	a.Created = 0
	a.Updated = 0
	return a
}

func (a *NgingCloudStorage) AsMap(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["Id"] = a.Id
		r["Name"] = a.Name
		r["Type"] = a.Type
		r["Key"] = a.Key
		r["Secret"] = a.Secret
		r["Bucket"] = a.Bucket
		r["Endpoint"] = a.Endpoint
		r["Region"] = a.Region
		r["Secure"] = a.Secure
		r["Baseurl"] = a.Baseurl
		r["Created"] = a.Created
		r["Updated"] = a.Updated
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "Id":
			r["Id"] = a.Id
		case "Name":
			r["Name"] = a.Name
		case "Type":
			r["Type"] = a.Type
		case "Key":
			r["Key"] = a.Key
		case "Secret":
			r["Secret"] = a.Secret
		case "Bucket":
			r["Bucket"] = a.Bucket
		case "Endpoint":
			r["Endpoint"] = a.Endpoint
		case "Region":
			r["Region"] = a.Region
		case "Secure":
			r["Secure"] = a.Secure
		case "Baseurl":
			r["Baseurl"] = a.Baseurl
		case "Created":
			r["Created"] = a.Created
		case "Updated":
			r["Updated"] = a.Updated
		}
	}
	return r
}

func (a *NgingCloudStorage) FromRow(row map[string]interface{}) {
	for key, value := range row {
		if _, ok := value.(db.RawValue); ok {
			continue
		}
		switch key {
		case "id":
			a.Id = param.AsUint(value)
		case "name":
			a.Name = param.AsString(value)
		case "type":
			a.Type = param.AsString(value)
		case "key":
			a.Key = param.AsString(value)
		case "secret":
			a.Secret = param.AsString(value)
		case "bucket":
			a.Bucket = param.AsString(value)
		case "endpoint":
			a.Endpoint = param.AsString(value)
		case "region":
			a.Region = param.AsString(value)
		case "secure":
			a.Secure = param.AsString(value)
		case "baseurl":
			a.Baseurl = param.AsString(value)
		case "created":
			a.Created = param.AsUint(value)
		case "updated":
			a.Updated = param.AsUint(value)
		}
	}
}

func (a *NgingCloudStorage) GetField(field string) interface{} {
	switch field {
	case "Id":
		return a.Id
	case "Name":
		return a.Name
	case "Type":
		return a.Type
	case "Key":
		return a.Key
	case "Secret":
		return a.Secret
	case "Bucket":
		return a.Bucket
	case "Endpoint":
		return a.Endpoint
	case "Region":
		return a.Region
	case "Secure":
		return a.Secure
	case "Baseurl":
		return a.Baseurl
	case "Created":
		return a.Created
	case "Updated":
		return a.Updated
	default:
		return nil
	}
}

func (a *NgingCloudStorage) GetAllFieldNames() []string {
	return []string{
		"Id",
		"Name",
		"Type",
		"Key",
		"Secret",
		"Bucket",
		"Endpoint",
		"Region",
		"Secure",
		"Baseurl",
		"Created",
		"Updated",
	}
}

func (a *NgingCloudStorage) HasField(field string) bool {
	switch field {
	case "Id":
		return true
	case "Name":
		return true
	case "Type":
		return true
	case "Key":
		return true
	case "Secret":
		return true
	case "Bucket":
		return true
	case "Endpoint":
		return true
	case "Region":
		return true
	case "Secure":
		return true
	case "Baseurl":
		return true
	case "Created":
		return true
	case "Updated":
		return true
	default:
		return false
	}
}

func (a *NgingCloudStorage) Set(key interface{}, value ...interface{}) {
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
			a.Id = param.AsUint(vv)
		case "Name":
			a.Name = param.AsString(vv)
		case "Type":
			a.Type = param.AsString(vv)
		case "Key":
			a.Key = param.AsString(vv)
		case "Secret":
			a.Secret = param.AsString(vv)
		case "Bucket":
			a.Bucket = param.AsString(vv)
		case "Endpoint":
			a.Endpoint = param.AsString(vv)
		case "Region":
			a.Region = param.AsString(vv)
		case "Secure":
			a.Secure = param.AsString(vv)
		case "Baseurl":
			a.Baseurl = param.AsString(vv)
		case "Created":
			a.Created = param.AsUint(vv)
		case "Updated":
			a.Updated = param.AsUint(vv)
		}
	}
}

func (a *NgingCloudStorage) AsRow(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["id"] = a.Id
		r["name"] = a.Name
		r["type"] = a.Type
		r["key"] = a.Key
		r["secret"] = a.Secret
		r["bucket"] = a.Bucket
		r["endpoint"] = a.Endpoint
		r["region"] = a.Region
		r["secure"] = a.Secure
		r["baseurl"] = a.Baseurl
		r["created"] = a.Created
		r["updated"] = a.Updated
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "id":
			r["id"] = a.Id
		case "name":
			r["name"] = a.Name
		case "type":
			r["type"] = a.Type
		case "key":
			r["key"] = a.Key
		case "secret":
			r["secret"] = a.Secret
		case "bucket":
			r["bucket"] = a.Bucket
		case "endpoint":
			r["endpoint"] = a.Endpoint
		case "region":
			r["region"] = a.Region
		case "secure":
			r["secure"] = a.Secure
		case "baseurl":
			r["baseurl"] = a.Baseurl
		case "created":
			r["created"] = a.Created
		case "updated":
			r["updated"] = a.Updated
		}
	}
	return r
}

func (a *NgingCloudStorage) ListPage(cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPage(a, cond, sorts...)
}

func (a *NgingCloudStorage) ListPageAs(recv interface{}, cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPageAs(a, recv, cond, sorts...)
}

func (a *NgingCloudStorage) ListPageByOffset(cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPageByOffset(a, cond, sorts...)
}

func (a *NgingCloudStorage) ListPageByOffsetAs(recv interface{}, cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPageByOffsetAs(a, recv, cond, sorts...)
}

func (a *NgingCloudStorage) BatchValidate(kvset map[string]interface{}) error {
	if kvset == nil {
		kvset = a.AsRow()
	}
	return DBI.Fields.BatchValidate(a.Short_(), kvset)
}

func (a *NgingCloudStorage) Validate(field string, value interface{}) error {
	return DBI.Fields.Validate(a.Short_(), field, value)
}
