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

type Slice_NgingTaskLog []*NgingTaskLog

func (s Slice_NgingTaskLog) Range(fn func(m factory.Model) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_NgingTaskLog) RangeRaw(fn func(m *NgingTaskLog) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_NgingTaskLog) GroupBy(keyField string) map[string][]*NgingTaskLog {
	r := map[string][]*NgingTaskLog{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		if _, y := r[vkey]; !y {
			r[vkey] = []*NgingTaskLog{}
		}
		r[vkey] = append(r[vkey], row)
	}
	return r
}

func (s Slice_NgingTaskLog) KeyBy(keyField string) map[string]*NgingTaskLog {
	r := map[string]*NgingTaskLog{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = row
	}
	return r
}

func (s Slice_NgingTaskLog) AsKV(keyField string, valueField string) param.Store {
	r := param.Store{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = dmap[valueField]
	}
	return r
}

func (s Slice_NgingTaskLog) Transform(transfers map[string]param.Transfer) []param.Store {
	r := make([]param.Store, len(s))
	for idx, row := range s {
		r[idx] = row.AsMap().Transform(transfers)
	}
	return r
}

func (s Slice_NgingTaskLog) FromList(data interface{}) Slice_NgingTaskLog {
	values, ok := data.([]*NgingTaskLog)
	if !ok {
		for _, value := range data.([]interface{}) {
			row := &NgingTaskLog{}
			row.FromRow(value.(map[string]interface{}))
			s = append(s, row)
		}
		return s
	}
	s = append(s, values...)

	return s
}

func NewNgingTaskLog(ctx echo.Context) *NgingTaskLog {
	m := &NgingTaskLog{}
	m.SetContext(ctx)
	return m
}

// NgingTaskLog 任务日志
type NgingTaskLog struct {
	base    factory.Base
	objects []*NgingTaskLog

	Id      uint64 `db:"id,omitempty,pk" bson:"id,omitempty" comment:"" json:"id" xml:"id"`
	TaskId  uint   `db:"task_id" bson:"task_id" comment:"任务ID" json:"task_id" xml:"task_id"`
	Output  string `db:"output" bson:"output" comment:"任务输出" json:"output" xml:"output"`
	Error   string `db:"error" bson:"error" comment:"错误信息" json:"error" xml:"error"`
	Status  string `db:"status" bson:"status" comment:"状态" json:"status" xml:"status"`
	Elapsed uint   `db:"elapsed" bson:"elapsed" comment:"消耗时间(毫秒)" json:"elapsed" xml:"elapsed"`
	Created uint   `db:"created" bson:"created" comment:"创建时间" json:"created" xml:"created"`
}

// - base function

func (a *NgingTaskLog) Trans() factory.Transactioner {
	return a.base.Trans()
}

func (a *NgingTaskLog) Use(trans factory.Transactioner) factory.Model {
	a.base.Use(trans)
	return a
}

func (a *NgingTaskLog) SetContext(ctx echo.Context) factory.Model {
	a.base.SetContext(ctx)
	return a
}

func (a *NgingTaskLog) EventON(on ...bool) factory.Model {
	a.base.EventON(on...)
	return a
}

func (a *NgingTaskLog) EventOFF(off ...bool) factory.Model {
	a.base.EventOFF(off...)
	return a
}

func (a *NgingTaskLog) Context() echo.Context {
	return a.base.Context()
}

func (a *NgingTaskLog) SetConnID(connID int) factory.Model {
	a.base.SetConnID(connID)
	return a
}

func (a *NgingTaskLog) ConnID() int {
	return a.base.ConnID()
}

func (a *NgingTaskLog) SetNamer(namer func(factory.Model) string) factory.Model {
	a.base.SetNamer(namer)
	return a
}

func (a *NgingTaskLog) Namer() func(factory.Model) string {
	return a.base.Namer()
}

func (a *NgingTaskLog) SetParam(param *factory.Param) factory.Model {
	a.base.SetParam(param)
	return a
}

func (a *NgingTaskLog) Param(mw func(db.Result) db.Result, args ...interface{}) *factory.Param {
	if a.base.Param() == nil {
		return a.NewParam().SetMiddleware(mw).SetArgs(args...)
	}
	return a.base.Param().SetMiddleware(mw).SetArgs(args...)
}

func (a *NgingTaskLog) New(structName string, connID ...int) factory.Model {
	return a.base.New(structName, connID...)
}

func (a *NgingTaskLog) Base_() factory.Baser {
	return &a.base
}

// - current function

func (a *NgingTaskLog) Objects() []*NgingTaskLog {
	if a.objects == nil {
		return nil
	}
	return a.objects[:]
}

func (a *NgingTaskLog) XObjects() Slice_NgingTaskLog {
	return Slice_NgingTaskLog(a.Objects())
}

func (a *NgingTaskLog) NewObjects() factory.Ranger {
	return &Slice_NgingTaskLog{}
}

func (a *NgingTaskLog) InitObjects() *[]*NgingTaskLog {
	a.objects = []*NgingTaskLog{}
	return &a.objects
}

func (a *NgingTaskLog) NewParam() *factory.Param {
	return factory.NewParam(factory.DefaultFactory).SetIndex(a.base.ConnID()).SetTrans(a.base.Trans()).SetCollection(a.Name_()).SetModel(a)
}

func (a *NgingTaskLog) Short_() string {
	return "nging_task_log"
}

func (a *NgingTaskLog) Struct_() string {
	return "NgingTaskLog"
}

func (a *NgingTaskLog) Name_() string {
	b := a
	if b == nil {
		b = &NgingTaskLog{}
	}
	if b.base.Namer() != nil {
		return WithPrefix(b.base.Namer()(b))
	}
	return WithPrefix(factory.TableNamerGet(b.Short_())(b))
}

func (a *NgingTaskLog) CPAFrom(source factory.Model) factory.Model {
	a.SetContext(source.Context())
	a.SetConnID(source.ConnID())
	a.SetNamer(source.Namer())
	return a
}

func (a *NgingTaskLog) Get(mw func(db.Result) db.Result, args ...interface{}) (err error) {
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

func (a *NgingTaskLog) List(recv interface{}, mw func(db.Result) db.Result, page, size int, args ...interface{}) (func() int64, error) {
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
		case *[]*NgingTaskLog:
			err = DBI.FireReaded(a, queryParam, Slice_NgingTaskLog(*v))
		case []*NgingTaskLog:
			err = DBI.FireReaded(a, queryParam, Slice_NgingTaskLog(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *NgingTaskLog) GroupBy(keyField string, inputRows ...[]*NgingTaskLog) map[string][]*NgingTaskLog {
	var rows Slice_NgingTaskLog
	if len(inputRows) > 0 {
		rows = Slice_NgingTaskLog(inputRows[0])
	} else {
		rows = Slice_NgingTaskLog(a.Objects())
	}
	return rows.GroupBy(keyField)
}

func (a *NgingTaskLog) KeyBy(keyField string, inputRows ...[]*NgingTaskLog) map[string]*NgingTaskLog {
	var rows Slice_NgingTaskLog
	if len(inputRows) > 0 {
		rows = Slice_NgingTaskLog(inputRows[0])
	} else {
		rows = Slice_NgingTaskLog(a.Objects())
	}
	return rows.KeyBy(keyField)
}

func (a *NgingTaskLog) AsKV(keyField string, valueField string, inputRows ...[]*NgingTaskLog) param.Store {
	var rows Slice_NgingTaskLog
	if len(inputRows) > 0 {
		rows = Slice_NgingTaskLog(inputRows[0])
	} else {
		rows = Slice_NgingTaskLog(a.Objects())
	}
	return rows.AsKV(keyField, valueField)
}

func (a *NgingTaskLog) ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error) {
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
		case *[]*NgingTaskLog:
			err = DBI.FireReaded(a, queryParam, Slice_NgingTaskLog(*v))
		case []*NgingTaskLog:
			err = DBI.FireReaded(a, queryParam, Slice_NgingTaskLog(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *NgingTaskLog) Insert() (pk interface{}, err error) {
	a.Created = uint(time.Now().Unix())
	a.Id = 0
	if len(a.Status) == 0 {
		a.Status = "success"
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

func (a *NgingTaskLog) Update(mw func(db.Result) db.Result, args ...interface{}) (err error) {

	if len(a.Status) == 0 {
		a.Status = "success"
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

func (a *NgingTaskLog) Updatex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {

	if len(a.Status) == 0 {
		a.Status = "success"
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

func (a *NgingTaskLog) UpdateByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (err error) {

	if len(a.Status) == 0 {
		a.Status = "success"
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

func (a *NgingTaskLog) UpdatexByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (affected int64, err error) {

	if len(a.Status) == 0 {
		a.Status = "success"
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

func (a *NgingTaskLog) UpdateField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (err error) {
	return a.UpdateFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *NgingTaskLog) UpdatexField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (affected int64, err error) {
	return a.UpdatexFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *NgingTaskLog) UpdateFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (err error) {

	if val, ok := kvset["status"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["status"] = "success"
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

func (a *NgingTaskLog) UpdatexFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (affected int64, err error) {

	if val, ok := kvset["status"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["status"] = "success"
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

func (a *NgingTaskLog) UpdateValues(mw func(db.Result) db.Result, keysValues *db.KeysValues, args ...interface{}) (err error) {
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

func (a *NgingTaskLog) Upsert(mw func(db.Result) db.Result, args ...interface{}) (pk interface{}, err error) {
	pk, err = a.Param(mw, args...).SetSend(a).Upsert(func() error {
		if len(a.Status) == 0 {
			a.Status = "success"
		}
		if !a.base.Eventable() {
			return nil
		}
		return DBI.Fire("updating", a, mw, args...)
	}, func() error {
		a.Created = uint(time.Now().Unix())
		a.Id = 0
		if len(a.Status) == 0 {
			a.Status = "success"
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

func (a *NgingTaskLog) Delete(mw func(db.Result) db.Result, args ...interface{}) (err error) {

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

func (a *NgingTaskLog) Deletex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {

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

func (a *NgingTaskLog) Count(mw func(db.Result) db.Result, args ...interface{}) (int64, error) {
	return a.Param(mw, args...).Count()
}

func (a *NgingTaskLog) Exists(mw func(db.Result) db.Result, args ...interface{}) (bool, error) {
	return a.Param(mw, args...).Exists()
}

func (a *NgingTaskLog) Reset() *NgingTaskLog {
	a.Id = 0
	a.TaskId = 0
	a.Output = ``
	a.Error = ``
	a.Status = ``
	a.Elapsed = 0
	a.Created = 0
	return a
}

func (a *NgingTaskLog) AsMap(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["Id"] = a.Id
		r["TaskId"] = a.TaskId
		r["Output"] = a.Output
		r["Error"] = a.Error
		r["Status"] = a.Status
		r["Elapsed"] = a.Elapsed
		r["Created"] = a.Created
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "Id":
			r["Id"] = a.Id
		case "TaskId":
			r["TaskId"] = a.TaskId
		case "Output":
			r["Output"] = a.Output
		case "Error":
			r["Error"] = a.Error
		case "Status":
			r["Status"] = a.Status
		case "Elapsed":
			r["Elapsed"] = a.Elapsed
		case "Created":
			r["Created"] = a.Created
		}
	}
	return r
}

func (a *NgingTaskLog) FromRow(row map[string]interface{}) {
	for key, value := range row {
		if _, ok := value.(db.RawValue); ok {
			continue
		}
		switch key {
		case "id":
			a.Id = param.AsUint64(value)
		case "task_id":
			a.TaskId = param.AsUint(value)
		case "output":
			a.Output = param.AsString(value)
		case "error":
			a.Error = param.AsString(value)
		case "status":
			a.Status = param.AsString(value)
		case "elapsed":
			a.Elapsed = param.AsUint(value)
		case "created":
			a.Created = param.AsUint(value)
		}
	}
}

func (a *NgingTaskLog) GetField(field string) interface{} {
	switch field {
	case "Id":
		return a.Id
	case "TaskId":
		return a.TaskId
	case "Output":
		return a.Output
	case "Error":
		return a.Error
	case "Status":
		return a.Status
	case "Elapsed":
		return a.Elapsed
	case "Created":
		return a.Created
	default:
		return nil
	}
}

func (a *NgingTaskLog) GetAllFieldNames() []string {
	return []string{
		"Id",
		"TaskId",
		"Output",
		"Error",
		"Status",
		"Elapsed",
		"Created",
	}
}

func (a *NgingTaskLog) HasField(field string) bool {
	switch field {
	case "Id":
		return true
	case "TaskId":
		return true
	case "Output":
		return true
	case "Error":
		return true
	case "Status":
		return true
	case "Elapsed":
		return true
	case "Created":
		return true
	default:
		return false
	}
}

func (a *NgingTaskLog) Set(key interface{}, value ...interface{}) {
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
		case "TaskId":
			a.TaskId = param.AsUint(vv)
		case "Output":
			a.Output = param.AsString(vv)
		case "Error":
			a.Error = param.AsString(vv)
		case "Status":
			a.Status = param.AsString(vv)
		case "Elapsed":
			a.Elapsed = param.AsUint(vv)
		case "Created":
			a.Created = param.AsUint(vv)
		}
	}
}

func (a *NgingTaskLog) AsRow(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["id"] = a.Id
		r["task_id"] = a.TaskId
		r["output"] = a.Output
		r["error"] = a.Error
		r["status"] = a.Status
		r["elapsed"] = a.Elapsed
		r["created"] = a.Created
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "id":
			r["id"] = a.Id
		case "task_id":
			r["task_id"] = a.TaskId
		case "output":
			r["output"] = a.Output
		case "error":
			r["error"] = a.Error
		case "status":
			r["status"] = a.Status
		case "elapsed":
			r["elapsed"] = a.Elapsed
		case "created":
			r["created"] = a.Created
		}
	}
	return r
}

func (a *NgingTaskLog) ListPage(cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPage(a, cond, sorts...)
}

func (a *NgingTaskLog) ListPageAs(recv interface{}, cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPageAs(a, recv, cond, sorts...)
}

func (a *NgingTaskLog) ListPageByOffset(cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPageByOffset(a, cond, sorts...)
}

func (a *NgingTaskLog) ListPageByOffsetAs(recv interface{}, cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPageByOffsetAs(a, recv, cond, sorts...)
}

func (a *NgingTaskLog) BatchValidate(kvset map[string]interface{}) error {
	if kvset == nil {
		kvset = a.AsRow()
	}
	return DBI.Fields.BatchValidate(a.Short_(), kvset)
}

func (a *NgingTaskLog) Validate(field string, value interface{}) error {
	return DBI.Fields.Validate(a.Short_(), field, value)
}
