// Code generated by 'yaegi extract github.com/webx-top/db'. DO NOT EDIT.

package lib

import (
	"context"
	"github.com/webx-top/db"
	"go/constant"
	"go/token"
	"reflect"
	"time"
)

func init() {
	Symbols["github.com/webx-top/db/db"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"After":                                  reflect.ValueOf(db.After),
		"And":                                    reflect.ValueOf(db.And),
		"Before":                                 reflect.ValueOf(db.Before),
		"Between":                                reflect.ValueOf(db.Between),
		"BuildSQL":                               reflect.ValueOf(db.BuildSQL),
		"ComparisonOperatorAfter":                reflect.ValueOf(db.ComparisonOperatorAfter),
		"ComparisonOperatorBefore":               reflect.ValueOf(db.ComparisonOperatorBefore),
		"ComparisonOperatorBetween":              reflect.ValueOf(db.ComparisonOperatorBetween),
		"ComparisonOperatorEqual":                reflect.ValueOf(db.ComparisonOperatorEqual),
		"ComparisonOperatorGreaterThan":          reflect.ValueOf(db.ComparisonOperatorGreaterThan),
		"ComparisonOperatorGreaterThanOrEqualTo": reflect.ValueOf(db.ComparisonOperatorGreaterThanOrEqualTo),
		"ComparisonOperatorIn":                   reflect.ValueOf(db.ComparisonOperatorIn),
		"ComparisonOperatorIs":                   reflect.ValueOf(db.ComparisonOperatorIs),
		"ComparisonOperatorIsNot":                reflect.ValueOf(db.ComparisonOperatorIsNot),
		"ComparisonOperatorLessThan":             reflect.ValueOf(db.ComparisonOperatorLessThan),
		"ComparisonOperatorLessThanOrEqualTo":    reflect.ValueOf(db.ComparisonOperatorLessThanOrEqualTo),
		"ComparisonOperatorLike":                 reflect.ValueOf(db.ComparisonOperatorLike),
		"ComparisonOperatorNone":                 reflect.ValueOf(db.ComparisonOperatorNone),
		"ComparisonOperatorNotBetween":           reflect.ValueOf(db.ComparisonOperatorNotBetween),
		"ComparisonOperatorNotEqual":             reflect.ValueOf(db.ComparisonOperatorNotEqual),
		"ComparisonOperatorNotIn":                reflect.ValueOf(db.ComparisonOperatorNotIn),
		"ComparisonOperatorNotLike":              reflect.ValueOf(db.ComparisonOperatorNotLike),
		"ComparisonOperatorNotRegExp":            reflect.ValueOf(db.ComparisonOperatorNotRegExp),
		"ComparisonOperatorOnOrAfter":            reflect.ValueOf(db.ComparisonOperatorOnOrAfter),
		"ComparisonOperatorOnOrBefore":           reflect.ValueOf(db.ComparisonOperatorOnOrBefore),
		"ComparisonOperatorRegExp":               reflect.ValueOf(db.ComparisonOperatorRegExp),
		"CompoundsPoolGet":                       reflect.ValueOf(db.CompoundsPoolGet),
		"CompoundsPoolRelease":                   reflect.ValueOf(db.CompoundsPoolRelease),
		"DefaultSettings":                        reflect.ValueOf(&db.DefaultSettings).Elem(),
		"EmptyCond":                              reflect.ValueOf(&db.EmptyCond).Elem(),
		"EnvEnableDebug":                         reflect.ValueOf(constant.MakeFromLiteral("\"UPPERIO_DB_DEBUG\"", token.STRING, 0)),
		"Eq":                                     reflect.ValueOf(db.Eq),
		"ErrAlreadyWithinTransaction":            reflect.ValueOf(&db.ErrAlreadyWithinTransaction).Elem(),
		"ErrCollectionDoesNotExist":              reflect.ValueOf(&db.ErrCollectionDoesNotExist).Elem(),
		"ErrGivingUpTryingToConnect":             reflect.ValueOf(&db.ErrGivingUpTryingToConnect).Elem(),
		"ErrMissingCollectionName":               reflect.ValueOf(&db.ErrMissingCollectionName).Elem(),
		"ErrMissingConditions":                   reflect.ValueOf(&db.ErrMissingConditions).Elem(),
		"ErrMissingConnURL":                      reflect.ValueOf(&db.ErrMissingConnURL).Elem(),
		"ErrMissingDatabaseName":                 reflect.ValueOf(&db.ErrMissingDatabaseName).Elem(),
		"ErrNoMoreRows":                          reflect.ValueOf(&db.ErrNoMoreRows).Elem(),
		"ErrNotConnected":                        reflect.ValueOf(&db.ErrNotConnected).Elem(),
		"ErrNotImplemented":                      reflect.ValueOf(&db.ErrNotImplemented).Elem(),
		"ErrQueryIsPending":                      reflect.ValueOf(&db.ErrQueryIsPending).Elem(),
		"ErrQueryLimitParam":                     reflect.ValueOf(&db.ErrQueryLimitParam).Elem(),
		"ErrQueryOffsetParam":                    reflect.ValueOf(&db.ErrQueryOffsetParam).Elem(),
		"ErrQuerySortParam":                      reflect.ValueOf(&db.ErrQuerySortParam).Elem(),
		"ErrSockerOrHost":                        reflect.ValueOf(&db.ErrSockerOrHost).Elem(),
		"ErrTooManyClients":                      reflect.ValueOf(&db.ErrTooManyClients).Elem(),
		"ErrUndefined":                           reflect.ValueOf(&db.ErrUndefined).Elem(),
		"ErrUnknownConditionType":                reflect.ValueOf(&db.ErrUnknownConditionType).Elem(),
		"ErrUnsupported":                         reflect.ValueOf(&db.ErrUnsupported).Elem(),
		"ErrUnsupportedDestination":              reflect.ValueOf(&db.ErrUnsupportedDestination).Elem(),
		"ErrUnsupportedType":                     reflect.ValueOf(&db.ErrUnsupportedType).Elem(),
		"ErrUnsupportedValue":                    reflect.ValueOf(&db.ErrUnsupportedValue).Elem(),
		"Func":                                   reflect.ValueOf(db.Func),
		"Gt":                                     reflect.ValueOf(db.Gt),
		"Gte":                                    reflect.ValueOf(db.Gte),
		"In":                                     reflect.ValueOf(db.In),
		"Is":                                     reflect.ValueOf(db.Is),
		"IsNot":                                  reflect.ValueOf(db.IsNot),
		"IsNotNull":                              reflect.ValueOf(db.IsNotNull),
		"IsNull":                                 reflect.ValueOf(db.IsNull),
		"Like":                                   reflect.ValueOf(db.Like),
		"Lt":                                     reflect.ValueOf(db.Lt),
		"Lte":                                    reflect.ValueOf(db.Lte),
		"NewCompounds":                           reflect.ValueOf(db.NewCompounds),
		"NewConstraint":                          reflect.ValueOf(db.NewConstraint),
		"NewKeysValues":                          reflect.ValueOf(db.NewKeysValues),
		"NewSettings":                            reflect.ValueOf(db.NewSettings),
		"NotBetween":                             reflect.ValueOf(db.NotBetween),
		"NotEq":                                  reflect.ValueOf(db.NotEq),
		"NotIn":                                  reflect.ValueOf(db.NotIn),
		"NotLike":                                reflect.ValueOf(db.NotLike),
		"NotRegExp":                              reflect.ValueOf(db.NotRegExp),
		"OnOrAfter":                              reflect.ValueOf(db.OnOrAfter),
		"OnOrBefore":                             reflect.ValueOf(db.OnOrBefore),
		"Op":                                     reflect.ValueOf(db.Op),
		"Open":                                   reflect.ValueOf(db.Open),
		"OperatorAnd":                            reflect.ValueOf(db.OperatorAnd),
		"OperatorNone":                           reflect.ValueOf(db.OperatorNone),
		"OperatorOr":                             reflect.ValueOf(db.OperatorOr),
		"Or":                                     reflect.ValueOf(db.Or),
		"Raw":                                    reflect.ValueOf(db.Raw),
		"RegExp":                                 reflect.ValueOf(db.RegExp),
		"RegisterAdapter":                        reflect.ValueOf(db.RegisterAdapter),
		"Table":                                  reflect.ValueOf(db.Table),

		// type definitions
		"AdapterFuncMap":     reflect.ValueOf((*db.AdapterFuncMap)(nil)),
		"Collection":         reflect.ValueOf((*db.Collection)(nil)),
		"Comparison":         reflect.ValueOf((*db.Comparison)(nil)),
		"ComparisonOperator": reflect.ValueOf((*db.ComparisonOperator)(nil)),
		"Compound":           reflect.ValueOf((*db.Compound)(nil)),
		"CompoundOperator":   reflect.ValueOf((*db.CompoundOperator)(nil)),
		"Compounds":          reflect.ValueOf((*db.Compounds)(nil)),
		"Cond":               reflect.ValueOf((*db.Cond)(nil)),
		"ConnectionURL":      reflect.ValueOf((*db.ConnectionURL)(nil)),
		"Constraint":         reflect.ValueOf((*db.Constraint)(nil)),
		"Constraints":        reflect.ValueOf((*db.Constraints)(nil)),
		"Database":           reflect.ValueOf((*db.Database)(nil)),
		"Function":           reflect.ValueOf((*db.Function)(nil)),
		"Intersection":       reflect.ValueOf((*db.Intersection)(nil)),
		"KeysValues":         reflect.ValueOf((*db.KeysValues)(nil)),
		"Logger":             reflect.ValueOf((*db.Logger)(nil)),
		"Marshaler":          reflect.ValueOf((*db.Marshaler)(nil)),
		"Method":             reflect.ValueOf((*db.Method)(nil)),
		"QueryStatus":        reflect.ValueOf((*db.QueryStatus)(nil)),
		"RawValue":           reflect.ValueOf((*db.RawValue)(nil)),
		"RequestURI":         reflect.ValueOf((*db.RequestURI)(nil)),
		"Result":             reflect.ValueOf((*db.Result)(nil)),
		"Settings":           reflect.ValueOf((*db.Settings)(nil)),
		"StdContext":         reflect.ValueOf((*db.StdContext)(nil)),
		"TableName":          reflect.ValueOf((*db.TableName)(nil)),
		"Tx":                 reflect.ValueOf((*db.Tx)(nil)),
		"Union":              reflect.ValueOf((*db.Union)(nil)),
		"Unmarshaler":        reflect.ValueOf((*db.Unmarshaler)(nil)),

		// interface wrapper definitions
		"_Collection":    reflect.ValueOf((*_github_com_webx_top_db_Collection)(nil)),
		"_Comparison":    reflect.ValueOf((*_github_com_webx_top_db_Comparison)(nil)),
		"_Compound":      reflect.ValueOf((*_github_com_webx_top_db_Compound)(nil)),
		"_ConnectionURL": reflect.ValueOf((*_github_com_webx_top_db_ConnectionURL)(nil)),
		"_Constraint":    reflect.ValueOf((*_github_com_webx_top_db_Constraint)(nil)),
		"_Constraints":   reflect.ValueOf((*_github_com_webx_top_db_Constraints)(nil)),
		"_Database":      reflect.ValueOf((*_github_com_webx_top_db_Database)(nil)),
		"_Function":      reflect.ValueOf((*_github_com_webx_top_db_Function)(nil)),
		"_Logger":        reflect.ValueOf((*_github_com_webx_top_db_Logger)(nil)),
		"_Marshaler":     reflect.ValueOf((*_github_com_webx_top_db_Marshaler)(nil)),
		"_Method":        reflect.ValueOf((*_github_com_webx_top_db_Method)(nil)),
		"_RawValue":      reflect.ValueOf((*_github_com_webx_top_db_RawValue)(nil)),
		"_RequestURI":    reflect.ValueOf((*_github_com_webx_top_db_RequestURI)(nil)),
		"_Result":        reflect.ValueOf((*_github_com_webx_top_db_Result)(nil)),
		"_Settings":      reflect.ValueOf((*_github_com_webx_top_db_Settings)(nil)),
		"_StdContext":    reflect.ValueOf((*_github_com_webx_top_db_StdContext)(nil)),
		"_TableName":     reflect.ValueOf((*_github_com_webx_top_db_TableName)(nil)),
		"_Tx":            reflect.ValueOf((*_github_com_webx_top_db_Tx)(nil)),
		"_Unmarshaler":   reflect.ValueOf((*_github_com_webx_top_db_Unmarshaler)(nil)),
	}
}

// _github_com_webx_top_db_Collection is an interface wrapper for Collection type
type _github_com_webx_top_db_Collection struct {
	IValue           interface{}
	WExists          func() bool
	WFind            func(a0 ...interface{}) db.Result
	WInsert          func(a0 interface{}) (interface{}, error)
	WInsertReturning func(a0 interface{}) error
	WName            func() string
	WTruncate        func() error
	WUpdateReturning func(a0 interface{}) error
}

func (W _github_com_webx_top_db_Collection) Exists() bool {
	return W.WExists()
}
func (W _github_com_webx_top_db_Collection) Find(a0 ...interface{}) db.Result {
	return W.WFind(a0...)
}
func (W _github_com_webx_top_db_Collection) Insert(a0 interface{}) (interface{}, error) {
	return W.WInsert(a0)
}
func (W _github_com_webx_top_db_Collection) InsertReturning(a0 interface{}) error {
	return W.WInsertReturning(a0)
}
func (W _github_com_webx_top_db_Collection) Name() string {
	return W.WName()
}
func (W _github_com_webx_top_db_Collection) Truncate() error {
	return W.WTruncate()
}
func (W _github_com_webx_top_db_Collection) UpdateReturning(a0 interface{}) error {
	return W.WUpdateReturning(a0)
}

// _github_com_webx_top_db_Comparison is an interface wrapper for Comparison type
type _github_com_webx_top_db_Comparison struct {
	IValue    interface{}
	WOperator func() db.ComparisonOperator
	WValue    func() interface{}
}

func (W _github_com_webx_top_db_Comparison) Operator() db.ComparisonOperator {
	return W.WOperator()
}
func (W _github_com_webx_top_db_Comparison) Value() interface{} {
	return W.WValue()
}

// _github_com_webx_top_db_Compound is an interface wrapper for Compound type
type _github_com_webx_top_db_Compound struct {
	IValue     interface{}
	WEmpty     func() bool
	WOperator  func() db.CompoundOperator
	WSentences func() []db.Compound
}

func (W _github_com_webx_top_db_Compound) Empty() bool {
	return W.WEmpty()
}
func (W _github_com_webx_top_db_Compound) Operator() db.CompoundOperator {
	return W.WOperator()
}
func (W _github_com_webx_top_db_Compound) Sentences() []db.Compound {
	return W.WSentences()
}

// _github_com_webx_top_db_ConnectionURL is an interface wrapper for ConnectionURL type
type _github_com_webx_top_db_ConnectionURL struct {
	IValue  interface{}
	WString func() string
}

func (W _github_com_webx_top_db_ConnectionURL) String() string {
	if W.WString == nil {
		return ""
	}
	return W.WString()
}

// _github_com_webx_top_db_Constraint is an interface wrapper for Constraint type
type _github_com_webx_top_db_Constraint struct {
	IValue interface{}
	WKey   func() interface{}
	WValue func() interface{}
}

func (W _github_com_webx_top_db_Constraint) Key() interface{} {
	return W.WKey()
}
func (W _github_com_webx_top_db_Constraint) Value() interface{} {
	return W.WValue()
}

// _github_com_webx_top_db_Constraints is an interface wrapper for Constraints type
type _github_com_webx_top_db_Constraints struct {
	IValue       interface{}
	WConstraints func() []db.Constraint
	WKeys        func() []interface{}
}

func (W _github_com_webx_top_db_Constraints) Constraints() []db.Constraint {
	return W.WConstraints()
}
func (W _github_com_webx_top_db_Constraints) Keys() []interface{} {
	return W.WKeys()
}

// _github_com_webx_top_db_Database is an interface wrapper for Database type
type _github_com_webx_top_db_Database struct {
	IValue                         interface{}
	WClearCache                    func()
	WClose                         func() error
	WCollection                    func(a0 string) db.Collection
	WCollections                   func() ([]string, error)
	WConnMaxLifetime               func() time.Duration
	WConnectionURL                 func() db.ConnectionURL
	WDriver                        func() interface{}
	WLog                           func(a0 *db.QueryStatus)
	WLogger                        func() db.Logger
	WLoggingElapsed                func() time.Duration
	WLoggingElapsedMs              func() uint32
	WLoggingEnabled                func() bool
	WMaxIdleConns                  func() int
	WMaxOpenConns                  func() int
	WName                          func() string
	WOpen                          func(a0 db.ConnectionURL) error
	WPing                          func() error
	WPreparedStatementCacheEnabled func() bool
	WSetConnMaxLifetime            func(a0 time.Duration)
	WSetLogger                     func(a0 db.Logger)
	WSetLogging                    func(a0 bool, a1 ...uint32)
	WSetLoggingElapsedMs           func(elapsedMs uint32)
	WSetMaxIdleConns               func(a0 int)
	WSetMaxOpenConns               func(a0 int)
	WSetPreparedStatementCache     func(a0 bool)
}

func (W _github_com_webx_top_db_Database) ClearCache() {
	W.WClearCache()
}
func (W _github_com_webx_top_db_Database) Close() error {
	return W.WClose()
}
func (W _github_com_webx_top_db_Database) Collection(a0 string) db.Collection {
	return W.WCollection(a0)
}
func (W _github_com_webx_top_db_Database) Collections() ([]string, error) {
	return W.WCollections()
}
func (W _github_com_webx_top_db_Database) ConnMaxLifetime() time.Duration {
	return W.WConnMaxLifetime()
}
func (W _github_com_webx_top_db_Database) ConnectionURL() db.ConnectionURL {
	return W.WConnectionURL()
}
func (W _github_com_webx_top_db_Database) Driver() interface{} {
	return W.WDriver()
}
func (W _github_com_webx_top_db_Database) Log(a0 *db.QueryStatus) {
	W.WLog(a0)
}
func (W _github_com_webx_top_db_Database) Logger() db.Logger {
	return W.WLogger()
}
func (W _github_com_webx_top_db_Database) LoggingElapsed() time.Duration {
	return W.WLoggingElapsed()
}
func (W _github_com_webx_top_db_Database) LoggingElapsedMs() uint32 {
	return W.WLoggingElapsedMs()
}
func (W _github_com_webx_top_db_Database) LoggingEnabled() bool {
	return W.WLoggingEnabled()
}
func (W _github_com_webx_top_db_Database) MaxIdleConns() int {
	return W.WMaxIdleConns()
}
func (W _github_com_webx_top_db_Database) MaxOpenConns() int {
	return W.WMaxOpenConns()
}
func (W _github_com_webx_top_db_Database) Name() string {
	return W.WName()
}
func (W _github_com_webx_top_db_Database) Open(a0 db.ConnectionURL) error {
	return W.WOpen(a0)
}
func (W _github_com_webx_top_db_Database) Ping() error {
	return W.WPing()
}
func (W _github_com_webx_top_db_Database) PreparedStatementCacheEnabled() bool {
	return W.WPreparedStatementCacheEnabled()
}
func (W _github_com_webx_top_db_Database) SetConnMaxLifetime(a0 time.Duration) {
	W.WSetConnMaxLifetime(a0)
}
func (W _github_com_webx_top_db_Database) SetLogger(a0 db.Logger) {
	W.WSetLogger(a0)
}
func (W _github_com_webx_top_db_Database) SetLogging(a0 bool, a1 ...uint32) {
	W.WSetLogging(a0, a1...)
}
func (W _github_com_webx_top_db_Database) SetLoggingElapsedMs(elapsedMs uint32) {
	W.WSetLoggingElapsedMs(elapsedMs)
}
func (W _github_com_webx_top_db_Database) SetMaxIdleConns(a0 int) {
	W.WSetMaxIdleConns(a0)
}
func (W _github_com_webx_top_db_Database) SetMaxOpenConns(a0 int) {
	W.WSetMaxOpenConns(a0)
}
func (W _github_com_webx_top_db_Database) SetPreparedStatementCache(a0 bool) {
	W.WSetPreparedStatementCache(a0)
}

// _github_com_webx_top_db_Function is an interface wrapper for Function type
type _github_com_webx_top_db_Function struct {
	IValue     interface{}
	WArguments func() []interface{}
	WName      func() string
}

func (W _github_com_webx_top_db_Function) Arguments() []interface{} {
	return W.WArguments()
}
func (W _github_com_webx_top_db_Function) Name() string {
	return W.WName()
}

// _github_com_webx_top_db_Logger is an interface wrapper for Logger type
type _github_com_webx_top_db_Logger struct {
	IValue interface{}
	WLog   func(a0 *db.QueryStatus)
}

func (W _github_com_webx_top_db_Logger) Log(a0 *db.QueryStatus) {
	W.WLog(a0)
}

// _github_com_webx_top_db_Marshaler is an interface wrapper for Marshaler type
type _github_com_webx_top_db_Marshaler struct {
	IValue     interface{}
	WMarshalDB func() (interface{}, error)
}

func (W _github_com_webx_top_db_Marshaler) MarshalDB() (interface{}, error) {
	return W.WMarshalDB()
}

// _github_com_webx_top_db_Method is an interface wrapper for Method type
type _github_com_webx_top_db_Method struct {
	IValue  interface{}
	WMethod func() string
}

func (W _github_com_webx_top_db_Method) Method() string {
	return W.WMethod()
}

// _github_com_webx_top_db_RawValue is an interface wrapper for RawValue type
type _github_com_webx_top_db_RawValue struct {
	IValue     interface{}
	WArguments func() []interface{}
	WEmpty     func() bool
	WOperator  func() db.CompoundOperator
	WRaw       func() string
	WSentences func() []db.Compound
	WString    func() string
}

func (W _github_com_webx_top_db_RawValue) Arguments() []interface{} {
	return W.WArguments()
}
func (W _github_com_webx_top_db_RawValue) Empty() bool {
	return W.WEmpty()
}
func (W _github_com_webx_top_db_RawValue) Operator() db.CompoundOperator {
	return W.WOperator()
}
func (W _github_com_webx_top_db_RawValue) Raw() string {
	return W.WRaw()
}
func (W _github_com_webx_top_db_RawValue) Sentences() []db.Compound {
	return W.WSentences()
}
func (W _github_com_webx_top_db_RawValue) String() string {
	if W.WString == nil {
		return ""
	}
	return W.WString()
}

// _github_com_webx_top_db_RequestURI is an interface wrapper for RequestURI type
type _github_com_webx_top_db_RequestURI struct {
	IValue      interface{}
	WRequestURI func() string
}

func (W _github_com_webx_top_db_RequestURI) RequestURI() string {
	return W.WRequestURI()
}

// _github_com_webx_top_db_Result is an interface wrapper for Result type
type _github_com_webx_top_db_Result struct {
	IValue        interface{}
	WAll          func(sliceOfStructs interface{}) error
	WAnd          func(a0 ...interface{}) db.Result
	WCallback     func(m interface{}) db.Result
	WClose        func() error
	WCount        func() (uint64, error)
	WCursor       func(cursorColumn string) db.Result
	WDelete       func() error
	WErr          func() error
	WExists       func() (bool, error)
	WForceIndex   func(index string) db.Result
	WGroup        func(a0 ...interface{}) db.Result
	WLimit        func(a0 int) db.Result
	WNext         func(ptrToStruct interface{}) bool
	WNextPage     func(cursorValue interface{}) db.Result
	WOffset       func(a0 int) db.Result
	WOne          func(ptrToStruct interface{}) error
	WOrderBy      func(a0 ...interface{}) db.Result
	WPage         func(pageNumber uint) db.Result
	WPaginate     func(pageSize uint) db.Result
	WPrevPage     func(cursorValue interface{}) db.Result
	WRelation     func(name string, fn interface{}) db.Result
	WSelect       func(a0 ...interface{}) db.Result
	WString       func() string
	WTotalEntries func() (uint64, error)
	WTotalPages   func() (uint, error)
	WUpdate       func(a0 interface{}) error
	WWhere        func(a0 ...interface{}) db.Result
}

func (W _github_com_webx_top_db_Result) All(sliceOfStructs interface{}) error {
	return W.WAll(sliceOfStructs)
}
func (W _github_com_webx_top_db_Result) And(a0 ...interface{}) db.Result {
	return W.WAnd(a0...)
}
func (W _github_com_webx_top_db_Result) Callback(m interface{}) db.Result {
	return W.WCallback(m)
}
func (W _github_com_webx_top_db_Result) Close() error {
	return W.WClose()
}
func (W _github_com_webx_top_db_Result) Count() (uint64, error) {
	return W.WCount()
}
func (W _github_com_webx_top_db_Result) Cursor(cursorColumn string) db.Result {
	return W.WCursor(cursorColumn)
}
func (W _github_com_webx_top_db_Result) Delete() error {
	return W.WDelete()
}
func (W _github_com_webx_top_db_Result) Err() error {
	return W.WErr()
}
func (W _github_com_webx_top_db_Result) Exists() (bool, error) {
	return W.WExists()
}
func (W _github_com_webx_top_db_Result) ForceIndex(index string) db.Result {
	return W.WForceIndex(index)
}
func (W _github_com_webx_top_db_Result) Group(a0 ...interface{}) db.Result {
	return W.WGroup(a0...)
}
func (W _github_com_webx_top_db_Result) Limit(a0 int) db.Result {
	return W.WLimit(a0)
}
func (W _github_com_webx_top_db_Result) Next(ptrToStruct interface{}) bool {
	return W.WNext(ptrToStruct)
}
func (W _github_com_webx_top_db_Result) NextPage(cursorValue interface{}) db.Result {
	return W.WNextPage(cursorValue)
}
func (W _github_com_webx_top_db_Result) Offset(a0 int) db.Result {
	return W.WOffset(a0)
}
func (W _github_com_webx_top_db_Result) One(ptrToStruct interface{}) error {
	return W.WOne(ptrToStruct)
}
func (W _github_com_webx_top_db_Result) OrderBy(a0 ...interface{}) db.Result {
	return W.WOrderBy(a0...)
}
func (W _github_com_webx_top_db_Result) Page(pageNumber uint) db.Result {
	return W.WPage(pageNumber)
}
func (W _github_com_webx_top_db_Result) Paginate(pageSize uint) db.Result {
	return W.WPaginate(pageSize)
}
func (W _github_com_webx_top_db_Result) PrevPage(cursorValue interface{}) db.Result {
	return W.WPrevPage(cursorValue)
}
func (W _github_com_webx_top_db_Result) Relation(name string, fn interface{}) db.Result {
	return W.WRelation(name, fn)
}
func (W _github_com_webx_top_db_Result) Select(a0 ...interface{}) db.Result {
	return W.WSelect(a0...)
}
func (W _github_com_webx_top_db_Result) String() string {
	if W.WString == nil {
		return ""
	}
	return W.WString()
}
func (W _github_com_webx_top_db_Result) TotalEntries() (uint64, error) {
	return W.WTotalEntries()
}
func (W _github_com_webx_top_db_Result) TotalPages() (uint, error) {
	return W.WTotalPages()
}
func (W _github_com_webx_top_db_Result) Update(a0 interface{}) error {
	return W.WUpdate(a0)
}
func (W _github_com_webx_top_db_Result) Where(a0 ...interface{}) db.Result {
	return W.WWhere(a0...)
}

// _github_com_webx_top_db_Settings is an interface wrapper for Settings type
type _github_com_webx_top_db_Settings struct {
	IValue                         interface{}
	WConnMaxLifetime               func() time.Duration
	WLog                           func(a0 *db.QueryStatus)
	WLogger                        func() db.Logger
	WLoggingElapsed                func() time.Duration
	WLoggingElapsedMs              func() uint32
	WLoggingEnabled                func() bool
	WMaxIdleConns                  func() int
	WMaxOpenConns                  func() int
	WPreparedStatementCacheEnabled func() bool
	WSetConnMaxLifetime            func(a0 time.Duration)
	WSetLogger                     func(a0 db.Logger)
	WSetLogging                    func(a0 bool, a1 ...uint32)
	WSetLoggingElapsedMs           func(elapsedMs uint32)
	WSetMaxIdleConns               func(a0 int)
	WSetMaxOpenConns               func(a0 int)
	WSetPreparedStatementCache     func(a0 bool)
}

func (W _github_com_webx_top_db_Settings) ConnMaxLifetime() time.Duration {
	return W.WConnMaxLifetime()
}
func (W _github_com_webx_top_db_Settings) Log(a0 *db.QueryStatus) {
	W.WLog(a0)
}
func (W _github_com_webx_top_db_Settings) Logger() db.Logger {
	return W.WLogger()
}
func (W _github_com_webx_top_db_Settings) LoggingElapsed() time.Duration {
	return W.WLoggingElapsed()
}
func (W _github_com_webx_top_db_Settings) LoggingElapsedMs() uint32 {
	return W.WLoggingElapsedMs()
}
func (W _github_com_webx_top_db_Settings) LoggingEnabled() bool {
	return W.WLoggingEnabled()
}
func (W _github_com_webx_top_db_Settings) MaxIdleConns() int {
	return W.WMaxIdleConns()
}
func (W _github_com_webx_top_db_Settings) MaxOpenConns() int {
	return W.WMaxOpenConns()
}
func (W _github_com_webx_top_db_Settings) PreparedStatementCacheEnabled() bool {
	return W.WPreparedStatementCacheEnabled()
}
func (W _github_com_webx_top_db_Settings) SetConnMaxLifetime(a0 time.Duration) {
	W.WSetConnMaxLifetime(a0)
}
func (W _github_com_webx_top_db_Settings) SetLogger(a0 db.Logger) {
	W.WSetLogger(a0)
}
func (W _github_com_webx_top_db_Settings) SetLogging(a0 bool, a1 ...uint32) {
	W.WSetLogging(a0, a1...)
}
func (W _github_com_webx_top_db_Settings) SetLoggingElapsedMs(elapsedMs uint32) {
	W.WSetLoggingElapsedMs(elapsedMs)
}
func (W _github_com_webx_top_db_Settings) SetMaxIdleConns(a0 int) {
	W.WSetMaxIdleConns(a0)
}
func (W _github_com_webx_top_db_Settings) SetMaxOpenConns(a0 int) {
	W.WSetMaxOpenConns(a0)
}
func (W _github_com_webx_top_db_Settings) SetPreparedStatementCache(a0 bool) {
	W.WSetPreparedStatementCache(a0)
}

// _github_com_webx_top_db_StdContext is an interface wrapper for StdContext type
type _github_com_webx_top_db_StdContext struct {
	IValue      interface{}
	WStdContext func() context.Context
}

func (W _github_com_webx_top_db_StdContext) StdContext() context.Context {
	return W.WStdContext()
}

// _github_com_webx_top_db_TableName is an interface wrapper for TableName type
type _github_com_webx_top_db_TableName struct {
	IValue     interface{}
	WTableName func() string
}

func (W _github_com_webx_top_db_TableName) TableName() string {
	return W.WTableName()
}

// _github_com_webx_top_db_Tx is an interface wrapper for Tx type
type _github_com_webx_top_db_Tx struct {
	IValue    interface{}
	WCommit   func() error
	WRollback func() error
}

func (W _github_com_webx_top_db_Tx) Commit() error {
	return W.WCommit()
}
func (W _github_com_webx_top_db_Tx) Rollback() error {
	return W.WRollback()
}

// _github_com_webx_top_db_Unmarshaler is an interface wrapper for Unmarshaler type
type _github_com_webx_top_db_Unmarshaler struct {
	IValue       interface{}
	WUnmarshalDB func(a0 interface{}) error
}

func (W _github_com_webx_top_db_Unmarshaler) UnmarshalDB(a0 interface{}) error {
	return W.WUnmarshalDB(a0)
}