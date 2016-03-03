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
	//"fmt"

	D "github.com/webx-top/blog/app/base/dbschema"
	X "github.com/webx-top/webx"
	. "github.com/webx-top/webx/lib/model"
)

func NewCategory(ctx *X.Context) *Category {
	return &Category{M: NewM(ctx)}
}

type Category struct {
	*M
}

func (a *Category) List(s *Select) (countFn func() int64, m []*D.Category, err error) {
	m = []*D.Category{}
	err = s.Do().Find(&m)
	if err != nil {
		return
	}
	countFn = func() int64 {
		return s.Count(D.Category{})
	}
	return
}

func (a *Category) Add(m *D.Category) (affected int64, err error) {
	affected, err = a.Sess().Insert(m)
	return
}

func (a *Category) Edit(id int, m *D.Category) (affected int64, err error) {
	affected, err = a.Sess().Id(id).Update(m)
	return
}

func (a *Category) Delete(id int) (affected int64, err error) {
	m := &D.Category{}
	affected, err = a.Sess().Where(`id=?`, id).Delete(m)
	return
}

func (a *Category) Get(id int) (m *D.Category, has bool, err error) {
	m = &D.Category{}
	has, err = a.DB.Id(id).Get(m)
	return
}

func (a *Category) Dir(pid int) (rs []D.Category) {
	m := D.Category{}
	var has bool
	has, _ = a.DB.Id(pid).Get(&m)
	r := make([]D.Category, 0)
	for has {
		r = append(r, m)
		if m.Pid <= 0 {
			break
		}
		pid := m.Pid
		m = D.Category{}
		has, _ = a.DB.Id(pid).Get(&m)
	}
	rs = make([]D.Category, 0)
	for i := len(r) - 1; i >= 0; i-- {
		rs = append(rs, r[i])
	}
	return
}
