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
	. "github.com/webx-top/webx/lib/model"
)

func NewAlbum(ctx *X.Context) *Album {
	return &Album{M: NewM(ctx)}
}

type Album struct {
	*M
}

func (a *Album) List(s *Select) (countFn func() int64, m []*D.Album, err error) {
	m = []*D.Album{}
	err = s.Do().Find(&m)
	if err != nil {
		return
	}
	countFn = func() int64 {
		return s.Count(D.Album{})
	}
	return
}

func (a *Album) Add(m *D.Album) (affected int64, err error) {
	affected, err = a.Sess().Insert(m)
	return
}

func (a *Album) Edit(id int, m *D.Album) (affected int64, err error) {
	affected, err = a.Sess().Id(id).Update(m)
	return
}

func (a *Album) Del(id int) (affected int64, err error) {
	m := &D.Album{}
	affected, err = a.Sess().Where(`id=?`, id).Delete(m)
	return
}

func (a *Album) Get(id int) (m *D.Album, has bool, err error) {
	m = &D.Album{}
	has, err = a.DB.Id(id).Get(m)
	return
}
