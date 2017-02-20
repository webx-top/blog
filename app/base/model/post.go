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
	"strings"

	D "github.com/webx-top/blog/app/base/dbschema"
	X "github.com/webx-top/webx"
	. "github.com/webx-top/webx/lib/model"
)

func NewPost(ctx *X.Context) *Post {
	return &Post{M: NewM(ctx), C: NewOcontent(ctx), typeName: `post`}
}

type Post struct {
	*M
	C        *Ocontent
	Uid      int
	typeName string
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
	tag := NewTag(a.M.Model.Context)
	if m.Tags != `` {
		tags := strings.SplitN(m.Tags, `,`, 5)
		var data []*D.Tag
		data, err = tag.AddNotExists(a.Uid, a.typeName, tags...)
		if err != nil {
			return
		}
		for _, v := range data {
			_, err = tag.UpdateTimes(v.Id, 1)
			if err != nil {
				return
			}
		}
	}
	otherContent := a.EditorContent(m)
	a.Trans(func() error {
		affected, err = a.Sess().Insert(m)
		if err != nil {
			return err
		}

		if m.Etype != `html` {
			occ := &D.Ocontent{
				RcId:    int64(m.Id),
				RcType:  `post`,
				Content: otherContent,
				Etype:   `markdown`,
			}
			_, err = a.C.Add(occ)
		}
		return err
	})
	return
}

func (a *Post) Delete(id int) (affected int64, err error) {
	m, has, err := a.Get(id)
	if err != nil {
		return
	}
	if !has {
		err = a.Context.Atoe(a.T(`数据不存在`))
		return
	}
	if m.Tags != `` {
		tag := NewTag(a.M.Model.Context)
		tags := strings.SplitN(m.Tags, `,`, -1)
		err = tag.DelExists(a.typeName, tags...)
		if err != nil {
			return
		}
	}

	a.Trans(func() error {
		affected, err = a.Sess().Id(id).Delete(&D.Post{})
		if err != nil {
			return err
		}
		_, err = a.C.DelByMaster(int64(id), `post`)
		return err
	})
	return
}

func (a *Post) Edit(id int, m *D.Post) (affected int64, err error) {
	oc, has, err := a.C.GetByMaster(int64(id), `post`)
	otherContent := a.EditorContent(m)
	old, has2, err := a.Get(id)
	if err != nil {
		return
	}
	if !has2 {
		err = a.Context.Atoe(a.T(`数据不存在`))
		return
	}
	a.Trans(func() error {
		affected, err = a.Sess().Id(id).Update(m)
		if err != nil {
			return err
		}
		if m.Etype != `html` {
			occ := &D.Ocontent{
				RcId:    int64(id),
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
				_, err = a.C.DelByMaster(int64(id), `post`)
			}
		}
		if err != nil {
			return err
		}
		tag := NewTag(a.M.Model.Context)
		tags := strings.Split(m.Tags, `,`)
		if old.Tags != `` {
			oldTags := strings.Split(old.Tags, `,`)
			var delTags []string
			for _, t := range oldTags {
				exists := false
				for _, v := range tags {
					if v == t {
						exists = true
						break
					}
				}
				if !exists {
					delTags = append(delTags, t)
				}
			}
			if len(delTags) > 0 {
				tag.DelExists(a.typeName, delTags...)
			}
		}
		_, err = tag.AddNotExists(a.Uid, a.typeName, tags...)
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
	return a.C.GetByMaster(int64(id), `post`)
}

func (a *Post) EditorContent(m *D.Post) (otherContent string) {
	m.Content, otherContent = EditorContent(a.Context, m.Etype, m.Content)
	return
}

func EditorContent(ctx *X.Context, etype string, editorContent string) (string, string) {
	var rawContent string
	switch etype {
	case `markdown`:
		rawContent = editorContent
		editorId := ctx.Form(`EditorId`)
		if len(editorId) > 0 {
			editorContent = ctx.Form(editorId + `-html-code`)
		} else {
			editorContent = ``
		}
	}
	return editorContent, rawContent
}
