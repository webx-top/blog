// Copyright 2015 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	"hash/crc32"
	"reflect"

	"github.com/coscms/xorm/core"
)

// Session keep a pointer to sql.DB and provides all execution of all
// kind of database operations.
type Session struct {
	db                     *core.DB
	Engine                 *Engine
	Tx                     *core.Tx
	Statement              Statement
	IsAutoCommit           bool
	IsCommitedOrRollbacked bool
	TransType              string
	IsAutoClose            bool

	// Automatically reset the statement after operations that execute a SQL
	// query such as Count(), Find(), Get(), ...
	AutoResetStatement bool

	// !nashtsai! storing these beans due to yet committed tx
	afterInsertBeans map[interface{}]*[]func(interface{})
	afterUpdateBeans map[interface{}]*[]func(interface{})
	afterDeleteBeans map[interface{}]*[]func(interface{})
	// --

	beforeClosures []func(interface{})
	afterClosures  []func(interface{})

	prepareStmt bool
	stmtCache   map[uint32]*core.Stmt //key: hash.Hash32 of (queryStr, len(queryStr))
	cascadeDeep int

	// !evalphobia! stored the last executed query on this session
	//beforeSQLExec func(string, ...interface{})
	lastSQL     string
	lastSQLArgs []interface{}
}

// Clone copy all the session's content and return a new session
func (session *Session) Clone() *Session {
	var sess = *session
	return &sess
}

// Init reset the session as the init status.
func (session *Session) Init() {
	session.Statement.Init()
	session.Statement.Engine = session.Engine
	session.IsAutoCommit = true
	session.IsCommitedOrRollbacked = false
	session.IsAutoClose = false
	session.AutoResetStatement = true
	session.prepareStmt = false

	// !nashtsai! is lazy init better?
	session.afterInsertBeans = make(map[interface{}]*[]func(interface{}), 0)
	session.afterUpdateBeans = make(map[interface{}]*[]func(interface{}), 0)
	session.afterDeleteBeans = make(map[interface{}]*[]func(interface{}), 0)
	session.beforeClosures = make([]func(interface{}), 0)
	session.afterClosures = make([]func(interface{}), 0)

	session.lastSQLArgs = []interface{}{}
}

// Close release the connection from pool
func (session *Session) Close() {
	for _, v := range session.stmtCache {
		v.Close()
	}

	if session.db != nil {
		// When Close be called, if session is a transaction and do not call
		// Commit or Rollback, then call Rollback.
		if session.Tx != nil && !session.IsCommitedOrRollbacked {
			session.Rollback()
		}
		session.Tx = nil
		session.stmtCache = nil
		session.db = nil

		session.Statement.Init()
		session.IsAutoCommit = true
		session.IsCommitedOrRollbacked = false
		session.IsAutoClose = false
		session.AutoResetStatement = true
		session.prepareStmt = false

		// processors
		session.afterInsertBeans = nil
		session.afterUpdateBeans = nil
		session.afterDeleteBeans = nil
		session.beforeClosures = nil
		session.afterClosures = nil
	}
}

func (session *Session) resetStatement() {
	if session.AutoResetStatement {
		session.Statement.Init()
	}
}

// Before Apply before Processor, affected bean is passed to closure arg
func (session *Session) Before(closures func(interface{})) *Session {
	if closures != nil {
		session.beforeClosures = append(session.beforeClosures, closures)
	}
	return session
}

// After Apply after Processor, affected bean is passed to closure arg
func (session *Session) After(closures func(interface{})) *Session {
	if closures != nil {
		session.afterClosures = append(session.afterClosures, closures)
	}
	return session
}

// Table can input a string or pointer to struct for special a table to operate.
func (session *Session) Table(tableNameOrBean interface{}) *Session {
	session.Statement.Table(tableNameOrBean)
	return session
}

// Alias set the table alias
func (session *Session) Alias(alias string) *Session {
	session.Statement.Alias(alias)
	return session
}

// Where provides custom query condition.
func (session *Session) Where(query interface{}, args ...interface{}) *Session {
	session.Statement.Where(query, args...)
	return session
}

// And provides custom query condition.
func (session *Session) And(query interface{}, args ...interface{}) *Session {
	session.Statement.And(query, args...)
	return session
}

// Or provides custom query condition.
func (session *Session) Or(query interface{}, args ...interface{}) *Session {
	session.Statement.Or(query, args...)
	return session
}

// Id will be deprecated, please use ID instead
func (session *Session) Id(id interface{}) *Session {
	session.Statement.Id(id)
	return session
}

// ID provides converting id as a query condition
func (session *Session) ID(id interface{}) *Session {
	session.Statement.Id(id)
	return session
}

// In provides a query string like "id in (1, 2, 3)"
func (session *Session) In(column string, args ...interface{}) *Session {
	session.Statement.In(column, args...)
	return session
}

// NotIn provides a query string like "id in (1, 2, 3)"
func (session *Session) NotIn(column string, args ...interface{}) *Session {
	session.Statement.NotIn(column, args...)
	return session
}

// Incr provides a query string like "count = count + 1"
func (session *Session) Incr(column string, arg ...interface{}) *Session {
	session.Statement.Incr(column, arg...)
	return session
}

// Decr provides a query string like "count = count - 1"
func (session *Session) Decr(column string, arg ...interface{}) *Session {
	session.Statement.Decr(column, arg...)
	return session
}

// SetExpr provides a query string like "column = {expression}"
func (session *Session) SetExpr(column string, expression string) *Session {
	session.Statement.SetExpr(column, expression)
	return session
}

// Select provides some columns to special
func (session *Session) Select(str string) *Session {
	session.Statement.Select(str)
	return session
}

// Cols provides some columns to special
func (session *Session) Cols(columns ...string) *Session {
	session.Statement.Cols(columns...)
	return session
}

// AllCols ask all columns
func (session *Session) AllCols() *Session {
	session.Statement.AllCols()
	return session
}

// MustCols specify some columns must use even if they are empty
func (session *Session) MustCols(columns ...string) *Session {
	session.Statement.MustCols(columns...)
	return session
}

// Distinct use for distinct columns. Caution: when you are using cache,
// distinct will not be cached because cache system need id,
// but distinct will not provide id
func (session *Session) Distinct(columns ...string) *Session {
	session.Statement.Distinct(columns...)
	return session
}

// ForUpdate Set Read/Write locking for UPDATE
func (session *Session) ForUpdate() *Session {
	session.Statement.IsForUpdate = true
	return session
}

// Limit provide limit and offset query condition
func (session *Session) Limit(limit int, start ...int) *Session {
	session.Statement.Limit(limit, start...)
	return session
}

// OrderBy provide order by query condition, the input parameter is the content
// after order by on a sql statement.
func (session *Session) OrderBy(order string) *Session {
	session.Statement.OrderBy(order)
	return session
}

// Desc provide desc order by query condition, the input parameters are columns.
func (session *Session) Desc(colNames ...string) *Session {
	session.Statement.Desc(colNames...)
	return session
}

// Asc provide asc order by query condition, the input parameters are columns.
func (session *Session) Asc(colNames ...string) *Session {
	session.Statement.Asc(colNames...)
	return session
}

// Join join_operator should be one of INNER, LEFT OUTER, CROSS etc - this will be prepended to JOIN
func (session *Session) Join(joinOperator string, tablename interface{}, condition string, args ...interface{}) *Session {
	session.Statement.Join(joinOperator, tablename, condition, args...)
	return session
}

// GroupBy Generate Group By statement
func (session *Session) GroupBy(keys string) *Session {
	session.Statement.GroupBy(keys)
	return session
}

// Having Generate Having statement
func (session *Session) Having(conditions string) *Session {
	session.Statement.Having(conditions)
	return session
}

// DB db return the wrapper of sql.DB
func (session *Session) DB() *core.DB {
	if session.db == nil {
		session.db = session.Engine.db
		session.stmtCache = make(map[uint32]*core.Stmt, 0)
	}
	return session.db
}

func (session *Session) canCache() bool {
	if session.Statement.RefTable == nil ||
		session.Statement.JoinStr() != "" ||
		session.Statement.RawSQL != "" ||
		session.Tx != nil ||
		len(session.Statement.selectStr) > 0 {
		return false
	}
	return true
}

func (session *Session) doPrepare(sqlStr string) (stmt *core.Stmt, err error) {
	crc := crc32.ChecksumIEEE([]byte(sqlStr))
	// TODO try hash(sqlStr+len(sqlStr))
	var has bool
	stmt, has = session.stmtCache[crc]
	if !has {
		stmt, err = session.DB().Prepare(sqlStr)
		if err != nil {
			return nil, err
		}
		session.stmtCache[crc] = stmt
	}
	return
}

// Ping test if database is ok
func (session *Session) Ping() error {
	defer session.resetStatement()
	if session.IsAutoClose {
		defer session.Close()
	}

	return session.DB().Ping()
}

func (session *Session) getField(dataStruct *reflect.Value, key string, table *core.Table, idx int) *reflect.Value {
	var col *core.Column
	if col = table.GetColumnIdx(key, idx); col == nil {
		//session.Engine.logger.Warnf("table %v has no column %v. %v", table.Name, key, table.ColumnsSeq())
		return nil
	}

	fieldValue, err := col.ValueOfV(dataStruct)
	if err != nil {
		session.Engine.logger.Error(err)
		return nil
	}

	if !fieldValue.IsValid() || !fieldValue.CanSet() {
		session.Engine.logger.Warnf("table %v's column %v is not valid or cannot set", table.Name, key)
		return nil
	}
	return fieldValue
}

// saveLastSQL stores executed query information
func (session *Session) saveLastSQL(sql string, args ...interface{}) {
	session.lastSQL = sql
	session.lastSQLArgs = args
	session.Engine.logSQL(sql, args...)
}

// LastSQL returns last query information
func (session *Session) LastSQL() (string, []interface{}) {
	return session.lastSQL, session.lastSQLArgs
}

// tbName get some table's table name
func (session *Session) tbNameNoSchema(table *core.Table) string {
	if len(session.Statement.AltTableName) > 0 {
		return session.Statement.AltTableName
	}

	return table.Name
}
