/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package database

import (
	//"fmt"
	"reflect"
	"strings"

	"github.com/coscms/xorm/core"
)

//验证字段
func (this *Orm) VerifyField(v interface{}, field string) string {
	if len(field) == 0 {
		return ``
	}
	fieldParts := strings.Split(field, `.`)
	count := len(fieldParts)
	if count == 1 {
		return this.SimpleVerifyField(v, field)
	}
	val := reflect.ValueOf(v)
	typ := val.Type()
	if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct {
		val = val.Elem()
	}
	vt := val.FieldByName(strings.Title(fieldParts[0]))
	if !vt.IsValid() {
		return ``
	}
	ret := this.SimpleVerifyField(vt.Interface(), fieldParts[1])
	if len(ret) == 0 {
		return ``
	}
	return field
}

//验证结构体字段并转为数据表字段
func (this *Orm) ToTableField(v interface{}, structField string) string {
	if len(structField) == 0 {
		return ``
	}
	fieldParts := strings.Split(structField, `.`)
	count := len(fieldParts)
	if count == 1 {
		return this.SimpleToTableField(v, structField)
	}
	val := reflect.ValueOf(v)
	typ := val.Type()
	if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct {
		val = val.Elem()
	}
	fieldParts[0] = strings.Title(fieldParts[0])
	vt := val.FieldByName(fieldParts[0])
	if !vt.IsValid() {
		return ``
	}
	table := this.TableInfo(v)
	if table.Table.Relation == nil {
		return ``
	}
	tableName := table.Table.Relation.AliasGetByStructField(fieldParts[0])
	if len(tableName) == 0 {
		return ``
	}
	return tableName + `.` + this.SimpleToTableField(vt.Interface(), fieldParts[1])
}

//验证结构体字段并转为数据表字段
func (this *Orm) SimpleToTableField(v interface{}, structField string) string {
	table := this.TableInfo(v)
	field := this.ColumnMapper.Obj2Table(structField)
	column := table.Table.GetColumn(field)
	if column == nil {
		return ``
	}
	if column.FieldName == structField {
		return column.Name
	}
	return ``
}

//验证模型结构体字段
func (this *Orm) SimpleVerifyField(v interface{}, field string) string {
	table := this.TableInfo(v)
	column := table.Table.GetColumn(this.ColumnMapper.Obj2Table(field))
	if column == nil {
		return ``
	}
	if column.FieldName == field {
		return field
	}
	return ``
}

//验证字段并返回有效的字段切片
func (this *Orm) VerifyFields(v interface{}, fields ...string) []string {
	ret := make([]string, 0)
	for _, field := range fields {
		field = this.VerifyField(v, field)
		if field != `` {
			ret = append(ret, field)
		}
	}
	return ret
}

//验证模型结构体字段并返回有效的字段切片
func (this *Orm) SimpleVerifyFields(v interface{}, fields ...string) []string {
	ret := make([]string, 0)
	table := this.TableInfo(v)
	for _, field := range fields {
		column := table.Table.GetColumn(this.ColumnMapper.Obj2Table(field))
		if column == nil {
			continue
		}
		if column.FieldName != field {
			continue
		}
		ret = append(ret, field)
	}
	return ret
}

//验证字段
func (this *Orm) VerifyTblField(v interface{}, field string) string {
	if field == `` {
		return ``
	}
	fieldParts := strings.Split(field, `.`)
	count := len(fieldParts)
	if count == 1 {
		return this.SimpleVerifyTblField(v, field)
	}
	val := reflect.ValueOf(v)
	typ := val.Type()
	if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct {
		val = val.Elem()
	}
	vt := val.FieldByName(strings.Title(fieldParts[0]))
	if !vt.IsValid() {
		return ``
	}
	ret := this.SimpleVerifyTblField(vt.Interface(), fieldParts[1])
	if ret == `` {
		return ``
	}
	return field
}

//验证数据表字段
func (this *Orm) SimpleVerifyTblField(v interface{}, field string) string {
	table := this.TableInfo(v)
	column := table.Table.GetColumn(field)
	if column == nil {
		return ``
	}
	if column.Name == field {
		return field
	}
	return ``
}

//验证数据表字段并返回有效的字段切片
func (this *Orm) VerifyTblFields(v interface{}, fields ...string) []string {
	ret := make([]string, 0)
	for _, field := range fields {
		field = this.VerifyTblField(v, field)
		if field != `` {
			ret = append(ret, field)
		}
	}
	return ret
}

//验证数据表字段并返回有效的字段切片
func (this *Orm) SimpleVerifyTblFields(v interface{}, fields ...string) []string {
	ret := make([]string, 0)
	table := this.TableInfo(v)
	for _, field := range fields {
		column := table.Table.GetColumn(field)
		if column == nil {
			continue
		}
		if column.Name != field {
			continue
		}
		ret = append(ret, field)
	}
	return ret
}

// ==========================
// 验证字段名
// ==========================
func (this *Orm) VerifyFieldsByMap(m interface{}, findFields map[string]map[string]interface{}, callback func(column *core.Column, info interface{}, prefix string)) string {
	tables := map[string]*core.Table{}
	var refValue reflect.Value
	var pkField string
	for parent, fields := range findFields {
		table, ok := tables[parent]
		if !ok {
			if !refValue.IsValid() {
				refValue = reflect.ValueOf(m)
				typ := refValue.Type()
				if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct {
					refValue = refValue.Elem()
				}
			}
			if parent != `` {
				vt := refValue.FieldByName(strings.Title(parent))
				if !vt.IsValid() {
					continue
				}
				if vt.Kind() == reflect.Ptr {
					if vt.IsNil() {
						vt.Set(reflect.New(vt.Type().Elem()))
					}
					vt = vt.Elem()
				}
				table = this.TableInfo(vt.Interface()).Table
			} else {
				table = this.TableInfo(m).Table
			}
			tables[parent] = table
		}
		if table == nil {
			continue
		}
		if pkField == `` {
			pks := table.PKColumns()
			if len(pks) > 0 {
				for _, col := range pks {
					if col.IsPrimaryKey && col.IsAutoIncrement {
						pkField = parent + `.` + col.Name
						break
					}
				}
			}
		}
		var prefix string
		if parent != `` {
			prefix = parent + `.`
		}
		for field, info := range fields {
			field = this.ColumnMapper.Obj2Table(field)
			column := table.GetColumn(field)
			//table.ColumnsSeq()-真实表名切片
			//column.FieldName-结构体字段名
			//column.Name-数据表列名
			//fmt.Printf("========== %#v => %v\n", column, field)
			if column != nil && column.Name == field {
				callback(column, info, prefix)
			}

		}
	}
	return pkField
}
