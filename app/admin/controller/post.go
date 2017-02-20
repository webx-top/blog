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
	"time"

	D "github.com/webx-top/blog/app/base/dbschema"
	"github.com/webx-top/blog/app/base/model"
	"github.com/webx-top/echo"
	X "github.com/webx-top/webx"
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

func (a *Post) Init(c echo.Context) error {
	a.Base = New(c)
	a.postM = model.NewPost(a.Context)
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

func (a *Post) validOk(m *D.Post) bool {
	valid := a.Valid()
	if r := valid.Required(m.Title, `Title`); !r.Ok {
		return false
	}
	if r := valid.Required(m.Content, `Content`); !r.Ok {
		return false
	}
	//valid.Required(m.Description, `Description`)
	if r := valid.Required(m.Catid, `Catid`); !r.Ok {
		return false
	}
	return true
}

func (a *Post) Add() error {
	m := &D.Post{}
	other := &D.Ocontent{}
	if a.IsPost() {
		err := a.Bind(m)
		if err != nil {
			return err
		}

		if a.validOk(m) {
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
				return a.GotoNext(`Index`)
			}
		}
	}
	a.Assign(`Detail`, m)
	a.Assign(`Other`, other)
	return a.Display(a.TmplPath(`Edit`))
}

func (a *Post) Edit() error {
	id := a.Formx(`id`).Int()
	m, has, err := a.postM.Get(id)
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
			affected, err := a.postM.Edit(m.Id, m)
			if err != nil {
				a.SetErr(err.Error())
			} else if affected < 1 {
				a.NotModified()
			} else {
				a.Done()
				return a.GotoNext(`Index`)
			}
		}
	}
	other, _, err := a.postM.GetOtherContent(m.Id)
	if err != nil {
		return err
	}
	a.Assign(`Detail`, m)
	a.Assign(`Other`, other)
	cateM := model.NewCategory(a.Context)
	a.Assign(`Breadcrumbs`, cateM.Dir(m.Catid))
	return a.Display()
}

func (a *Post) Delete() error {
	id := a.Formx(`id`).Int()
	if id < 1 {
		return a.NotFoundData().GotoNext(`Index`)
	}
	affected, err := a.postM.Delete(id)
	if err != nil {
		return err
	}
	if affected < 1 {
		return a.NotFoundData().GotoNext(`Index`)
	}
	a.Done()
	return a.GotoNext(`Index`)
}

func (a *Post) View() error {
	return a.Display()
}
