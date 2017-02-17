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
	"fmt"

	"github.com/coscms/xorm"
	listClient "github.com/webx-top/webx/lib/client/list"
	"github.com/webx-top/webx/lib/database"
)

func NewSelect(orm *database.Orm, c listClient.Client) *Select {
	s := &Select{
		Orm:    orm,
		Params: make([]interface{}, 0),
		Client: c,
	}
	return s
}

type Select struct {
	//偏移值
	Offset int64

	//查询数据量
	Limit int64

	//排序
	OrderBy string

	//查询条件
	Condition string

	//查询条件中的参数值
	Params []interface{}

	//group by 子句
	GroupBy string

	//having 子句
	Having string

	//表对象或名称
	Table interface{}

	//表别名
	Alias string

	*database.Orm
	listClient.Client
	Error error
}

// 生成查询条件。 参数 1-表名, 2-别名
func (a *Select) Do(args ...interface{}) *xorm.Session {
	return a.GenSess(args...).OrderBy(a.OrderBy).Limit(int(a.Limit), int(a.Offset))
}

// 添加查询条件参数值
func (a *Select) AddParam(args ...interface{}) *Select {
	a.Params = append(a.Params, args...)
	return a
}

func (a *Select) SetPage(page int64, size int64) {
	if page < 1 {
		page = 1
	}
	if size > 1000 {
		size = 1000
	} else if size < 1 {
		size = 10
	}
	a.Offset = (page - 1) * size
	a.Limit = size
}

// 从客户端获取查询条件
func (a *Select) FromClient(gen bool, fields ...string) *Select {
	a.OrderBy = a.Client.GenOrderBy()
	a.Offset = a.Client.Offset()
	a.Limit = a.Client.PageSize()
	if !gen {
		return a
	}
	a.GenSearchCond(fields...)
	return a
}

// 生成查询条件
func (a *Select) GenSearchCond(fields ...string) *Select {
	sch := a.Client.GenSearch(fields...)
	if len(sch) > 0 {
		if len(a.Condition) > 0 {
			a.Condition += ` AND `
		}
		a.Condition += sch
	}
	return a
}

// 从设置客户端数据
func (a *Select) SetProcesser(processer func(*Select) (func() int64, interface{}, error)) *Select {
	countFn, data, err := processer(a)
	a.Error = err
	return a.SetClient(countFn, data)
}

// 从设置客户端数据
func (a *Select) SetClient(countFn func() int64, data interface{}) *Select {
	a.Client.SetCount(countFn).Data(data)
	return a
}

// 生成查询条件。 参数 1-表名, 2-别名
func (a *Select) GenSess(args ...interface{}) *xorm.Session {
	s := a.Orm.Replica().NewSession()
	s.IsAutoClose = true
	switch len(args) {
	case 2:
		alias, _ := args[1].(string)
		if args[0] == nil {
			s = s.Alias(alias)
		} else {
			if tableName, ok := args[0].(string); ok {
				tableName = a.Orm.TableName(tableName)
				s = s.Table(tableName).Alias(alias)
				a.Table = tableName
			} else {
				s = s.Table(args[0]).Alias(alias)
				a.Table = args[0]
			}
		}
		a.Alias = alias
	case 1:
		if tableName, ok := args[0].(string); ok {
			tableName = a.Orm.TableName(tableName)
			s = s.Table(tableName)
			a.Table = tableName
		} else {
			s = s.Table(args[0])
			a.Table = args[0]
		}
	default:
		if a.Table != nil {
			s = s.Table(a.Table)
		}
		if a.Alias != `` {
			s = s.Alias(a.Alias)
		}
	}
	s = s.Where(a.Condition, a.Params...).GroupBy(a.GroupBy)
	if a.Having != `` {
		s = s.Having(a.Having)
	}
	return s
}

// 统计数据量
func (a *Select) Count(m interface{}, args ...interface{}) int64 {
	count, err := a.GenSess(args...).Count(m)
	if err != nil {
		fmt.Println(err)
	}
	return count
}

func (s *Select) ForList(bean interface{}, listFn func(*Select) (func() int64, interface{}, error), searchFields string, condition string, args ...interface{}) {
	if len(condition) > 0 {
		s.Condition = condition
		if len(args) > 0 {
			s.AddParam(args...)
		}
	}
	if len(searchFields) > 0 {
		s.FromClient(true, searchFields)
	} else {
		s.FromClient(true)
	}
	countFn, data, _ := listFn(s)
	s.Client.SetCount(countFn).Data(data)
}
