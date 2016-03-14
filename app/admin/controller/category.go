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
	//"strings"

	//"github.com/webx-top/blog/app/admin/lib"
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
	view   X.Mapper
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

func (a *Category) validate(m *D.Category) (bool, map[string]string) {
	ok, es, valid := a.Valid(nil)
	valid.Required(m.Name, `Name`)
	ok = valid.HasErrors() == false
	es = valid.ErrMap()
	for key, msg := range es {
		a.SetErr(msg, key)
		break
	}
	return ok, es
}

func (a *Category) Add() error {
	m := &D.Category{}
	errs := make(map[string]string)
	if a.IsPost() {
		err := a.Bind(m)
		if err != nil {
			return err
		}

		if ok, es := a.validate(m); !ok {
			errs = es
		} else {
			affected, err := a.cateM.Add(m)
			if err != nil {
				a.SetErr(err.Error())
			} else if affected < 1 {
				a.NotModified()
			} else {
				a.Done()
				return a.Redirect(a.Url("Category", "Index"))
			}
		}
	}
	a.Assign(`Detail`, m)
	a.Assign(`Errors`, errs)
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
	errs := make(map[string]string)
	if a.IsPost() {
		err = a.Bind(m)
		if err != nil {
			return err
		}
		if ok, es := a.validate(m); !ok {
			errs = es
		} else {
			affected, err := a.cateM.Edit(m.Id, m)
			if err != nil {
				a.SetErr(err.Error())
			} else if affected < 1 {
				a.NotModified()
			} else {
				a.Done()
			}
		}
	}
	a.Assign(`Detail`, m)
	a.Assign(`Errors`, errs)
	a.Assign(`Breadcrumbs`, a.cateM.Dir(m.Pid))
	return a.Display()
}

func (a *Category) Delete() error {
	id := com.Int(a.Form(`id`))
	if id < 1 {
		return a.NotFoundData().Display()
	}
	affected, err := a.cateM.Delete(id)
	if err != nil {
		return err
	}
	if affected < 1 {
		return a.NotFoundData().Display()
	}
	a.Done()
	return a.Display()
}

func (a *Category) View() error {
	return a.Display()
}
