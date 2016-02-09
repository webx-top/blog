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
	//"fmt"
	"math"
	"strings"

	//"github.com/webx-top/echo"
	"github.com/coscms/xorm/core"
	"github.com/webx-top/blog/app/base/lib/database"
	X "github.com/webx-top/webx"
	"github.com/webx-top/webx/lib/com"
)

func New(c *X.Context, orm *database.Orm, m interface{}) *DataTable {
	a := &DataTable{Context: c, Orm: orm}
	a.PageSize = com.Int64(c.Form(`length`))
	a.Offset = com.Int64(c.Form(`start`))
	if a.PageSize < 1 || a.PageSize > 1000 {
		a.PageSize = 10
	}
	a.Fields = make([]string, 0)
	a.TableFields = make([]string, 0)
	a.Searches = make(map[string]string)
	a.Page = (a.Offset + a.PageSize) / a.PageSize
	a.Orders = Sorts{}
	a.fieldsInfo = make(map[string]*core.Column)
	a.Search = c.Form(`search[value]`)
	a.Draw = c.Form(`draw`)
	table := orm.TableInfo(m)
	if table == nil {
		return a
	}
	pks := table.PKColumns()
	if len(pks) > 0 {
		for _, col := range pks {
			if col.IsPrimaryKey && col.IsAutoIncrement {
				a.IdFieldName = col.Name
				break
			}
		}
	}
	var fm []string = strings.Split(`columns[0][data]`, `0`)
	c.AutoParseForm()
	for k, _ := range c.Request().Form {
		if !strings.HasPrefix(k, fm[0]) || !strings.HasSuffix(k, fm[1]) {
			continue
		}
		idx := strings.TrimSuffix(k, fm[1])
		idx = strings.TrimPrefix(idx, fm[0])

		//要查询的所有字段
		field := c.Form(k)

		column := table.GetColumn(field)
		if column != nil && column.FieldName == field {
			a.Fields = append(a.Fields, field)
			field = a.Orm.Engine.ColumnMapper.Obj2Table(field)
			a.TableFields = append(a.TableFields, field)

			//搜索本字段
			kw := c.Form(`columns[` + idx + `][search][value]`)
			if kw != `` {
				a.Searches[field] = kw
			}
			a.fieldsInfo[field] = column
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
		column = table.GetColumn(field)
		if column != nil && column.FieldName == field {
			continue
		}
		sort := c.Form(`order[` + idx + `][dir]`)
		if sort != `asc` {
			sort = `desc`
		}
		a.Orders.Insert(com.Int(idx), field, a.Orm.Engine.ColumnMapper.Obj2Table(field), sort)
	}
	a.OrderBy = a.Orders.Sql()
	//a.Form(`search[regex]`)=="false"
	//columns[0][search][regex]=false / columns[0][search][value]
	return a
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
	Draw        string //DataTabels发起的请求标识
	PageSize    int64  //每页数据量
	Page        int64
	Offset      int64             //数据偏移值
	Fields      []string          //查询的字段
	TableFields []string          //查询的字段
	Orders      Sorts             //字段和排序方式
	OrderBy     string            //ORDER BY 语句
	Search      string            //搜索关键字
	Searches    map[string]string //搜索某个字段
	totalPages  int64             //总页数
	IdFieldName string
	fieldsInfo  map[string]*core.Column
}

//总页数
func (a *DataTable) Pages(totalRows int64) int64 {
	if totalRows <= 0 {
		a.totalPages = 1
	} else {
		a.totalPages = int64(math.Ceil(float64(totalRows) / float64(a.PageSize)))
	}
	return a.totalPages
}

//结果数据
func (a *DataTable) Data(totalRows int64, data interface{}) (r *map[string]interface{}) {
	r = &map[string]interface{}{
		"draw":            a.Draw,
		"recordsTotal":    totalRows,
		"recordsFiltered": totalRows,
		"data":            data,
	}
	a.Context.AssignX(r)
	return
}

//生成 ORDER BY 子句
func (a *DataTable) GenOrderBy(args ...func(string, string) string) string {
	return a.Orders.Sql(args...)
}

//生成搜索条件
func (a *DataTable) GenSearch(fields ...string) string {
	var sql string
	var lnk string
	var sqle, lnke string
	for field, keywords := range a.Searches {
		var cond string
		column, ok := a.fieldsInfo[field]
		if !ok {
			continue
		}
		if column.SQLType.IsText() {
			switch column.SQLType.Name {
			case core.Enum, core.Set, core.Char, core.Uuid:
				cond = a.Orm.EqField(field, keywords)
				if cond != `` {
					sqle += lnke + cond
					lnke = ` AND `
				}
			default:
				cond = a.Orm.SearchField(field, keywords, a.IdFieldName, a.IdFieldName)
				if cond != `` {
					sql += lnk + cond
					lnk = ` AND `
				}
			}
		} else if column.SQLType.IsNumeric() {
			switch column.SQLType.Name {
			case core.Bool, core.Serial, core.BigSerial:
				cond = a.Orm.EqField(field, keywords)
				if cond != `` {
					sqle += lnke + cond
					lnke = ` AND `
				}
			default:
				cond = a.Orm.RangeField(field, keywords)
				if cond != `` {
					sql += lnk + cond
					lnk = ` AND `
				}
			}
		}
	}
	if sqle != `` {
		if sql != `` {
			sql = sqle + lnke + sql
		} else {
			sql = sqle
		}
	}
	if sql != `` {
		sql = `(` + sql + `)`
	}
	if a.Search != `` && len(fields) > 0 {
		cond := a.Orm.SearchField(strings.Join(fields, `,`), a.Search, a.IdFieldName)
		if cond != `` {
			if sql != `` {
				sql += ` AND ` + cond
			} else {
				sql = cond
			}
		}
	}
	return sql
}
