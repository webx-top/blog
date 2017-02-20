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
	"github.com/webx-top/echo"
	X "github.com/webx-top/webx"
)

type Comment struct {
	index  X.Mapper
	add    X.Mapper
	edit   X.Mapper
	delete X.Mapper
	view   X.Mapper
	*Base
	cmtM *model.Comment
}

func (a *Comment) Init(c echo.Context) error {
	a.Base = New(c)
	a.cmtM = model.NewComment(a.Context)
	return nil
}

func (a *Comment) Index_HTML() error {
	return a.Display()
}

func (a *Comment) Index() error {
	sel := a.cmtM.NewSelect(&model.CommentWithPost{})
	sel.AddParam(a.User.Id).FromClient(true, "Comment.content")
	countFn, data, _ := a.cmtM.ListWithPost(sel)
	sel.Client.SetCount(countFn).Data(data)
	return a.Display()
}

func (a *Comment) Add() error {
	m := &D.Comment{}
	other := &D.Ocontent{}
	if a.IsPost() {
		err := a.Bind(m)
		if err != nil {
			return err
		}

		if a.ValidOk(m) {
			affected, err := a.cmtM.Add(m)
			if err != nil {
				a.SetErr(err.Error())
			} else if affected < 1 {
				a.NotModified()
			} else {
				a.Done()
				a.GotoNext(`Index`)
			}
		}
	}
	a.Assign(`Detail`, m)
	a.Assign(`Other`, other)
	return a.Display(a.TmplPath(`Edit`))
}

func (a *Comment) Edit() error {
	id := a.Formx(`id`).Int64()
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
			return a.GotoNext(`Index`)
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
	id := a.Formx(`id`).Int64()
	if id < 1 {
		return a.NotFoundData().GotoNext(`Index`)
	}
	affected, err := a.cmtM.Delete(id)
	if err != nil {
		return err
	}
	if affected < 1 {
		return a.NotFoundData().GotoNext(`Index`)
	}
	a.Done()
	return a.GotoNext(`Index`)
}

func (a *Comment) View() error {
	return a.Display()
}
