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
	"time"

	"github.com/webx-top/blog/app/admin/lib"
	D "github.com/webx-top/blog/app/base/dbschema"
	"github.com/webx-top/blog/app/base/model"
	X "github.com/webx-top/webx"
	"github.com/webx-top/webx/lib/com"
)

func init() {
	lib.App.Reg(&Category{}).Auto()
}

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
	if a.Format != `html` {
		sel := a.cateM.NewSelect(&D.Category{})
		sel.Condition = `uid=?`
		sel.AddP(a.User.Id).FromClient(true, "title")
		countFn, data, _ := a.cateM.List(sel)
		sel.Client.SetCount(countFn).Data(data)
	}
	return a.Display()
}

func (a *Category) Add() error {
	m := &D.Category{}
	errs := make(map[string]string)
	if a.IsPost() {
		err := a.Bind(m)
		if err != nil {
			return err
		}

		if ok, es, _ := a.Valid(m); !ok {
			errs = es
		} else {
			m.Uid = a.User.Id
			m.Uname = a.User.Uname
			t := time.Now().Local()
			m.Year = t.Year()
			m.Month = com.Int(t.Month().String())
			affected, err := a.cateM.Add(m)
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
		affected, err := a.cateM.Edit(m.Id, m)
		if err != nil {
			a.SetErr(err.Error())
		} else if affected < 1 {
			a.NotModified()
		} else {
			a.Done()
		}
	}
	other, _, err := a.cateM.GetOtherContent(m.Id)
	if err != nil {
		return err
	}
	a.Assign(`Detail`, m)
	a.Assign(`Other`, other)
	return a.Display()
}

func (a *Category) Delete() error {
	return a.Display()
}

func (a *Category) View() error {
	return a.Display()
}
