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
package datatable

import (
	"fmt"
	"math"
	"strings"

	"github.com/coscms/xorm/core"
	"github.com/webx-top/com"
	X "github.com/webx-top/webx"
	listClient "github.com/webx-top/webx/lib/client/list"
	"github.com/webx-top/webx/lib/database"
)

func init() {
	_ = fmt.Sprint
	listClient.Reg(`dataTable`, func() listClient.Client {
		return New()
	})
}

func New() listClient.Client {
	return &DataTable{
		autoSearchPk: true,
	}
}

type Sort struct {
	Field      string
	TableField string
	Sort       string
}

type Sorts []*Sort

func (a Sorts) Each(f func(string, string)) {
	for _, v := range a {
		if v != nil {
			f(v.TableField, v.Sort)
		}
	}
}

func (a *Sorts) Insert(index int, field string, tableField string, sort string) {
	length := len(*a)
	if length > index {
		(*a)[index] = &Sort{Field: field, TableField: tableField, Sort: sort}
	} else if index <= 10 {
		for i := length; i <= index; i++ {
			if i == index {
				*a = append(*a, &Sort{Field: field, TableField: tableField, Sort: sort})
			} else {
				*a = append(*a, nil)
			}
		}
	}
}

func (a Sorts) Sql(args ...func(string, string) string) (r string) {
	var fn func(string, string) string
	if len(args) > 0 {
		fn = args[0]
	} else {
		fn = func(field string, sort string) string {
			return field + ` ` + sort
		}
	}
	a.Each(func(field string, sort string) {
		r += fn(field, sort) + `,`
	})

	if r != `` {
		r = strings.TrimSuffix(r, `,`)
	}
	return
}

type DataTable struct {
	*X.Context
	*database.Orm
	draw         string //DataTabels发起的请求标识
	pageSize     int64  //每页数据量
	totalRows    int64
	page         int64
	offset       int64             //数据偏移值
	fields       []string          //查询的字段
	tableFields  []string          //查询的字段
	orders       Sorts             //字段和排序方式
	search       string            //搜索关键字
	searches     map[string]string //搜索某个字段
	totalPages   int64             //总页数
	pkName       string            //primary key name
	fieldsInfo   map[string]*core.Column
	countFn      func() int64
	autoSearchPk bool
}

type fieldInfo struct {
	index string
	sort  string
}

func (a *DataTable) Init(c *X.Context, orm *database.Orm, m interface{}, args ...string) listClient.Client {
	pagerowsKey := `length`
	totalrowsKey := `totalrows`
	offsetKey := `start`
	switch len(args) {
	case 3:
		pagerowsKey = args[2]
		fallthrough
	case 2:
		totalrowsKey = args[1]
		fallthrough
	case 1:
		offsetKey = args[0]
	}
	a.Context = c
	a.Orm = orm
	a.pageSize = com.Int64(c.Form(pagerowsKey))
	a.offset = com.Int64(c.Form(offsetKey))
	if a.pageSize < 1 || a.pageSize > 1000 {
		a.pageSize = 10
	}
	if a.offset < 0 {
		a.offset = 0
	}
	a.fields = make([]string, 0)
	a.tableFields = make([]string, 0)
	a.searches = make(map[string]string)
	a.page = (a.offset + a.pageSize) / a.pageSize
	a.orders = Sorts{}
	a.fieldsInfo = make(map[string]*core.Column)
	a.search = c.Form(`search[value]`)
	a.draw = c.Form(`draw`)
	a.totalRows = com.Int64(a.Context.Form(totalrowsKey))
	if a.totalRows < 1 && a.countFn != nil {
		a.totalRows = a.countFn()
	}

	fm := strings.Split(`columns[0][data]`, `0`)

	// ==========================
	// 获取客户端提交的字段名
	// ==========================
	findFields := map[string]map[string]interface{}{}
	for k, _ := range c.Request().Form().All() {
		if !strings.HasPrefix(k, fm[0]) || !strings.HasSuffix(k, fm[1]) {
			continue
		}
		idx := strings.TrimSuffix(k, fm[1])
		idx = strings.TrimPrefix(idx, fm[0])

		//要查询的所有字段
		field := c.Form(k)
		fieldParts := strings.Split(field, `.`)
		var parent string
		if len(fieldParts) > 1 {
			parent = fieldParts[0]
			field = fieldParts[1]
		}
		if pv, ok := findFields[parent]; !ok {
			findFields[parent] = map[string]interface{}{}
		} else {
			if _, ok := pv[field]; !ok {
				findFields[parent][field] = fieldInfo{index: idx, sort: ``}
			}
		}

		//要排序的字段
		fidx := c.Form(`order[` + idx + `][column]`)
		if fidx == `` {
			continue
		}
		field = c.Form(fm[0] + fidx + fm[1])
		if field == `` {
			continue
		}
		fieldParts = strings.Split(field, `.`)
		parent = ``
		if len(fieldParts) > 1 {
			parent = fieldParts[0]
			field = fieldParts[1]
		}
		sort := c.Form(`order[` + idx + `][dir]`)
		if sort != `asc` {
			sort = `desc`
		}
		if _, ok := findFields[parent]; !ok {
			findFields[parent] = map[string]interface{}{}
		}
		findFields[parent][field] = fieldInfo{index: fidx, sort: sort}
	}
	//com.Dump(findFields)
	a.pkName = orm.VerifyFieldsByMap(m, findFields, func(column *core.Column, inf interface{}, prefix string) {
		info := inf.(fieldInfo)
		fullField := prefix + column.FieldName
		colName := column.Name
		fullColName := prefix + colName

		a.fields = append(a.fields, fullField)
		a.tableFields = append(a.tableFields, fullColName)

		//搜索本字段
		kw := c.Form(`columns[` + info.index + `][search][value]`)
		if kw != `` {
			a.searches[fullColName] = kw
		}
		a.fieldsInfo[fullColName] = column
		if info.sort != `` {
			a.orders.Insert(com.Int(info.index), fullField, fullColName, info.sort)
		}
	})
	//com.Dump(a.searches)
	//a.Form(`search[regex]`)=="false"
	//columns[0][search][regex]=false / columns[0][search][value]
	return a
}

func (a *DataTable) SetCount(fn func() int64) listClient.Client {
	a.countFn = fn
	if a.totalRows < 1 && a.countFn != nil {
		a.totalRows = a.countFn()
	}
	return a
}

func (a *DataTable) PageSize() int64 {
	return a.pageSize
}

func (a *DataTable) Offset() int64 {
	return a.offset
}

// Pages 总页数
func (a *DataTable) Pages() int64 {
	if a.totalRows <= 0 {
		a.totalPages = 1
	} else {
		a.totalPages = int64(math.Ceil(float64(a.totalRows) / float64(a.pageSize)))
	}
	return a.totalPages
}

// Data 结果数据
func (a *DataTable) Data(data interface{}, args ...string) (r *map[string]interface{}) {
	r = &map[string]interface{}{
		"draw":            a.draw,
		"recordsTotal":    a.totalRows,
		"recordsFiltered": a.totalRows,
		"data":            data,
	}
	if len(args) > 0 {
		a.Context.Assign(args[0], r)
	} else {
		a.Context.AssignX(r)
	}
	return
}

// GenOrderBy 生成 ORDER BY 子句
func (a *DataTable) GenOrderBy(args ...func(string, string) string) string {
	var fn = func(field string, sort string) string {
		return a.Orm.Quote(field) + ` ` + sort
	}
	if len(args) > 0 {
		fn = args[0]
	}
	return a.orders.Sql(fn)
}

func (a *DataTable) SearchPk(on bool) listClient.Client {
	a.autoSearchPk = on
	return a
}

// GenSearch 生成搜索条件
func (a *DataTable) GenSearch(fields ...string) string {
	var sql, lnk, sqle, lnke string
	for field, keywords := range a.searches {
		var cond string
		column, ok := a.fieldsInfo[field]
		if !ok {
			continue
		}
		if column.SQLType.IsText() {
			switch column.SQLType.Name {
			case core.Enum, core.Set, core.Char, core.Uuid:
				cond = a.Orm.EqField(field, keywords)
			default:
				cond = a.Orm.SearchField(field, keywords)
			}
			if len(cond) > 0 {
				sql += lnk + cond
				lnk = ` AND `
			}
		} else if column.SQLType.IsNumeric() {
			switch column.SQLType.Name {
			case core.Bool, core.Serial, core.BigSerial:
				cond = a.Orm.EqField(field, keywords)
			//case core.TinyInt: cond = a.Orm.RangeField(field, keywords)
			default:
				if strings.Contains(keywords, ` - `) {
					a.Orm.GenDateRangeSql(&sql, field, keywords)
				} else if database.IsCompareField(keywords) {
					cond = a.Orm.CompareField(field, keywords)
				} else {
					cond = a.Orm.RangeField(field, keywords)
				}
			}

			if len(cond) > 0 {
				sql += lnk + cond
				lnk = ` AND `
			}
		}
	}
	if len(sqle) > 0 {
		if len(sql) > 0 {
			sql = sqle + lnke + sql
		} else {
			sql = sqle
		}
	}
	if len(sql) > 0 {
		sql = `(` + sql + `)`
	}
	if len(a.search) > 0 && len(fields) > 0 {
		var cond string
		if a.autoSearchPk {
			cond = a.Orm.SearchFields(fields, a.search, a.pkName)
		} else {
			cond = a.Orm.SearchFields(fields, a.search)
		}
		if len(cond) > 0 {
			if len(sql) > 0 {
				sql += ` AND ` + cond
			} else {
				sql = cond
			}
		}
	}
	return sql
}
