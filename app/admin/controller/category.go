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
	D "github.com/webx-top/blog/app/base/dbschema"
	"github.com/webx-top/blog/app/base/model"
	X "github.com/webx-top/webx"
	"github.com/webx-top/webx/lib/com"
)

type Category struct {
	index  X.Mapper
	add    X.Mapper
	edit   X.Mapper
	delete X.Mapper
	*Base
	cateM *model.Category
}

func (a *Category) Init(c *X.Context) error {
	a.Base = New(c)
	a.cateM = model.NewCategory(c)
	return nil
}

func (a *Category) Index() error {
	pid := com.Int(a.Query(`pid`))
	nid := com.Int(a.Query(`ignore`))
	sel := a.cateM.NewSelect(&D.Category{})
	sel.Condition = `pid=?`
	sel.AddParam(pid)
	if nid > 0 {
		sel.Condition += ` AND id!=?`
		sel.AddParam(nid)
	}
	sel.FromClient(true, "Name")
	countFn, data, _ := a.cateM.List(sel)
	sel.Client.SetCount(countFn).Data(data)
	a.Assign(`Breadcrumbs`, a.cateM.Dir(pid))
	return a.Display()
}

func (a *Category) Index_HTML() error {
	pid := com.Int(a.Query(`pid`))
	a.Assign(`Breadcrumbs`, a.cateM.Dir(pid))
	return a.Display()
}

func (a *Category) validOk(m *D.Category) bool {
	valid := a.Valid()
	if r := valid.Required(m.Name, `Name`); !r.Ok {
		return r.Ok
	}
	return true
}

func (a *Category) Add() error {
	m := &D.Category{}
	if a.IsPost() {
		err := a.Bind(m)
		if err != nil {
			return err
		}

		if a.validOk(m) {
			affected, err := a.cateM.Add(m)
			if err != nil {
				a.SetErr(err.Error())
			} else if affected < 1 {
				a.NotModified()
			} else {
				a.Done()
				return a.GotoNext("Category", "Index")
			}
		}
	}
	a.Assign(`Detail`, m)
	return a.Display(a.TmplPath(`Edit`))
}

func (a *Category) Edit() error {
	id := com.Int(a.Form(`id`))
	m, has, err := a.cateM.Get(id)
	if err != nil {
		return err
	} else if !has {
		return a.NotFoundData().Display()
	}
	if a.IsPost() {
		err = a.Bind(m)
		if err != nil {
			return err
		}
		if a.validOk(m) {
			affected, err := a.cateM.Edit(m.Id, m)
			if err != nil {
				a.SetErr(err.Error())
			} else if affected < 1 {
				a.NotModified()
			} else {
				a.Done()
				return a.GotoNext("Category", "Index")
			}
		}
	}
	a.Assign(`Detail`, m)
	a.Assign(`Breadcrumbs`, a.cateM.Dir(m.Pid))
	return a.Display()
}

func (a *Category) Delete() error {
	id := com.Int(a.Form(`id`))
	if id < 1 {
		return a.NotFoundData().GotoNext(`Index`)
	}
	affected, err := a.cateM.Delete(id)
	if err != nil {
		return err
	}
	if affected < 1 {
		return a.NotFoundData().GotoNext(`Index`)
	}
	a.Done()
	return a.GotoNext(`Index`)
}
