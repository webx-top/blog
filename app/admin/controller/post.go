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
package controller

import (
	//"fmt"
	"strings"

	"github.com/webx-top/blog/app/admin/lib"
	X "github.com/webx-top/webx"
	"github.com/webx-top/webx/lib/com"
)

func init() {
	lib.App.Reg(&Post{}).Auto()
}

type Post struct {
	index X.Mapper
	*Base
}

func (a *Post) Init(c *X.Context) error {
	a.Base = New(c)
	return nil
}

func (a *Post) Index() error {
	var pageSize int = com.Int(a.Form(`length`))
	var offset int = com.Int(a.Form(`start`))
	if pageSize < 1 || pageSize > 1000 {
		pageSize = 10
	}
	var fields []string = make([]string, 0)
	var orderBy string
	var fm []string = strings.Split(`columns[0][data]`, `0`)
	a.AutoParseForm()
	for k, _ := range a.Request().Form {
		if !strings.HasPrefix(k, fm[0]) || !strings.HasSuffix(k, fm[1]) {
			continue
		}
		//要查询的所有字段
		field := a.Form(k)
		fields = append(fields, field)

		//要排序的字段
		idx := strings.TrimSuffix(k, fm[1])
		idx = strings.TrimPrefix(idx, fm[0])

		fidx := a.Form(`order[` + idx + `][column]`)
		if fidx != `` {
			field := a.Form(fm[0] + fidx + fm[1])
			if field == `` {
				continue
			}
			sort := a.Form(`order[` + idx + `][dir]`)
			if sort != `asc` {
				sort = `desc`
			}
			orderBy += field + ` ` + sort + `,`
		}
	}
	if orderBy != `` {
		orderBy = strings.TrimSuffix(orderBy, `,`)
	}
	var search string = a.Form(`search[value]`)
	_ = search
	_ = offset
	//a.Form(`search[regex]`)=="false"
	a.AssignX(&map[string]interface{}{
		"draw":            a.Form(`draw`),
		"recordsTotal":    100000000,
		"recordsFiltered": 50000,
		"data": []map[string]string{
			map[string]string{"id": "id", "title": "title1", "priority": "priority1", "status": "status1"},
		},
	})
	//columns[0][search][regex]=false / columns[0][search][value]
	return nil
}
