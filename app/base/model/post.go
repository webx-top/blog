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
package model

import (
	//"errors"
	//"strings"

	D "github.com/webx-top/blog/app/base/dbschema"
	X "github.com/webx-top/webx"
	//"github.com/webx-top/webx/lib/com"
)

func NewPost(ctx *X.Context) *Post {
	return &Post{M: NewM(ctx), C: NewOcontent(ctx)}
}

type Post struct {
	*M
	C *Ocontent
}

func (a *Post) List(s *Select) (countFn func() int64, m []*D.Post, err error) {
	m = []*D.Post{}
	err = s.Do().Find(&m)
	if err != nil {
		return
	}
	countFn = func() int64 {
		return s.Count(D.Post{})
	}
	return
}

func (a *Post) Add(m *D.Post) (affected int64, err error) {
	affected, err = a.DB.Insert(m)
	return
}

func (a *Post) Edit(id int, m *D.Post) (affected int64, err error) {
	oc, has, err := a.C.GetByMaster(id, `post`)
	otherContent := a.EditorContent(m)
	a.Trans(func() error {
		affected, err = a.Sess().Id(id).Update(m)
		if err != nil {
			return err
		}
		if m.Etype != `html` {
			occ := &D.Ocontent{
				RcId:    id,
				RcType:  `post`,
				Content: otherContent,
				Etype:   `markdown`,
			}
			if has {
				_, err = a.C.Edit(oc.Id, occ)
			} else {
				_, err = a.C.Add(occ)
			}
		} else {
			if has {
				_, err = a.C.DelByMaster(id, `post`)
			}
		}
		return err
	})
	return
}

func (a *Post) Get(id int) (m *D.Post, has bool, err error) {
	m = &D.Post{}
	has, err = a.DB.Id(id).Get(m)
	return
}

func (a *Post) GetOtherContent(id int) (m *D.Ocontent, has bool, err error) {
	return a.C.GetByMaster(id, `post`)
}

func (a *Post) EditorContent(m *D.Post) (otherContent string) {
	switch m.Etype {
	case `markdown`:
		otherContent = m.Content
		editorId := a.M.Context.Form(`editorId`)
		if editorId != `` {
			m.Content = a.M.Context.Form(editorId + `-html-code`)
		} else {
			m.Content = ``
		}
	}
	return
}
