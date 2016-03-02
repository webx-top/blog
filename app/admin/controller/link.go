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

type Link struct {
	index  X.Mapper
	add    X.Mapper
	edit   X.Mapper
	delete X.Mapper
	view   X.Mapper
	*Base
	lnkM *model.Link
}

func (a *Link) Init(c *X.Context) error {
	a.Base = New(c)
	a.lnkM = model.NewLink(c)
	return nil
}

func (a *Link) Index() error {
	if a.Format != `html` {
		sel := a.lnkM.NewSelect(&D.Link{})
		sel.Condition = `uid=?`
		sel.AddParam(a.User.Id).FromClient(true, "title")
		countFn, data, _ := a.lnkM.List(sel)
		sel.Client.SetCount(countFn).Data(data)
	}
	return a.Display()
}

func (a *Link) Add() error {
	m := &D.Link{}
	errs := make(map[string]string)
	if a.IsPost() {
		err := a.Bind(m)
		if err != nil {
			return err
		}

		if ok, es, _ := a.Valid(m); !ok {
			errs = es
		} else {
			affected, err := a.lnkM.Add(m)
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

func (a *Link) Edit() error {
	id := com.Int(a.Form(`id`))
	m, has, err := a.lnkM.Get(id)
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
		affected, err := a.lnkM.Edit(m.Id, m)
		if err != nil {
			a.SetErr(err.Error())
		} else if affected < 1 {
			a.NotModified()
		} else {
			a.Done()
		}
	}
	a.Assign(`Detail`, m)
	return a.Display()
}

func (a *Link) Delete() error {
	id := com.Int(a.Form(`id`))
	if id < 1 {
		return a.NotFoundData().Display()
	}
	affected, err := a.lnkM.Delete(id)
	if err != nil {
		return err
	}
	if affected < 1 {
		return a.NotFoundData().Display()
	}
	a.Done()
	return a.Display()
}

func (a *Link) View() error {
	return a.Display()
}
