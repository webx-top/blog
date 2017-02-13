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
package list

import (
	"math"

	"github.com/webx-top/com"
	X "github.com/webx-top/webx"
	"github.com/webx-top/webx/lib/client/list/pagination"
	"github.com/webx-top/webx/lib/database"
)

type defaultClient struct {
	*X.Context
	*database.Orm
	pageRows   int64
	totalRows  int64
	totalPages int64
	offset     int64
	pageno     int64
	countFn    func() int64
}

func (a *defaultClient) Init(ctx *X.Context, orm *database.Orm, m interface{}, args ...string) Client {
	pagerowsKey := `pagerows`
	totalrowsKey := `totalrows`
	pageKey := `page`
	switch len(args) {
	case 3:
		pagerowsKey = args[2]
		fallthrough
	case 2:
		totalrowsKey = args[1]
		fallthrough
	case 1:
		pageKey = args[0]
	}
	a.Context = ctx
	a.Orm = orm
	a.pageRows = com.Int64(a.Context.Form(pagerowsKey))
	if a.pageRows < 1 || a.pageRows > 1000 {
		a.pageRows = 10
	}
	a.totalRows = com.Int64(a.Context.Form(totalrowsKey))
	a.pageno = com.Int64(a.Context.Form(pageKey))
	if a.pageno < 1 {
		a.pageno = 1
	}
	a.offset = (a.pageno - 1) * a.pageRows
	return a
}

func (a *defaultClient) PageSize() int64 {
	return a.pageRows
}

func (a *defaultClient) Offset() int64 {
	return a.offset
}

func (a *defaultClient) SetCount(fn func() int64) Client {
	a.countFn = fn
	if a.totalRows < 1 && a.countFn != nil {
		a.totalRows = a.countFn()
	}
	return a
}

//总页数
func (a *defaultClient) Pages() int64 {
	if a.totalRows <= 0 {
		a.totalPages = 1
	} else {
		a.totalPages = int64(math.Ceil(float64(a.totalRows) / float64(a.pageRows)))
	}
	return a.totalPages
}

//结果数据
func (a *defaultClient) Data(data interface{}, args ...string) *map[string]interface{} {
	r := &map[string]interface{}{
		"data":       data,
		"pageRows":   a.pageRows,
		"totalRows":  a.totalRows,
		"totalPages": a.Pages(),
		"offset":     a.offset,
		"page":       a.pageno,
	}
	pagebar := func(links int64, args ...string) interface{} {
		p := pagination.New(a.Context)
		p.CheckData(data).SetPage(a.totalRows, a.pageRows, a.pageno).Ready(links)
		if len(args) > 0 {
			p.SetTmpl(args[0])
		}
		return p
	}
	if len(args) > 0 {
		a.Context.Assign(args[0], r)
		a.Context.SetFunc(`Pagebar_`+args[0], pagebar)
	} else {
		a.Context.AssignX(r)
		a.Context.SetFunc(`Pagebar`, pagebar)
	}
	return r
}

//生成 ORDER BY 子句
func (a *defaultClient) GenOrderBy(args ...func(string, string) string) string {
	return ""
}

//是否自动搜索主键字段
func (a *defaultClient) SearchPk(bool) Client {
	return a
}

//生成搜索条件
func (a *defaultClient) GenSearch(args ...string) string {
	return ""
}

type Client interface {
	//初始化
	Init(*X.Context, *database.Orm, interface{}, ...string) Client

	//设置统计功能函数
	SetCount(func() int64) Client

	//每页数据量
	PageSize() int64

	//分页查询时的偏移值
	Offset() int64

	//总页数
	Pages() int64

	//结果数据
	Data(interface{}, ...string) *map[string]interface{}

	//生成 ORDER BY 子句
	GenOrderBy(...func(string, string) string) string

	//是否自动搜索主键字段
	SearchPk(bool) Client

	//生成搜索条件
	GenSearch(...string) string
}
