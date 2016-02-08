package database

import (
	"database/sql"

	. "github.com/coscms/xorm"
	"github.com/coscms/xorm/core"
)

type Ormer interface {
	// Method Init reset the session as the init status.
	Init()

	// Method Close release the connection from pool
	Close()

	// Prepare
	Prepare() Ormer

	// Method Sql provides raw sql input parameter. When you have a complex SQL statement
	// and cannot use Where, Id, In and etc. Methods to describe, you can use Sql.
	Sql(querystring string, args ...interface{}) Ormer

	// Method Where provides custom query condition.
	Where(querystring string, args ...interface{}) Ormer

	// Method Where provides custom query condition.
	And(querystring string, args ...interface{}) Ormer

	// Method Where provides custom query condition.
	Or(querystring string, args ...interface{}) Ormer

	// Method Id provides converting id as a query condition
	Id(id interface{}) Ormer

	// Apply before Processor, affected bean is passed to closure arg
	Before(closures func(interface{})) Ormer

	// Apply after Processor, affected bean is passed to closure arg
	After(closures func(interface{})) Ormer

	// Method core.Table can input a string or pointer to struct for special a table to operate.
	Table(tableNameOrBean interface{}) Ormer

	// set the table alias
	Alias(alias string) Ormer

	// Method In provides a query string like "id in (1, 2, 3)"
	In(column string, args ...interface{}) Ormer

	// Method In provides a query string like "count = count + 1"
	Incr(column string, arg ...interface{}) Ormer

	// Method Decr provides a query string like "count = count - 1"
	Decr(column string, arg ...interface{}) Ormer

	// Method SetExpr provides a query string like "column = {expression}"
	SetExpr(column string, expression string) Ormer

	// Method Cols provides some columns to special
	Select(str string) Ormer

	// Method Cols provides some columns to special
	Cols(columns ...string) Ormer

	AllCols() Ormer

	MustCols(columns ...string) Ormer

	NoCascade() Ormer

	// Xorm automatically retrieve condition according struct, but
	// if struct has bool field, it will ignore them. So use UseBool
	// to tell system to do not ignore them.
	// If no paramters, it will use all the bool field of struct, or
	// it will use paramters's columns
	UseBool(columns ...string) Ormer

	// use for distinct columns. Caution: when you are using cache,
	// distinct will not be cached because cache system need id,
	// but distinct will not provide id
	Distinct(columns ...string) Ormer

	// Set Read/Write locking for UPDATE
	ForUpdate() Ormer

	// Only not use the paramters as select or update columns
	Omit(columns ...string) Ormer

	// Set null when column is zero-value and nullable for update
	Nullable(columns ...string) Ormer

	// Method NoAutoTime means do not automatically give created field and updated field
	// the current time on the current session temporarily
	NoAutoTime() Ormer

	NoAutoCondition(no ...bool) Ormer

	// Method Limit provide limit and offset query condition
	Limit(limit int, start ...int) Ormer

	// Method OrderBy provide order by query condition, the input parameter is the content
	// after order by on a sql statement.
	OrderBy(order string) Ormer

	// Method Desc provide desc order by query condition, the input parameters are columns.
	Desc(colNames ...string) Ormer

	// Method Asc provide asc order by query condition, the input parameters are columns.
	Asc(colNames ...string) Ormer

	// Method StoreEngine is only avialble mysql dialect currently
	StoreEngine(storeEngine string) Ormer

	// Method Charset is only avialble mysql dialect currently
	Charset(charset string) Ormer

	// Method Cascade indicates if loading sub Struct
	Cascade(trueOrFalse ...bool) Ormer

	// Method NoCache ask this session do not retrieve data from cache system and
	// get data from database directly.
	NoCache() Ormer

	//The join_operator should be one of INNER, LEFT OUTER, CROSS etc - this will be prepended to JOIN
	Join(join_operator string, tablename interface{}, condition string) Ormer

	// Generate Group By statement
	GroupBy(keys string) Ormer

	// Generate Having statement
	Having(conditions string) Ormer

	DB() *core.DB

	// Begin a transaction
	Begin() error

	// When using transaction, you can rollback if any error
	Rollback() error

	// When using transaction, Commit will commit all operations.
	Commit() error

	// Exec raw sql
	Exec(sqlStr string, args ...interface{}) (sql.Result, error)

	// this function create a table according a bean
	CreateTable(bean interface{}) error

	// create indexes
	CreateIndexes(bean interface{}) error

	// create uniques
	CreateUniques(bean interface{}) error

	// drop indexes
	DropIndexes(bean interface{}) error

	// drop table will drop table if exist, if drop failed, it will return error
	DropTable(beanOrTableName interface{}) error

	JoinColumns(cols []*core.Column, includeTableName bool) string

	// Return sql.Rows compatible Rows obj, as a forward Iterator object for iterating record by record, bean's non-empty fields
	// are conditions.
	Rows(bean interface{}) (*Rows, error)

	// Iterate record by record handle records from table, condiBeans's non-empty fields
	// are conditions. beans could be []Struct, []*Struct, map[int64]Struct
	// map[int64]*Struct
	Iterate(bean interface{}, fun IterFunc) error

	// get retrieve one record from database, bean's non-empty fields
	// will be as conditions
	Get(bean interface{}) (bool, error)

	// Count counts the records. bean's non-empty fields
	// are conditions.
	Count(bean interface{}) (int64, error)

	// Find retrieve records from table, condiBeans's non-empty fields
	// are conditions. beans could be []Struct, []*Struct, map[int64]Struct
	// map[int64]*Struct
	Find(rowsSlicePtr interface{}, condiBean ...interface{}) error

	// Test if database is ok
	Ping() error

	IsTableExist(beanOrTableName interface{}) (bool, error)

	IsTableEmpty(bean interface{}) (bool, error)

	// Exec a raw sql and return records as []map[string][]byte
	Query(sqlStr string, paramStr ...interface{}) (resultsSlice []map[string][]byte, err error)

	// insert one or more beans
	Insert(beans ...interface{}) (int64, error)

	// Insert multiple records
	InsertMulti(rowsSlicePtr interface{}) (int64, error)

	// Method InsertOne insert only one struct into database as a record.
	// The in parameter bean must a struct or a point to struct. The return
	// parameter is inserted and error
	InsertOne(bean interface{}) (int64, error)

	// Update records, bean's non-empty fields are updated contents,
	// condiBean' non-empty filds are conditions
	// CAUTION:
	//        1.bool will defaultly be updated content nor conditions
	//         You should call UseBool if you have bool to use.
	//        2.float32 & float64 may be not inexact as conditions
	Update(bean interface{}, condiBean ...interface{}) (int64, error)

	// Delete records, bean's non-empty fields are conditions
	Delete(bean interface{}) (int64, error)

	// LastSQL returns last query information
	LastSQL() (string, []interface{})

	Sync2(beans ...interface{}) error

	// Always disable struct tag "deleted"
	Unscoped() Ormer
}
