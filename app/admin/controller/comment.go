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
	lib.App.Reg(&Comment{}).Auto()
}

type Comment struct {
	index  X.Mapper
	add    X.Mapper
	edit   X.Mapper
	delete X.Mapper
	view   X.Mapper
	*Base
	cmtM *model.Comment
}

func (a *Comment) Init(c *X.Context) error {
	a.Base = New(c)
	a.cmtM = model.NewComment(c)
	return nil
}

func (a *Comment) Index() error {
	if a.Format != `html` {
		sel := a.cmtM.NewSelect(&D.Comment{})
		sel.Condition = `uid=?`
		sel.AddP(a.User.Id).FromClient(true, "title")
		countFn, data, _ := a.cmtM.List(sel)
		sel.Client.SetCount(countFn).Data(data)
	}
	return a.Display()
}

func (a *Comment) Add() error {
	m := &D.Comment{}
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
			affected, err := a.cmtM.Add(m)
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

func (a *Comment) Edit() error {
	id := com.Int(a.Form(`id`))
	m, has, err := a.cmtM.Get(id)
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
		affected, err := a.cmtM.Edit(m.Id, m)
		if err != nil {
			a.SetErr(err.Error())
		} else if affected < 1 {
			a.NotModified()
		} else {
			a.Done()
		}
	}
	other, _, err := a.cmtM.GetOtherContent(m.Id)
	if err != nil {
		return err
	}
	a.Assign(`Detail`, m)
	a.Assign(`Other`, other)
	return a.Display()
}

func (a *Comment) Delete() error {
	return a.Display()
}

func (a *Comment) View() error {
	return a.Display()
}
