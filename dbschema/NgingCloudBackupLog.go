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

type Slice_NgingCloudBackupLog []*NgingCloudBackupLog

func (s Slice_NgingCloudBackupLog) Range(fn func(m factory.Model) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_NgingCloudBackupLog) RangeRaw(fn func(m *NgingCloudBackupLog) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_NgingCloudBackupLog) GroupBy(keyField string) map[string][]*NgingCloudBackupLog {
	r := map[string][]*NgingCloudBackupLog{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		if _, y := r[vkey]; !y {
			r[vkey] = []*NgingCloudBackupLog{}
		}
		r[vkey] = append(r[vkey], row)
	}
	return r
}

func (s Slice_NgingCloudBackupLog) KeyBy(keyField string) map[string]*NgingCloudBackupLog {
	r := map[string]*NgingCloudBackupLog{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = row
	}
	return r
}

func (s Slice_NgingCloudBackupLog) AsKV(keyField string, valueField string) param.Store {
	r := param.Store{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = dmap[valueField]
	}
	return r
}

func (s Slice_NgingCloudBackupLog) Transform(transfers map[string]param.Transfer) []param.Store {
	r := make([]param.Store, len(s))
	for idx, row := range s {
		r[idx] = row.AsMap().Transform(transfers)
	}
	return r
}

func (s Slice_NgingCloudBackupLog) FromList(data interface{}) Slice_NgingCloudBackupLog {
	values, ok := data.([]*NgingCloudBackupLog)
	if !ok {
		for _, value := range data.([]interface{}) {
			row := &NgingCloudBackupLog{}
			row.FromRow(value.(map[string]interface{}))
			s = append(s, row)
		}
		return s
	}
	s = append(s, values...)

	return s
}

func NewNgingCloudBackupLog(ctx echo.Context) *NgingCloudBackupLog {
	m := &NgingCloudBackupLog{}
	m.SetContext(ctx)
	return m
}

// NgingCloudBackupLog 云备份日志
type NgingCloudBackupLog struct {
	base    factory.Base
	objects []*NgingCloudBackupLog

	Id         uint64 `db:"id,omitempty,pk" bson:"id,omitempty" comment:"ID" json:"id" xml:"id"`
	BackupId   uint   `db:"backup_id" bson:"backup_id" comment:"云备份规则ID" json:"backup_id" xml:"backup_id"`
	BackupType string `db:"backup_type" bson:"backup_type" comment:"备份方式(full-全量备份;change-文件更改时触发的备份)" json:"backup_type" xml:"backup_type"`
	BackupFile string `db:"backup_file" bson:"backup_file" comment:"需要备份本地文件" json:"backup_file" xml:"backup_file"`
	RemoteFile string `db:"remote_file" bson:"remote_file" comment:"保存到远程的文件路径" json:"remote_file" xml:"remote_file"`
	Operation  string `db:"operation" bson:"operation" comment:"操作" json:"operation" xml:"operation"`
	Error      string `db:"error" bson:"error" comment:"错误信息" json:"error" xml:"error"`
	Status     string `db:"status" bson:"status" comment:"状态" json:"status" xml:"status"`
	Size       uint64 `db:"size" bson:"size" comment:"文件大小(字节)" json:"size" xml:"size"`
	Elapsed    uint   `db:"elapsed" bson:"elapsed" comment:"消耗时间(毫秒)" json:"elapsed" xml:"elapsed"`
	Created    uint   `db:"created" bson:"created" comment:"创建时间" json:"created" xml:"created"`
}

// - base function

func (a *NgingCloudBackupLog) Trans() factory.Transactioner {
	return a.base.Trans()
}

func (a *NgingCloudBackupLog) Use(trans factory.Transactioner) factory.Model {
	a.base.Use(trans)
	return a
}

func (a *NgingCloudBackupLog) SetContext(ctx echo.Context) factory.Model {
	a.base.SetContext(ctx)
	return a
}

func (a *NgingCloudBackupLog) EventON(on ...bool) factory.Model {
	a.base.EventON(on...)
	return a
}

func (a *NgingCloudBackupLog) EventOFF(off ...bool) factory.Model {
	a.base.EventOFF(off...)
	return a
}

func (a *NgingCloudBackupLog) Context() echo.Context {
	return a.base.Context()
}

func (a *NgingCloudBackupLog) SetConnID(connID int) factory.Model {
	a.base.SetConnID(connID)
	return a
}

func (a *NgingCloudBackupLog) ConnID() int {
	return a.base.ConnID()
}

func (a *NgingCloudBackupLog) SetNamer(namer func(factory.Model) string) factory.Model {
	a.base.SetNamer(namer)
	return a
}

func (a *NgingCloudBackupLog) Namer() func(factory.Model) string {
	return a.base.Namer()
}

func (a *NgingCloudBackupLog) SetParam(param *factory.Param) factory.Model {
	a.base.SetParam(param)
	return a
}

func (a *NgingCloudBackupLog) Param(mw func(db.Result) db.Result, args ...interface{}) *factory.Param {
	if a.base.Param() == nil {
		return a.NewParam().SetMiddleware(mw).SetArgs(args...)
	}
	return a.base.Param().SetMiddleware(mw).SetArgs(args...)
}

func (a *NgingCloudBackupLog) New(structName string, connID ...int) factory.Model {
	return a.base.New(structName, connID...)
}

func (a *NgingCloudBackupLog) Base_() factory.Baser {
	return &a.base
}

// - current function

func (a *NgingCloudBackupLog) Objects() []*NgingCloudBackupLog {
	if a.objects == nil {
		return nil
	}
	return a.objects[:]
}

func (a *NgingCloudBackupLog) XObjects() Slice_NgingCloudBackupLog {
	return Slice_NgingCloudBackupLog(a.Objects())
}

func (a *NgingCloudBackupLog) NewObjects() factory.Ranger {
	return &Slice_NgingCloudBackupLog{}
}

func (a *NgingCloudBackupLog) InitObjects() *[]*NgingCloudBackupLog {
	a.objects = []*NgingCloudBackupLog{}
	return &a.objects
}

func (a *NgingCloudBackupLog) NewParam() *factory.Param {
	return factory.NewParam(factory.DefaultFactory).SetIndex(a.base.ConnID()).SetTrans(a.base.Trans()).SetCollection(a.Name_()).SetModel(a)
}

func (a *NgingCloudBackupLog) Short_() string {
	return "nging_cloud_backup_log"
}

func (a *NgingCloudBackupLog) Struct_() string {
	return "NgingCloudBackupLog"
}

func (a *NgingCloudBackupLog) Name_() string {
	b := a
	if b == nil {
		b = &NgingCloudBackupLog{}
	}
	if b.base.Namer() != nil {
		return WithPrefix(b.base.Namer()(b))
	}
	return WithPrefix(factory.TableNamerGet(b.Short_())(b))
}

func (a *NgingCloudBackupLog) CPAFrom(source factory.Model) factory.Model {
	a.SetContext(source.Context())
	a.SetConnID(source.ConnID())
	a.SetNamer(source.Namer())
	return a
}

func (a *NgingCloudBackupLog) Get(mw func(db.Result) db.Result, args ...interface{}) (err error) {
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

func (a *NgingCloudBackupLog) List(recv interface{}, mw func(db.Result) db.Result, page, size int, args ...interface{}) (func() int64, error) {
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
		case *[]*NgingCloudBackupLog:
			err = DBI.FireReaded(a, queryParam, Slice_NgingCloudBackupLog(*v))
		case []*NgingCloudBackupLog:
			err = DBI.FireReaded(a, queryParam, Slice_NgingCloudBackupLog(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *NgingCloudBackupLog) GroupBy(keyField string, inputRows ...[]*NgingCloudBackupLog) map[string][]*NgingCloudBackupLog {
	var rows Slice_NgingCloudBackupLog
	if len(inputRows) > 0 {
		rows = Slice_NgingCloudBackupLog(inputRows[0])
	} else {
		rows = Slice_NgingCloudBackupLog(a.Objects())
	}
	return rows.GroupBy(keyField)
}

func (a *NgingCloudBackupLog) KeyBy(keyField string, inputRows ...[]*NgingCloudBackupLog) map[string]*NgingCloudBackupLog {
	var rows Slice_NgingCloudBackupLog
	if len(inputRows) > 0 {
		rows = Slice_NgingCloudBackupLog(inputRows[0])
	} else {
		rows = Slice_NgingCloudBackupLog(a.Objects())
	}
	return rows.KeyBy(keyField)
}

func (a *NgingCloudBackupLog) AsKV(keyField string, valueField string, inputRows ...[]*NgingCloudBackupLog) param.Store {
	var rows Slice_NgingCloudBackupLog
	if len(inputRows) > 0 {
		rows = Slice_NgingCloudBackupLog(inputRows[0])
	} else {
		rows = Slice_NgingCloudBackupLog(a.Objects())
	}
	return rows.AsKV(keyField, valueField)
}

func (a *NgingCloudBackupLog) ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error) {
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
		case *[]*NgingCloudBackupLog:
			err = DBI.FireReaded(a, queryParam, Slice_NgingCloudBackupLog(*v))
		case []*NgingCloudBackupLog:
			err = DBI.FireReaded(a, queryParam, Slice_NgingCloudBackupLog(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *NgingCloudBackupLog) Insert() (pk interface{}, err error) {
	a.Created = uint(time.Now().Unix())
	a.Id = 0
	if len(a.BackupType) == 0 {
		a.BackupType = "change"
	}
	if len(a.Operation) == 0 {
		a.Operation = "none"
	}
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

func (a *NgingCloudBackupLog) Update(mw func(db.Result) db.Result, args ...interface{}) (err error) {

	if len(a.BackupType) == 0 {
		a.BackupType = "change"
	}
	if len(a.Operation) == 0 {
		a.Operation = "none"
	}
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

func (a *NgingCloudBackupLog) Updatex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {

	if len(a.BackupType) == 0 {
		a.BackupType = "change"
	}
	if len(a.Operation) == 0 {
		a.Operation = "none"
	}
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

func (a *NgingCloudBackupLog) UpdateByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (err error) {

	if len(a.BackupType) == 0 {
		a.BackupType = "change"
	}
	if len(a.Operation) == 0 {
		a.Operation = "none"
	}
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

func (a *NgingCloudBackupLog) UpdatexByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (affected int64, err error) {

	if len(a.BackupType) == 0 {
		a.BackupType = "change"
	}
	if len(a.Operation) == 0 {
		a.Operation = "none"
	}
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

func (a *NgingCloudBackupLog) UpdateField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (err error) {
	return a.UpdateFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *NgingCloudBackupLog) UpdatexField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (affected int64, err error) {
	return a.UpdatexFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *NgingCloudBackupLog) UpdateFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (err error) {

	if val, ok := kvset["backup_type"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["backup_type"] = "change"
		}
	}
	if val, ok := kvset["operation"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["operation"] = "none"
		}
	}
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

func (a *NgingCloudBackupLog) UpdatexFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (affected int64, err error) {

	if val, ok := kvset["backup_type"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["backup_type"] = "change"
		}
	}
	if val, ok := kvset["operation"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["operation"] = "none"
		}
	}
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

func (a *NgingCloudBackupLog) UpdateValues(mw func(db.Result) db.Result, keysValues *db.KeysValues, args ...interface{}) (err error) {
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

func (a *NgingCloudBackupLog) Upsert(mw func(db.Result) db.Result, args ...interface{}) (pk interface{}, err error) {
	pk, err = a.Param(mw, args...).SetSend(a).Upsert(func() error {
		if len(a.BackupType) == 0 {
			a.BackupType = "change"
		}
		if len(a.Operation) == 0 {
			a.Operation = "none"
		}
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
		if len(a.BackupType) == 0 {
			a.BackupType = "change"
		}
		if len(a.Operation) == 0 {
			a.Operation = "none"
		}
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

func (a *NgingCloudBackupLog) Delete(mw func(db.Result) db.Result, args ...interface{}) (err error) {

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

func (a *NgingCloudBackupLog) Deletex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {

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

func (a *NgingCloudBackupLog) Count(mw func(db.Result) db.Result, args ...interface{}) (int64, error) {
	return a.Param(mw, args...).Count()
}

func (a *NgingCloudBackupLog) Exists(mw func(db.Result) db.Result, args ...interface{}) (bool, error) {
	return a.Param(mw, args...).Exists()
}

func (a *NgingCloudBackupLog) Reset() *NgingCloudBackupLog {
	a.Id = 0
	a.BackupId = 0
	a.BackupType = ``
	a.BackupFile = ``
	a.RemoteFile = ``
	a.Operation = ``
	a.Error = ``
	a.Status = ``
	a.Size = 0
	a.Elapsed = 0
	a.Created = 0
	return a
}

func (a *NgingCloudBackupLog) AsMap(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["Id"] = a.Id
		r["BackupId"] = a.BackupId
		r["BackupType"] = a.BackupType
		r["BackupFile"] = a.BackupFile
		r["RemoteFile"] = a.RemoteFile
		r["Operation"] = a.Operation
		r["Error"] = a.Error
		r["Status"] = a.Status
		r["Size"] = a.Size
		r["Elapsed"] = a.Elapsed
		r["Created"] = a.Created
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "Id":
			r["Id"] = a.Id
		case "BackupId":
			r["BackupId"] = a.BackupId
		case "BackupType":
			r["BackupType"] = a.BackupType
		case "BackupFile":
			r["BackupFile"] = a.BackupFile
		case "RemoteFile":
			r["RemoteFile"] = a.RemoteFile
		case "Operation":
			r["Operation"] = a.Operation
		case "Error":
			r["Error"] = a.Error
		case "Status":
			r["Status"] = a.Status
		case "Size":
			r["Size"] = a.Size
		case "Elapsed":
			r["Elapsed"] = a.Elapsed
		case "Created":
			r["Created"] = a.Created
		}
	}
	return r
}

func (a *NgingCloudBackupLog) FromRow(row map[string]interface{}) {
	for key, value := range row {
		if _, ok := value.(db.RawValue); ok {
			continue
		}
		switch key {
		case "id":
			a.Id = param.AsUint64(value)
		case "backup_id":
			a.BackupId = param.AsUint(value)
		case "backup_type":
			a.BackupType = param.AsString(value)
		case "backup_file":
			a.BackupFile = param.AsString(value)
		case "remote_file":
			a.RemoteFile = param.AsString(value)
		case "operation":
			a.Operation = param.AsString(value)
		case "error":
			a.Error = param.AsString(value)
		case "status":
			a.Status = param.AsString(value)
		case "size":
			a.Size = param.AsUint64(value)
		case "elapsed":
			a.Elapsed = param.AsUint(value)
		case "created":
			a.Created = param.AsUint(value)
		}
	}
}

func (a *NgingCloudBackupLog) GetField(field string) interface{} {
	switch field {
	case "Id":
		return a.Id
	case "BackupId":
		return a.BackupId
	case "BackupType":
		return a.BackupType
	case "BackupFile":
		return a.BackupFile
	case "RemoteFile":
		return a.RemoteFile
	case "Operation":
		return a.Operation
	case "Error":
		return a.Error
	case "Status":
		return a.Status
	case "Size":
		return a.Size
	case "Elapsed":
		return a.Elapsed
	case "Created":
		return a.Created
	default:
		return nil
	}
}

func (a *NgingCloudBackupLog) GetAllFieldNames() []string {
	return []string{
		"Id",
		"BackupId",
		"BackupType",
		"BackupFile",
		"RemoteFile",
		"Operation",
		"Error",
		"Status",
		"Size",
		"Elapsed",
		"Created",
	}
}

func (a *NgingCloudBackupLog) HasField(field string) bool {
	switch field {
	case "Id":
		return true
	case "BackupId":
		return true
	case "BackupType":
		return true
	case "BackupFile":
		return true
	case "RemoteFile":
		return true
	case "Operation":
		return true
	case "Error":
		return true
	case "Status":
		return true
	case "Size":
		return true
	case "Elapsed":
		return true
	case "Created":
		return true
	default:
		return false
	}
}

func (a *NgingCloudBackupLog) Set(key interface{}, value ...interface{}) {
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
		case "BackupId":
			a.BackupId = param.AsUint(vv)
		case "BackupType":
			a.BackupType = param.AsString(vv)
		case "BackupFile":
			a.BackupFile = param.AsString(vv)
		case "RemoteFile":
			a.RemoteFile = param.AsString(vv)
		case "Operation":
			a.Operation = param.AsString(vv)
		case "Error":
			a.Error = param.AsString(vv)
		case "Status":
			a.Status = param.AsString(vv)
		case "Size":
			a.Size = param.AsUint64(vv)
		case "Elapsed":
			a.Elapsed = param.AsUint(vv)
		case "Created":
			a.Created = param.AsUint(vv)
		}
	}
}

func (a *NgingCloudBackupLog) AsRow(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["id"] = a.Id
		r["backup_id"] = a.BackupId
		r["backup_type"] = a.BackupType
		r["backup_file"] = a.BackupFile
		r["remote_file"] = a.RemoteFile
		r["operation"] = a.Operation
		r["error"] = a.Error
		r["status"] = a.Status
		r["size"] = a.Size
		r["elapsed"] = a.Elapsed
		r["created"] = a.Created
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "id":
			r["id"] = a.Id
		case "backup_id":
			r["backup_id"] = a.BackupId
		case "backup_type":
			r["backup_type"] = a.BackupType
		case "backup_file":
			r["backup_file"] = a.BackupFile
		case "remote_file":
			r["remote_file"] = a.RemoteFile
		case "operation":
			r["operation"] = a.Operation
		case "error":
			r["error"] = a.Error
		case "status":
			r["status"] = a.Status
		case "size":
			r["size"] = a.Size
		case "elapsed":
			r["elapsed"] = a.Elapsed
		case "created":
			r["created"] = a.Created
		}
	}
	return r
}

func (a *NgingCloudBackupLog) ListPage(cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPage(a, cond, sorts...)
}

func (a *NgingCloudBackupLog) ListPageAs(recv interface{}, cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPageAs(a, recv, cond, sorts...)
}

func (a *NgingCloudBackupLog) ListPageByOffset(cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPageByOffset(a, cond, sorts...)
}

func (a *NgingCloudBackupLog) ListPageByOffsetAs(recv interface{}, cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPageByOffsetAs(a, recv, cond, sorts...)
}

func (a *NgingCloudBackupLog) BatchValidate(kvset map[string]interface{}) error {
	if kvset == nil {
		kvset = a.AsRow()
	}
	return DBI.Fields.BatchValidate(a.Short_(), kvset)
}

func (a *NgingCloudBackupLog) Validate(field string, value interface{}) error {
	return DBI.Fields.Validate(a.Short_(), field, value)
}
