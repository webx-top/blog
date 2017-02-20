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
	D "github.com/webx-top/blog/app/base/dbschema"
	X "github.com/webx-top/webx"
	. "github.com/webx-top/webx/lib/model"
)

func NewComment(ctx *X.Context) *Comment {
	return &Comment{M: NewM(ctx), C: NewOcontent(ctx)}
}

type Comment struct {
	*M
	C *Ocontent
}

type CommentWithPost struct {
	Comment *D.Comment `xorm:"extends"`
	Post    *D.Post    `xorm:"extends" rel:"LEFT:Post.id=Comment.rc_id"`
}

func (a *Comment) List(s *Select) (countFn func() int64, m []*D.Comment, err error) {
	m = []*D.Comment{}
	err = s.Do().Find(&m)
	if err != nil {
		return
	}
	countFn = func() int64 {
		return s.Count(D.Comment{})
	}
	return
}

func (a *Comment) ListWithPost(s *Select) (countFn func() int64, m []*CommentWithPost, err error) {
	m = []*CommentWithPost{}
	err = s.Do().Find(&m)
	if err != nil {
		return
	}
	countFn = func() int64 {
		return s.Count(&CommentWithPost{})
	}
	return
}

func (a *Comment) Add(m *D.Comment) (affected int64, err error) {
	otherContent := a.EditorContent(m)
	a.Trans(func() error {
		affected, err = a.Sess().Insert(m)
		if err != nil {
			return err
		}

		if m.Etype != `html` {
			occ := &D.Ocontent{
				RcId:    int64(m.Id),
				RcType:  `comment`,
				Content: otherContent,
				Etype:   `markdown`,
			}
			_, err = a.C.Add(occ)
		}
		return err
	})
	return
}

func (a *Comment) EditorContent(m *D.Comment) (otherContent string) {
	m.Content, otherContent = EditorContent(a.Context, m.Etype, m.Content)
	return
}

func (a *Comment) Edit(id int64, m *D.Comment) (affected int64, err error) {
	oc, has, err := a.C.GetByMaster(int64(id), `post`)
	otherContent := a.EditorContent(m)
	a.Trans(func() error {
		affected, err = a.Sess().Id(id).Update(m)
		if err != nil {
			return err
		}
		if m.Etype != `html` {
			occ := &D.Ocontent{
				RcId:    id,
				RcType:  `comment`,
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
				_, err = a.C.DelByMaster(id, `comment`)
			}
		}
		return err
	})
	return
}

func (a *Comment) Delete(id int64) (affected int64, err error) {
	m := &D.Comment{}
	a.Trans(func() error {
		affected, err = a.Sess().Where(`id=?`, id).Delete(m)
		if err != nil {
			return err
		}
		_, err = a.C.DelByMaster(id, `comment`)
		return err
	})
	return
}

func (a *Comment) Get(id int64) (m *D.Comment, has bool, err error) {
	m = &D.Comment{}
	has, err = a.DB.Id(id).Get(m)
	return
}

func (a *Comment) GetOtherContent(id int64) (m *D.Ocontent, has bool, err error) {
	return a.C.GetByMaster(id, `comment`)
}
