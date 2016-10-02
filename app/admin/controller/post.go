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

	//"github.com/webx-top/blog/app/admin/lib"
	D "github.com/webx-top/blog/app/base/dbschema"
	"github.com/webx-top/blog/app/base/model"
	X "github.com/webx-top/webx"
	"github.com/webx-top/webx/lib/com"
)

type Post struct {
	index  X.Mapper
	add    X.Mapper
	edit   X.Mapper
	delete X.Mapper
	view   X.Mapper
	*Base
	postM *model.Post
}

func (a *Post) Init(c *X.Context) error {
	a.Base = New(c)
	a.postM = model.NewPost(c)
	return nil
}

func (a *Post) Before() error {
	err := a.Base.Before()
	if err != nil {
		return err
	}
	a.postM.Uid = a.User.Id
	return nil
}

func (a *Post) Index() error {
	sel := a.postM.NewSelect(&D.Post{})
	sel.Condition = `uid=?`
	sel.AddParam(a.User.Id).FromClient(true, "title")
	countFn, data, _ := a.postM.List(sel)
	sel.Client.SetCount(countFn).Data(data)
	return a.Display()
}

func (a *Post) Index_HTML() error {
	return a.Display()
}

func (a *Post) validate(m *D.Post) (bool, map[string]string) {
	ok, es, valid := a.Valid(nil)
	valid.Required(m.Title, `Title`)
	valid.Required(m.Content, `Content`)
	//valid.Required(m.Description, `Description`)
	valid.Required(m.Catid, `Catid`)
	ok = valid.HasError() == false
	es = valid.ErrMap()
	for key, msg := range es {
		a.SetErr(msg, key)
		break
	}
	return ok, es
}

func (a *Post) Add() error {
	m := &D.Post{}
	other := &D.Ocontent{}
	errs := make(map[string]string)
	if a.IsPost() {
		err := a.Bind(m)
		if err != nil {
			return err
		}

		if ok, es := a.validate(m); !ok {
			errs = es
		} else {
			m.Uid = a.User.Id
			m.Uname = a.User.Uname
			t := time.Now().Local()
			m.Year = t.Year()
			m.Month = int(t.Month())
			affected, err := a.postM.Add(m)
			if err != nil {
				a.SetErr(err.Error())
			} else if affected < 1 {
				a.NotModified()
			} else {
				a.Done()
				return a.Redirect(a.BuildURL(`Post`, `Index`))
			}
		}
	}
	a.Assign(`Detail`, m)
	a.Assign(`Other`, other)
	a.Assign(`Errors`, errs)
	return a.Display(a.TmplPath(`Edit`))
}

func (a *Post) Edit() error {
	id := com.Int(a.Form(`id`))
	m, has, err := a.postM.Get(id)
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
			affected, err := a.postM.Edit(m.Id, m)
			if err != nil {
				a.SetErr(err.Error())
			} else if affected < 1 {
				a.NotModified()
			} else {
				a.Done()
			}
		}
	}
	other, _, err := a.postM.GetOtherContent(m.Id)
	if err != nil {
		return err
	}
	a.Assign(`Detail`, m)
	a.Assign(`Other`, other)
	a.Assign(`Errors`, errs)
	cateM := model.NewCategory(a.Context)
	a.Assign(`Breadcrumbs`, cateM.Dir(m.Catid))
	return a.Display()
}

func (a *Post) Delete() error {
	id := com.Int(a.Form(`id`))
	if id < 1 {
		return a.NotFoundData().Redir(a.NextURL(`Index`))
	}
	affected, err := a.postM.Delete(id)
	if err != nil {
		return err
	}
	if affected < 1 {
		return a.NotFoundData().Redir(a.NextURL(`Index`))
	}
	a.Done()
	return a.Redir(a.NextURL(`Index`))
}

func (a *Post) View() error {
	return a.Display()
}
