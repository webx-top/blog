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
	X "github.com/webx-top/webx"
	"github.com/webx-top/webx/lib/com"
)

func New(c *X.Context) *dataTable {
	a := &dataTable{Context: c}
	a.PageSize = com.Int(c.Form(`length`))
	a.Offset = com.Int(c.Form(`start`))
	if a.PageSize < 1 || a.PageSize > 1000 {
		a.PageSize = 10
	}
	a.Fields = make([]string, 0)
	a.Orders = make(map[string]string)
	var fm []string = strings.Split(`columns[0][data]`, `0`)
	c.AutoParseForm()
	for k, _ := range c.Request().Form {
		if !strings.HasPrefix(k, fm[0]) || !strings.HasSuffix(k, fm[1]) {
			continue
		}
		//要查询的所有字段
		field := c.Form(k)
		a.Fields = append(a.Fields, field)

		//要排序的字段
		idx := strings.TrimSuffix(k, fm[1])
		idx = strings.TrimPrefix(idx, fm[0])

		fidx := c.Form(`order[` + idx + `][column]`)
		if fidx != `` {
			field := c.Form(fm[0] + fidx + fm[1])
			if field == `` {
				continue
			}
			sort := c.Form(`order[` + idx + `][dir]`)
			if sort != `asc` {
				sort = `desc`
			}
			a.Orders[field] = sort
			a.OrderBy += field + ` ` + sort + `,`
		}
	}
	if a.OrderBy != `` {
		a.OrderBy = strings.TrimSuffix(a.OrderBy, `,`)
	}
	a.Search = c.Form(`search[value]`)
	a.Draw = c.Form(`draw`)
	//a.Form(`search[regex]`)=="false"
	//columns[0][search][regex]=false / columns[0][search][value]
	return a
}

type dataTable struct {
	*X.Context
	Draw       string            //DataTabels发起的请求标识
	PageSize   int               //每页数据量
	Offset     int               //数据偏移值
	Fields     []string          //查询的字段
	Orders     map[string]string //字段和排序方式
	OrderBy    string            //ORDER BY 语句
	Search     string            //搜索关键字
	totalPages int               //总页数
}

//总页数
func (a *dataTable) Pages(totalRows int) int {
	if totalRows <= 0 {
		a.totalPages = 1
	} else {
		a.totalPages = int(math.Ceil(float64(totalRows) / float64(a.PageSize)))
	}
	return a.totalPages
}

//结果数据
func (a *dataTable) Data(totalRows int, data interface{}) (r *map[string]interface{}) {
	r = &map[string]interface{}{
		"draw":            a.Draw,
		"recordsTotal":    totalRows,
		"recordsFiltered": totalRows,
		"data":            data,
	}
	a.Context.AssignX(r)
	return
}
