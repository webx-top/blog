package xorm

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/coscms/xorm/core"
)

type SelectRows []*ResultSet

func (this SelectRows) GetRow() (result *ResultSet) {
	if len(this) > 0 {
		result = this[0]
	}
	return
}

func (this SelectRows) GetOne() (result string) {
	results := this.GetRow()
	if results != nil {
		result = results.GetString(0)
	}
	return
}

// =====================================
// 定义ResultSet
// =====================================
func NewResultSet() *ResultSet {
	return &ResultSet{
		Fields:    make([]string, 0),
		Values:    make([]*reflect.Value, 0),
		NameIndex: make(map[string]int),
		Length:    0,
	}
}

type ResultSet struct {
	Fields    []string
	Values    []*reflect.Value
	NameIndex map[string]int
	Length    int
}

func (r *ResultSet) Get(index int) interface{} {
	if index >= r.Length {
		return nil
	}
	return (*r.Values[index]).Interface()
}

func (r *ResultSet) GetByName(name string) interface{} {
	if index, ok := r.NameIndex[name]; ok {
		return r.Get(index)
	}
	return nil
}

func (r *ResultSet) GetString(index int) string {
	if index >= r.Length {
		return ""
	}
	str, err := reflect2value(r.Values[index])
	if err != nil {
		log.Println(err)
	}
	return str
}

func (r *ResultSet) GetStringByName(name string) string {
	if index, ok := r.NameIndex[name]; ok {
		return r.GetString(index)
	}
	return ""
}

func (r *ResultSet) GetInt64(index int) int64 {
	v := r.Get(index)
	if v == nil {
		return 0
	}
	if val, ok := v.(int64); ok {
		return val
	}
	if val, ok := v.(float64); ok {
		return int64(val)
	}
	if val, ok := v.(int); ok {
		return int64(val)
	}
	return 0
}

func (r *ResultSet) GetInt64ByName(name string) int64 {
	if index, ok := r.NameIndex[name]; ok {
		return r.GetInt64(index)
	}
	return 0
}

func (r *ResultSet) GetInt(index int) int {
	return int(r.GetInt64(index))
}

func (r *ResultSet) GetIntByName(name string) int {
	return int(r.GetInt64ByName(name))
}

func (r *ResultSet) GetFloat64(index int) float64 {
	v := r.Get(index)
	if v == nil {
		return 0
	}
	if val, ok := v.([]uint8); ok {
		r, _ := strconv.ParseFloat(string(val), 64)
		return r
	}
	if val, ok := v.(float64); ok {
		return val
	}
	if val, ok := v.(float32); ok {
		return float64(val)
	}
	return 0
}

func (r *ResultSet) GetFloat64ByName(name string) float64 {
	if index, ok := r.NameIndex[name]; ok {
		return r.GetFloat64(index)
	}
	return 0
}

func (r *ResultSet) GetFloat32(index int) float32 {
	return float32(r.GetFloat64(index))
}

func (r *ResultSet) GetFloat32ByName(name string) float32 {
	return float32(r.GetFloat64ByName(name))
}

func (r *ResultSet) GetBool(index int) bool {
	v := r.Get(index)
	if v == nil {
		return false
	}
	if val, ok := v.(bool); ok {
		return val
	}
	return false
}

func (r *ResultSet) GetBoolByName(name string) bool {
	if index, ok := r.NameIndex[name]; ok {
		return r.GetBool(index)
	}
	return false
}

func (r *ResultSet) GetTime(index int) time.Time {
	var t time.Time
	if v := r.Get(index); v != nil {
		t, _ = v.(time.Time)
	}
	return t
}

func (r *ResultSet) GetTimeByName(name string) time.Time {
	var t time.Time
	if index, ok := r.NameIndex[name]; ok {
		t = r.GetTime(index)
	}
	return t
}

func (r *ResultSet) Set(index int, value interface{}) bool {
	if index >= r.Length {
		return false
	}
	rawValue := reflect.Indirect(reflect.ValueOf(value))
	r.Values[index] = &rawValue
	return true
}

func (r *ResultSet) SetByName(name string, value interface{}) bool {
	if index, ok := r.NameIndex[name]; ok {
		return r.Set(index, value)
	} else {
		r.NameIndex[name] = len(r.Values)
		r.Fields = append(r.Fields, name)
		rawValue := reflect.Indirect(reflect.ValueOf(value))
		r.Values = append(r.Values, &rawValue)
		r.Length = len(r.Values)
	}
	return true
}

// =====================================
// 增加Session结构体中的方法
// =====================================
func (session *Session) QueryStr(sqlStr string, paramStr ...interface{}) ([]map[string]string, error) {
	defer session.resetStatement()
	if session.IsAutoClose {
		defer session.Close()
	}

	return session.queryStr(sqlStr, paramStr...)
}

func (session *Session) QueryRaw(sqlStr string, paramStr ...interface{}) ([]map[string]interface{}, error) {
	defer session.resetStatement()
	if session.IsAutoClose {
		defer session.Close()
	}

	return session.queryInterface(sqlStr, paramStr...)
}

/**
 * Exec a raw sql and return records as []*ResultSet
 * @param  string					SQL
 * @param  ...interface{}			params
 * @return []*ResultSet,error
 * @author AdamShen (swh@admpub.com)
 */
func (session *Session) Q(sqlStr string, paramStr ...interface{}) (resultsSlice []*ResultSet, err error) {

	defer session.resetStatement()
	if session.IsAutoClose {
		defer session.Close()
	}

	resultsSlice = make([]*ResultSet, 0)
	rows, err := session.queryRows(sqlStr, paramStr...)
	if rows != nil {
		if err == nil {
			resultsSlice, err = rows2ResultSetSlice(rows)
		}
		rows.Close()
	}
	return
}

/**
 * 逐行执行回调函数
 * @param  func(*core.Rows) callback		callback func
 * @param  string sqlStr 					SQL
 * @param  ...interface{} paramStr			params
 * @return error
 * @author AdamShen (swh@admpub.com)
 * @example
 * QCallback(func(rows *core.Rows){
 * 	if err := rows.Scan(bean); err != nil {
 *		return
 *	}
 *	//.....
 * },"SELECT * FROM shop WHERE type=?","vip")
 */
func (session *Session) QCallback(callback func(*core.Rows, []string), sqlStr string, paramStr ...interface{}) (err error) {

	defer session.resetStatement()
	if session.IsAutoClose {
		defer session.Close()
	}

	rows, err := session.queryRows(sqlStr, paramStr...)
	if rows != nil {
		if err == nil {
			var fields []string
			fields, err = rows.Columns()
			if err != nil {
				return err
			}
			for rows.Next() {
				callback(rows, fields)
			}
		}
		rows.Close()
	}
	return
}

func (session *Session) queryInterface(sqlStr string, params ...interface{}) (result []map[string]interface{}, err error) {
	rows, err := session.queryRows(sqlStr, params...)
	if err != nil {
		return nil, err
	}
	result, err = rows2MapInterface(rows)
	rows.Close()
	return
}

func (session *Session) queryRows(sqlStr string, paramStr ...interface{}) (rows *core.Rows, err error) {
	session.queryPreprocess(&sqlStr, paramStr...)

	if session.IsAutoCommit {
		return session.innerQueryRows(session.DB(), sqlStr, paramStr...)
	}
	return session.txQueryRows(session.Tx, sqlStr, paramStr...)
}

func (session *Session) txQueryRows(tx *core.Tx, sqlStr string, params ...interface{}) (rows *core.Rows, err error) {
	rows, err = tx.Query(sqlStr, params...)
	if err != nil {
		return nil, err
	}
	return
}

func (session *Session) innerQueryRows(db *core.DB, sqlStr string, params ...interface{}) (rows *core.Rows, err error) {
	stmt, rows, err := session.Engine.logSQLQueryTime(sqlStr, params, func() (*core.Stmt, *core.Rows, error) {
		stmt, err := db.Prepare(sqlStr)
		if err != nil {
			return stmt, nil, err
		}
		rows, err := stmt.Query(params...)

		return stmt, rows, err
	})
	if stmt != nil {
		stmt.Close()
	}
	if err != nil {
		return nil, err
	}
	return
}

// =====================================
// 增加Engine结构体中的方法
// =====================================

func (this *Engine) QueryStr(sql string, paramStr ...interface{}) []map[string]string {
	session := this.NewSession()
	defer session.Close()
	result, err := session.QueryStr(sql, paramStr...)
	if err != nil {
		this.TLogger.Base.Error(err)
	}
	return result
}

func (this *Engine) QueryRaw(sql string, paramStr ...interface{}) []map[string]interface{} {
	session := this.NewSession()
	defer session.Close()
	result, err := session.QueryRaw(sql, paramStr...)
	if err != nil {
		this.TLogger.Base.Error(err)
	}
	return result
}

// =======================
// 原生SQL查询
// =======================
func (this *Engine) RawQuery(sql string, paramStr ...interface{}) (resultsSlice []*ResultSet, err error) {
	session := this.NewSession()
	defer session.Close()
	resultsSlice, err = session.Q(sql, paramStr...)
	return
}

func (this *Engine) RawQueryCallback(callback func(*core.Rows, []string), sql string, paramStr ...interface{}) (err error) {
	session := this.NewSession()
	defer session.Close()
	err = session.QCallback(callback, sql, paramStr...)
	return
}

/**
 * 查询键值对
 */
func (this *Engine) RawQueryKv(key string, val string, sql string, paramStr ...interface{}) map[string]string {
	var results map[string]string = make(map[string]string, 0)
	err := this.RawQueryCallback(func(rows *core.Rows, fields []string) {
		var result map[string]string = make(map[string]string)
		StrRowProcessing(rows, fields, func(data string, index int, fieldName string) {
			result[fieldName] = data
		})
		if k, ok := result[key]; ok {
			if v, ok := result[val]; ok {
				results[k] = v
			}
		}
	}, sql, paramStr...)
	if err != nil {
		this.TLogger.Base.Error(err)
	}
	return results
}

func (this *Engine) RawQueryAllKvs(key string, sql string, paramStr ...interface{}) map[string][]map[string]string {
	var results map[string][]map[string]string = make(map[string][]map[string]string, 0)
	err := this.RawQueryCallback(func(rows *core.Rows, fields []string) {
		var result map[string]string = make(map[string]string)
		StrRowProcessing(rows, fields, func(data string, index int, fieldName string) {
			result[fieldName] = data
		})
		if k, ok := result[key]; ok {
			if _, ok := results[k]; !ok {
				results[k] = make([]map[string]string, 0)
			}
			results[k] = append(results[k], result)
		}
	}, sql, paramStr...)
	if err != nil {
		this.TLogger.Base.Error(err)
	}
	return results
}

// -----------------------
// ResultSet结果
// -----------------------
func (this *Engine) GetRows(sql string, params ...interface{}) []*ResultSet {
	sql = this.ReplaceTablePrefix(sql)
	result, err := this.RawQuery(sql, params...)
	if err != nil {
		this.TLogger.Base.Error(err)
	}
	return result
}

func (this *Engine) GetRow(sql string, params ...interface{}) (result *ResultSet) {
	sql = this.ReplaceTablePrefix(sql)
	results, err := this.RawQuery(sql+" LIMIT 1", params...)
	if err != nil {
		this.TLogger.Base.Error(err)
		return
	}
	if len(results) > 0 {
		result = results[0]
	}
	return
}

func (this *Engine) GetOne(sql string, params ...interface{}) (result string) {
	results := this.GetRow(sql, params...)
	if results != nil {
		result = results.GetString(0)
	}
	return
}

// RawSelect("*","member","id=?",1)
// RawSelect("*","member","status=? AND sex=?",1,1)
// RawSelect("*","`~member` a,`~order` b","a.status=? AND b.status=?",1,1)
func (this *Engine) RawSelect(fields string, table string, where string, params ...interface{}) SelectRows {
	if fields == "" {
		fields = "*"
	} else {
		fields = this.ReplaceTablePrefix(fields)
	}
	sql := `SELECT ` + fields + ` FROM ` + this.fullTableName(table) + ` WHERE ` + this.ReplaceTablePrefix(where)
	if len(params) == 1 {
		switch params[0].(type) {
		case []interface{}:
			return this.GetRows(sql, params[0].([]interface{})...)
		}
	}
	return SelectRows(this.GetRows(sql, params...))
}

// -----------------------
// map结果
// -----------------------
func (this *Engine) RawFetchAll(fields string, table string, where string, params ...interface{}) []map[string]string {
	if fields == "" {
		fields = "*"
	} else {
		fields = this.ReplaceTablePrefix(fields)
	}
	sql := `SELECT ` + fields + ` FROM ` + this.fullTableName(table) + ` WHERE ` + this.ReplaceTablePrefix(where)
	if len(params) == 1 {
		switch params[0].(type) {
		case []interface{}:
			return this.RawQueryStr(sql, params[0].([]interface{})...)
		}
	}
	return this.RawQueryStr(sql, params...)
}

func (this *Engine) RawFetch(fields string, table string, where string, params ...interface{}) (result map[string]string) {
	if fields == "" {
		fields = "*"
	} else {
		fields = this.ReplaceTablePrefix(fields)
	}
	sql := `SELECT ` + fields + ` FROM ` + this.fullTableName(table) + ` WHERE ` + this.ReplaceTablePrefix(where) + ` LIMIT 1`
	if len(params) == 1 {
		switch params[0].(type) {
		case []interface{}:
			results := this.RawQueryStr(sql, params[0].([]interface{})...)
			if len(results) > 0 {
				result = results[0]
			}
			return
		}
	}
	results := this.RawQueryStr(sql, params...)
	if len(results) > 0 {
		result = results[0]
	}
	return
}

/**
 * 查询基于指定字段值为键名的map
 */
func (this *Engine) RawQueryKvs(key string, sql string, paramStr ...interface{}) map[string]map[string]string {
	if key == "" {
		key = "id"
	}
	var results map[string]map[string]string = make(map[string]map[string]string, 0)
	err := this.RawQueryCallback(func(rows *core.Rows, fields []string) {
		var result map[string]string = make(map[string]string)
		StrRowProcessing(rows, fields, func(data string, index int, fieldName string) {
			result[fieldName] = data
		})
		if k, ok := result[key]; ok {
			results[k] = result
		}
	}, sql, paramStr...)
	if err != nil {
		this.TLogger.Base.Error(err)
	}
	return results
}

/**
 * 查询[]map[string]string
 */
func (this *Engine) RawQueryStr(sql string, paramStr ...interface{}) []map[string]string {
	var results []map[string]string = make([]map[string]string, 0)
	err := this.RawQueryCallback(func(rows *core.Rows, fields []string) {
		var result map[string]string = make(map[string]string)
		StrRowProcessing(rows, fields, func(data string, index int, fieldName string) {
			result[fieldName] = data
		})
		results = append(results, result)
	}, sql, paramStr...)
	if err != nil {
		this.TLogger.Base.Error(err)
	}
	return results
}

// -----------------------
// 写操作
// -----------------------
func (this *Engine) RawInsert(table string, sets map[string]interface{}) (lastId int64) {
	fields := ""
	values := ""
	params := make([]interface{}, 0)
	delim := ""
	for k, v := range sets {
		fields += delim + this.Quote(k)
		values += delim + "?"
		params = append(params, v)
		delim = ","
	}
	sql := `INSERT INTO ` + this.fullTableName(table) + ` (` + fields + `) VALUES (` + values + `)`
	return this.RawExec(sql, true, params...)
}

func (this *Engine) RawOnlyInsert(table string, sets map[string]interface{}) (sql.Result, error) {
	fields := ""
	values := ""
	params := make([]interface{}, 0)
	delim := ""
	for k, v := range sets {
		fields += delim + this.Quote(k)
		values += delim + "?"
		params = append(params, v)
		delim = ","
	}
	sql := `INSERT INTO ` + this.fullTableName(table) + ` (` + fields + `) VALUES (` + values + `)`
	return this.RawExecr(sql, params...)
}

func (this *Engine) RawBatchInsert(table string, multiSets []map[string]interface{}) (sql.Result, error) {
	fields := ""
	values := ""
	params := make([]interface{}, 0)
	keyIdx := map[string]int{}
	length := 0
	delim := ""
	for i, sets := range multiSets {
		innerDelim := ""
		values += delim + "("
		if i == 0 {
			idx := 0
			for k, v := range sets {
				keyIdx[k] = idx
				fields += innerDelim + this.Quote(k)
				values += innerDelim + "?"
				params = append(params, v)
				innerDelim = ","
				idx++
			}
			length = idx
		} else {
			innerParams := make([]interface{}, length)
			for k, idx := range keyIdx {
				v, _ := sets[k]
				values += innerDelim + "?"
				innerParams[idx] = v
				innerDelim = ","
			}
			params = append(params, innerParams...)
		}
		values += ")"
		delim = ","
	}
	if values == `` {
		return nil, nil
	}
	sql := `INSERT INTO ` + this.fullTableName(table) + ` (` + fields + `) VALUES ` + values
	return this.RawExecr(sql, params...)
}

func (this *Engine) RawReplace(table string, sets map[string]interface{}) int64 {
	fields := ""
	values := ""
	params := make([]interface{}, 0)
	delim := ""
	for k, v := range sets {
		fields += delim + this.Quote(k)
		values += delim + "?"
		params = append(params, v)
		delim = ","
	}
	sql := `REPLACE INTO ` + this.fullTableName(table) + ` (` + fields + `) VALUES (` + values + `)`
	return this.RawExec(sql, false, params...)
}

func (this *Engine) RawUpdate(table string, sets map[string]interface{}, where string, args ...interface{}) int64 {
	set := ""
	params := make([]interface{}, 0)
	delim := ""
	for k, v := range sets {
		set += delim + this.Quote(k) + "=?"
		params = append(params, v)
		delim = ","
	}
	if len(args) > 0 {
		isAloneSlice := false
		if len(args) == 1 {
			switch args[0].(type) {
			case []interface{}:
				params = append(params, args[0].([]interface{})...)
				isAloneSlice = true
			}
		}
		if !isAloneSlice {
			for _, v := range args {
				params = append(params, v)
			}
		}
	}
	sql := `UPDATE ` + this.fullTableName(table) + ` SET ` + set + ` WHERE ` + where

	return this.RawExec(sql, false, params...)
}

func (this *Engine) RawDelete(table string, where string, params ...interface{}) int64 {
	sql := `DELETE FROM ` + this.fullTableName(table) + ` WHERE ` + where
	if len(params) == 1 {
		switch params[0].(type) {
		case []interface{}:
			return this.RawExec(sql, false, params[0].([]interface{})...)
		}
	}
	return this.RawExec(sql, false, params...)
}

func (this *Engine) RawExec(sql string, retId bool, params ...interface{}) (affected int64) {
	if result, err := this.Exec(sql, params...); err == nil {
		if retId {
			affected, err = result.LastInsertId()
		} else {
			affected, err = result.RowsAffected()
		}
		if err != nil {
			this.TLogger.Base.Error(err)
		}
	} else {
		this.TLogger.Base.Error(err)
	}
	return
}

func (this *Engine) RawExecr(sql string, params ...interface{}) (result sql.Result, err error) {
	result, err = this.Exec(sql, params...)
	if err != nil {
		this.TLogger.Base.Error(err)
	}
	return
}

func (this *Engine) ReplaceTablePrefix(sql string) (r string) {
	r = strings.Replace(sql, "~", this.TablePrefix, -1)
	return
}

func (this *Engine) TableName(table string) string {
	return this.TablePrefix + table + this.TableSuffix
}

func (this *Engine) fullTableName(table string) string {
	if table[0] != '`' && table[0] != '~' {
		table = this.Quote(this.TableName(table))
	}
	table = this.ReplaceTablePrefix(table)
	return table
}

func (this *Engine) QuoteValue(s string) string {
	return "'" + AddSlashes(s) + "'"
}

func (this *Engine) QuoteKey(s string) string {
	return this.Quote(s)
}

// =====================================
// 函数
// =====================================
func rows2ResultSetSlice(rows *core.Rows) (resultsSlice []*ResultSet, err error) {
	fields, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		result, err := row2ResultSet(rows, fields)
		if err != nil {
			return nil, err
		}
		resultsSlice = append(resultsSlice, result)
	}

	return resultsSlice, nil
}

func row2ResultSet(rows *core.Rows, fields []string) (result *ResultSet, err error) {
	//result := make(map[string]string)
	result = NewResultSet()
	err = RowProcessing(rows, fields, func(rawValue *reflect.Value, index int, fieldName string) error {
		//if row is null then ignore
		if (*rawValue).Interface() == nil {
			return nil
		}
		result.NameIndex[fieldName] = len(result.Fields)
		result.Fields = append(result.Fields, fieldName)
		result.Values = append(result.Values, rawValue)
		return nil
	})
	result.Length = len(result.Values)
	return result, err
}

func row2Interface(rows *core.Rows, fields []string) (result map[string]interface{}, err error) {
	result = make(map[string]interface{})
	err = RowProcessing(rows, fields, func(rawValue *reflect.Value, index int, fieldName string) error {
		//if row is null then ignore
		if (*rawValue).Interface() == nil {
			return nil
		}
		result[fieldName] = (*rawValue).Interface()
		return nil
	})
	return
}

func rows2MapInterface(rows *core.Rows) (result []map[string]interface{}, err error) {
	fields, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		ret, err := row2Interface(rows, fields)
		if err != nil {
			return nil, err
		}
		result = append(result, ret)
	}
	return
}

//获取一行中每一列字符串数据
func StrRowProcessing(rows *core.Rows, fields []string, fn func(data string, index int, fieldName string)) (err error) {
	return RowProcessing(rows, fields, func(rawValue *reflect.Value, index int, fieldName string) error {
		//if row is null then ignore
		if (*rawValue).Interface() == nil {
			return nil
		}
		if data, err := value2String(rawValue); err == nil {
			fn(data, index, fieldName)
		} else {
			return err
		}
		return nil
	})
}

//获取一行中每一列reflect.Value数据
func RowProcessing(rows *core.Rows, fields []string, fn func(data *reflect.Value, index int, fieldName string) error) (err error) {
	length := len(fields)
	scanResultContainers := make([]interface{}, length)
	for i := 0; i < length; i++ {
		var resultContainer interface{}
		scanResultContainers[i] = &resultContainer
	}
	if err := rows.Scan(scanResultContainers...); err != nil {
		return err
	}
	for ii, key := range fields {
		rawValue := reflect.Indirect(reflect.ValueOf(scanResultContainers[ii]))
		fn(&rawValue, ii, key)
	}
	return nil
}

//根据core.Rows来查询结果
func getResultSliceByRows(rows *core.Rows, erre error) (resultsSlice []map[string][]byte, err error) {
	resultsSlice = make([]map[string][]byte, 0)
	if rows != nil {
		if erre == nil {
			resultsSlice, err = rows2maps(rows)
		}
		rows.Close()
	}
	return
}

//替换sql中的占位符
func BuildSqlResult(sqlStr string, args interface{}) string {
	if args, ok := args.([]interface{}); ok {
		for _, v := range args {
			val := ""
			switch v.(type) {
			case []interface{}:
				vals := v.([]interface{})
				delim := ""
				for _, v := range vals {
					rv := fmt.Sprintf("%v", v)
					rv = AddSlashes(rv)
					val += delim + "'" + rv + "'"
					delim = ","
				}
				val = strings.Replace(val, "'", `\'`, -1)
			default:
				val = fmt.Sprintf("%v", v)
				val = AddSlashes(val)
				val = "'" + val + "'"
			}
			sqlStr = strings.Replace(sqlStr, "?", val, 1)
		}
	}
	//fmt.Printf("%v\n", sqlStr)
	return sqlStr
}

func AddSlashes(s string, args ...rune) string {
	b := []rune{'\\', '\''}
	if len(args) > 0 {
		b = append(b, args...)
	}
	return AddCSlashes(s, b...)
}

func AddCSlashes(s string, b ...rune) string {
	r := []rune{}
	for _, v := range []rune(s) {
		for _, f := range b {
			if v == f {
				r = append(r, '\\')
				break
			}
		}
		r = append(r, v)
	}
	return strings.TrimRight(string(r), `\`)
}
