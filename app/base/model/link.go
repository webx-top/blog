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

func NewLink(ctx *X.Context) *Link {
	return &Link{M: NewM(ctx)}
}

type Link struct {
	*M
}

func (a *Link) List(s *Select) (countFn func() int64, m []*D.Link, err error) {
	m = []*D.Link{}
	err = s.Do().Find(&m)
	if err != nil {
		return
	}
	countFn = func() int64 {
		return s.Count(D.Link{})
	}
	return
}

func (a *Link) Add(m *D.Link) (affected int64, err error) {
	affected, err = a.Sess().Insert(m)
	return
}

func (a *Link) Edit(id int, m *D.Link) (affected int64, err error) {
	affected, err = a.Sess().Id(id).Update(m)
	return
}

func (a *Link) Delete(id int) (affected int64, err error) {
	m := &D.Link{}
	affected, err = a.Sess().Where(`id=?`, id).Delete(m)
	return
}

func (a *Link) Get(id int) (m *D.Link, has bool, err error) {
	m = &D.Link{}
	has, err = a.DB.Id(id).Get(m)
	return
}
