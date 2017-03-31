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
package model

import (
	"errors"

	"github.com/coscms/xorm"
	"github.com/webx-top/echo/logger"
	X "github.com/webx-top/webx"
	editClient "github.com/webx-top/webx/lib/client/edit"
	listClient "github.com/webx-top/webx/lib/client/list"
	"github.com/webx-top/webx/lib/database"
)

var ErrNotModified = errors.New(`not modified`)

func New(db *database.Orm, ctx *X.Context) *Model {
	return &Model{
		DB:      db,
		Context: ctx,
		Logger:  ctx.Object().Echo().Logger(),
	}
}

type Model struct {
	DB *database.Orm
	*X.Context
	logger.Logger
}

/*
Merge 合并源数据到目标数据中
template:
	collectFn:=func(iteratorFn func(interface{}, int)){
		for idx,v:=range {{.targets}}{
			iteratorFn(v.{{.readAttr}},idx)
		}
	}
	if session,fn:=this.Merge(collectFn,`{{.relKey}}`); session!=nil {
		tm:=[]*{{.sources}}{}
		session{{.condition}}.Find(&tm)
		for _,v := range tm{
			fn(v.{{.relAttr}},func(idx int){
				{{if .saveManyToAttr}}
				{{.targets}}[idx].{{.saveManyToAttr}}=append({{.targets}}[idx].{{.saveManyToAttr}},v) //一对多
				{{else if .fromSourceAttr}}
				{{.targets}}[idx].{{.saveOneToAttr}}=v.{{.fromSourceAttr}} //一对一
				{{else}}
				{{.targets}}[idx].{{.saveOneToAttr}}=v //一对一
				{{end}}
			})
		}
	}
模板变量:
targets			- 目标实例名
readAttr		- 目标实例中关联源数据值的属性名(外键名)
relKey			- 源数据表中的主键名
sources			- 源实例名
condition		- 源查询条件
relAttr			- 源实例中被关联的属性名
saveManyToAttr	- 目标实例中接收源数据结果的属性名（一对多时使用）
saveOneToAttr	- 同上（一对一时使用）
fromSourceAttr  - 指定仅仅采用源数据的单个属性名

example:

比如：
type Source struct{
	Id int
	Name string
}
type Target struct{
	Id int
	Title string
	Rid int //保存Source.Id的值(外键)
	Source *Source
}
targets:=[]*Target{}
获取Target中Rid所关联的数据并保存到Target.Source中

	collectFn:=func(iteratorFn func(interface{}, int)){
		for idx,v:=range targets{
			iteratorFn(v.Rid,idx)
		}
	}
	if session,fn:=this.Merge(collectFn,`id`); session!=nil {
		tm:=[]*Source{}
		session.Find(&tm)
		for _,v := range tm{
			fn(v.id,func(idx int){
				targets[idx].Source=append(targets[idx].Source,v)
			})
		}
	}

@param collectFn 收集目标数据(收集器(关联值,索引值))
@param relKey    in查询字段
@return xorm会话
@return 接收源数据(关联值,遍历符合关联值的目标数据)
*/
func (this *Model) Merge(collectFn func(func(interface{}, int)), relKey string) (*xorm.Session, func(interface{}, func(int))) {
	_ids := []interface{}{}
	_map := map[interface{}][]int{}
	collectFn(func(id interface{}, idx int) {
		if _, ok := _map[id]; !ok {
			_map[id] = []int{idx}
			_ids = append(_ids, id)
		} else {
			_map[id] = append(_map[id], idx)
		}
	})
	if len(_ids) > 0 {
		ses := this.DB.In(relKey, _ids...)
		fn := func(id interface{}, mergeFn func(int)) {
			if indexes, ok := _map[id]; ok {
				for _, idx := range indexes {
					mergeFn(idx)
				}
			}
		}
		return ses, fn
	}
	return nil, nil
}

// =====================================
// TransContext
// =====================================
type TransContext struct {
	*xorm.Session
	Error error
}

func (t *TransContext) SetError(err error, affected ...int64) *TransContext {
	if err != nil {
		t.Error = err
	} else if len(affected) > 0 && affected[0] < 1 {
		t.Error = ErrNotModified
	}
	return t
}

func NewTransContext(s *xorm.Session) *TransContext {
	return &TransContext{Session: s}
}

func (this *Model) Begin() *TransContext {
	ss, ok := this.transContext()
	if ok {
		ss.Close()
	}
	ss = NewTransContext(this.DB.NewSession())
	err := ss.Begin()
	if err != nil {
		this.Logger.Error(err)
	}
	this.Context.Set(`webx:transContext`, ss)
	return ss
}

// HasBegun 事务是否已经开始
func (this *Model) HasBegun() bool {
	_, ok := this.transContext()
	return ok
}

func (this *Model) transContext() (ss *TransContext, ok bool) {
	ss, ok = this.Context.Get(`webx:transContext`).(*TransContext)
	return
}

func (this *Model) TransSession() *xorm.Session {
	if ss, ok := this.transContext(); ok {
		return ss.Session
	}
	return nil
}

// TSess transContext
func (this *Model) TSess() *TransContext {
	ss, ok := this.transContext()
	if !ok {
		return this.Begin()
	}
	return ss
}

func (this *Model) Trans(fn func() error) *database.Orm {
	ss, ok := this.transContext()
	if !ok {
		ss = this.Begin()
	} else {
		if ss.Error != nil {
			return this.DB
		}
	}
	ss.Error = fn()
	if !ok {
		this.End(ss)
	}
	return this.DB
}

func (this *Model) Sess() (sess *xorm.Session) { // transContext or Session
	ss, ok := this.transContext()
	if !ok {
		sess = this.DB.NewSession()
		sess.IsAutoClose = true
	} else {
		sess = ss.Session
	}
	return
}

func (this *Model) End(args ...*TransContext) (err error) {
	var ss *TransContext
	if len(args) > 0 {
		ss = args[0]
	} else {
		ss, _ = this.transContext()
	}
	if ss == nil {
		return nil
	}
	if ss.Error == nil {
		err = ss.Commit()
	} else {
		err = ss.Rollback()
	}
	if err != nil {
		this.Logger.Error(err)
	}
	ss.Close()
	this.Context.Delete(`webx:transContext`)
	return
}

func (this *Model) NewSelect(m interface{}) *Select {
	return NewSelect(this.DB, this.NewListClient(m))
}

func (this *Model) NewClient(m interface{}, args ...string) *Conversion {
	typ := `list`
	if len(args) > 0 {
		typ = args[0]
	}
	cv := &Conversion{}
	switch typ {
	case `edit`:
		cv.Original = this.NewEditClient(m)
	case `list`:
		fallthrough
	default:
		cv.Original = this.NewListClient(m)
	}
	return cv
}

func (this *Model) NewListClient(m interface{}) listClient.Client {
	clientName := this.Context.Form(`client`)
	if len(clientName) == 0 {
		this.Logger.Debug(`Parameter client is empty in URL: ` + this.Context.Request().URI())
	}
	c := listClient.Get(clientName)
	return c.Init(this.Context, this.DB, m)
}

func (this *Model) NotModified() error {
	return ErrNotModified
}

func (this *Model) NewEditClient(m interface{}) editClient.Client {
	clientName := this.Context.Form(`client`)
	if len(clientName) == 0 {
		this.Logger.Debug(`Parameter client is empty in URL: ` + this.Context.Request().URI())
	}
	c := editClient.Get(clientName)
	return c.Init(this.Context, this.DB, m)
}

type Conversion struct {
	Original interface{}
}

func (c *Conversion) EditClient() editClient.Client {
	return c.Original.(editClient.Client)
}

func (c *Conversion) ListClient() listClient.Client {
	return c.Original.(listClient.Client)
}
