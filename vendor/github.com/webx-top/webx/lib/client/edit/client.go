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
package edit

import (
	"errors"

	X "github.com/webx-top/webx"
	//"github.com/webx-top/com"
	"github.com/webx-top/webx/lib/database"
)

type defaultClient struct {
	*X.Context
	*database.Orm
	model       interface{}
	structField string
	tableField  string
	value       string
	primary     string
}

func (a *defaultClient) Init(ctx *X.Context, orm *database.Orm, m interface{}, args ...string) Client {
	a.Context = ctx
	a.Orm = orm
	a.model = m
	a.structField = ctx.Form(`field`)
	a.value = ctx.Form(`value`)
	a.primary = ctx.Form(`id`)
	return a
}

func (a *defaultClient) Do(fn func(string, string, string) error, validField ...bool) error {
	if len(a.structField) < 1 {
		return errors.New(`Invalid field name: missing paramter field`)
	}
	if len(a.primary) == 0 {
		return errors.New(`Primary key value is invalid: missing paramter id`)
	}
	if len(validField) < 1 || validField[0] {
		a.tableField = a.Orm.ToTableField(a.model, a.structField)
		if len(a.tableField) < 1 {
			return errors.New(`Invalid field name: ` + a.structField)
		}
	} else {
		a.tableField = a.structField
	}
	a.Context.MapData(a.model, map[string][]string{
		a.structField: []string{a.value},
	})
	return fn(a.primary, a.tableField, a.value)
}

type Client interface {
	//初始化
	Init(*X.Context, *database.Orm, interface{}, ...string) Client
	Do(func(string, string, string) error, ...bool) error
}
