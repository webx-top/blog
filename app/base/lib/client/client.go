package client

import (
	"math"

	"github.com/webx-top/blog/app/base/lib/database"
	X "github.com/webx-top/webx"
	"github.com/webx-top/webx/lib/com"
)

var clients map[string]func() Client = make(map[string]func() Client)

func Reg(name string, c func() Client) {
	clients[name] = c
}

func Get(name string) Client {
	fn, _ := clients[name]
	if fn == nil {
		fn = func() Client {
			return &defaultClient{}
		}
	}
	return fn()
}

func Delete(name string) {
	if _, ok := clients[name]; ok {
		delete(clients, name)
	}
}

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

func (a *defaultClient) Init(ctx *X.Context, orm *database.Orm, m interface{}) Client {
	a.Context = ctx
	a.Orm = orm
	a.pageRows = com.Int64(a.Context.Form(`pagerows`))
	if a.pageRows < 1 || a.pageRows > 1000 {
		a.pageRows = 10
	}
	a.totalRows = com.Int64(a.Context.Form(`totalrows`))
	a.pageno = com.Int64(a.Context.Form(`page`))
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
func (a *defaultClient) Data(data interface{}) *map[string]interface{} {
	r := &map[string]interface{}{
		"data":       data,
		"pageRows":   a.pageRows,
		"totalRows":  a.totalRows,
		"totalPages": a.Pages(),
		"offset":     a.offset,
		"page":       a.pageno,
	}
	a.Context.AssignX(r)
	return r
}

//生成 ORDER BY 子句
func (a *defaultClient) GenOrderBy(args ...func(string, string) string) string {
	return ""
}

//生成搜索条件
func (a *defaultClient) GenSearch(args ...string) string {
	return ""
}

type Client interface {
	Init(*X.Context, *database.Orm, interface{}) Client

	SetCount(func() int64) Client

	PageSize() int64

	Offset() int64

	//总页数
	Pages() int64

	//结果数据
	Data(interface{}) *map[string]interface{}

	//生成 ORDER BY 子句
	GenOrderBy(...func(string, string) string) string

	//生成搜索条件
	GenSearch(...string) string
}
