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
	"math"
	"strings"

	//"github.com/webx-top/echo"
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
	a.Page = (a.Offset + a.PageSize) / a.PageSize
	a.Orders = Sorts{}
	var fm []string = strings.Split(`columns[0][data]`, `0`)
	c.AutoParseForm()
	for k, _ := range c.Request().Form {
		if !strings.HasPrefix(k, fm[0]) || !strings.HasSuffix(k, fm[1]) {
			continue
		}
		//要查询的所有字段
		field := c.Form(k)
		if a.Orm.VerifyField(m, field) != `` {
			a.Fields = append(a.Fields, field)
		}

		//要排序的字段
		idx := strings.TrimSuffix(k, fm[1])
		idx = strings.TrimPrefix(idx, fm[0])

		fidx := c.Form(`order[` + idx + `][column]`)
		if fidx == `` {
			continue
		}
		field = c.Form(fm[0] + fidx + fm[1])
		if field == `` {
			continue
		}
		if a.Orm.VerifyField(m, field) == `` {
			continue
		}
		sort := c.Form(`order[` + idx + `][dir]`)
		if sort != `asc` {
			sort = `desc`
		}
		a.Orders.Insert(com.Int(idx), field, sort)
	}
	a.OrderBy = a.Orders.Sql()
	a.Search = c.Form(`search[value]`)
	a.Draw = c.Form(`draw`)
	//a.Form(`search[regex]`)=="false"
	//columns[0][search][regex]=false / columns[0][search][value]
	return a
}

type Sort struct {
	Field string
	Sort  string
}

type Sorts []*Sort

func (a Sorts) Each(f func(string, string)) {
	for _, v := range a {
		if v != nil {
			f(v.Field, v.Sort)
		}
	}
}

func (a *Sorts) Insert(index int, field string, sort string) {
	length := len(*a)
	if length > index {
		(*a)[index] = &Sort{Field: field, Sort: sort}
	} else if index <= 10 {
		for i := length; i <= index; i++ {
			if i == index {
				*a = append(*a, &Sort{Field: field, Sort: sort})
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
	Draw       string //DataTabels发起的请求标识
	PageSize   int64  //每页数据量
	Page       int64
	Offset     int64    //数据偏移值
	Fields     []string //查询的字段
	Orders     Sorts    //字段和排序方式
	OrderBy    string   //ORDER BY 语句
	Search     string   //搜索关键字
	totalPages int64    //总页数
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
