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
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	X "github.com/webx-top/webx"
)

type Setting struct {
	index  X.Mapper
	add    X.Mapper
	edit   X.Mapper
	delete X.Mapper
	view   X.Mapper
	*Base
	confM *model.Config
}

func (a *Setting) Init(c echo.Context) error {
	a.Base = New(c)
	a.confM = model.NewConfig(a.Context)
	return nil
}

func (a *Setting) Index_HTML() error {
	return a.Display()
}

func (a *Setting) Index() error {
	sel := a.confM.NewSelect(&D.Config{})
	sel.Condition = `uid=?`
	sel.AddParam(a.User.Id).FromClient(true, "title")
	countFn, data, _ := a.confM.List(sel)
	sel.Client.SetCount(countFn).Data(data)
	return a.Display()
}

func (a *Setting) Add() error {
	m := &D.Config{}
	if a.IsPost() {
		err := a.Bind(m)
		if err != nil {
			return err
		}

		if a.ValidOk(m) {
			affected, err := a.confM.Add(m)
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
	return a.Display(a.TmplPath(`Edit`))
}

func (a *Setting) Edit() error {
	id := com.Int(a.Form(`id`))
	m, has, err := a.confM.Get(id)
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
		affected, err := a.confM.Edit(m.Id, m)
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

func (a *Setting) Delete() error {
	return a.Display()
}

func (a *Setting) View() error {
	return a.Display()
}
