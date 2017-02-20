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
	return &Comment{M: NewM(ctx)}
}

type Comment struct {
	*M
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
	affected, err = a.Sess().Insert(m)
	return
}

func (a *Comment) Edit(id int64, m *D.Comment) (affected int64, err error) {
	affected, err = a.Sess().Id(id).Update(m)
	return
}

func (a *Comment) Delete(id int64) (affected int64, err error) {
	m := &D.Comment{}
	affected, err = a.Sess().Where(`id=?`, id).Delete(m)
	return
}

func (a *Comment) Get(id int64) (m *D.Comment, has bool, err error) {
	m = &D.Comment{}
	has, err = a.DB.Id(id).Get(m)
	return
}
