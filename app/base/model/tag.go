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
	"math"

	D "github.com/webx-top/blog/app/base/dbschema"
	X "github.com/webx-top/webx"
	. "github.com/webx-top/webx/lib/model"
)

func NewTag(ctx *X.Context) *Tag {
	return &Tag{M: NewM(ctx)}
}

type Tag struct {
	*M
}

func (a *Tag) List(s *Select) (countFn func() int64, m []*D.Tag, err error) {
	m = []*D.Tag{}
	err = s.Do().Find(&m)
	if err != nil {
		return
	}
	countFn = func() int64 {
		return s.Count(D.Tag{})
	}
	return
}

func (a *Tag) Add(m *D.Tag) (affected int64, err error) {
	affected, err = a.Sess().Insert(m)
	return
}

func (a *Tag) Edit(id int, m *D.Tag) (affected int64, err error) {
	affected, err = a.Sess().Id(id).Update(m)
	return
}

func (a *Tag) Delete(id int) (affected int64, err error) {
	m := &D.Tag{}
	affected, err = a.Sess().Where(`id=?`, id).Delete(m)
	return
}

func (a *Tag) UpdateTimes(id int, n int) (affected int64, err error) {
	m := &D.Tag{}
	if n > 0 {
		affected, err = a.Sess().Id(id).Incr(`times`, n).Update(m)
	} else {
		affected, err = a.Sess().Id(id).Decr(`times`, math.Abs(float64(id))).Update(m)
	}
	return
}

func (a *Tag) AddNotExists(uid int, rcType string, tags ...string) (m []*D.Tag, err error) {
	m = []*D.Tag{}
	params := make([]interface{}, len(tags))
	for k, v := range tags {
		params[k] = v
	}
	err = a.Sess().Where(`rc_type=?`, rcType).In(`name`, params...).Find(&m)
	rs := make([]string, 0)
	if err != nil {
		return
	}
	for _, tag := range tags {
		exists := false
		for _, v := range m {
			if v.Name == tag {
				exists = true
				break
			}
		}
		if !exists {
			rs = append(rs, tag)
		}
	}
	for _, tag := range rs {
		if tag == `` {
			continue
		}
		_m := &D.Tag{
			Name:   tag,
			Uid:    uid,
			RcType: rcType,
			Times:  1,
		}
		_, err = a.Add(_m)
		if err != nil {
			return
		}
	}
	return
}

func (a *Tag) DelExists(rcType string, tags ...string) (err error) {
	m := []*D.Tag{}
	params := make([]interface{}, len(tags))
	for k, v := range tags {
		params[k] = v
	}
	err = a.Sess().Where(`rc_type=?`, rcType).In(`name`, params...).Find(&m)
	if err != nil {
		return
	}
	for _, v := range m {
		if v.Times <= 1 {
			_, err = a.Delete(v.Id)
			if err != nil {
				return
			}
			continue
		}
		_, err = a.UpdateTimes(v.Id, -1)
		if err != nil {
			return
		}
	}
	return
}

func (a *Tag) Get(id int) (m *D.Tag, has bool, err error) {
	m = &D.Tag{}
	has, err = a.DB.Id(id).Get(m)
	return
}
