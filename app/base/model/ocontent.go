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
)

func NewOcontent(ctx *X.Context) *Ocontent {
	return &Ocontent{M: NewM(ctx)}
}

type Ocontent struct {
	*M
}

func (a *Ocontent) Add(m *D.Ocontent) (affected int64, err error) {
	affected, err = a.Sess().Insert(m)
	return
}

func (a *Ocontent) Edit(id int64, m *D.Ocontent) (affected int64, err error) {
	affected, err = a.Sess().Id(id).Update(m)
	return
}

func (a *Ocontent) GetByMaster(rid int64, rtype string) (m *D.Ocontent, has bool, err error) {
	m = &D.Ocontent{}
	has, err = a.DB.Where(`rc_id=? AND rc_type=?`, rid, rtype).Get(m)
	return
}

func (a *Ocontent) DelByMaster(rid int64, rtype string) (affected int64, err error) {
	m := &D.Ocontent{}
	affected, err = a.Sess().Where(`rc_id=? AND rc_type=?`, rid, rtype).Delete(m)
	return
}

func (a *Ocontent) Get(id int64) (m *D.Ocontent, has bool, err error) {
	m = &D.Ocontent{}
	has, err = a.DB.Id(id).Get(m)
	return
}
