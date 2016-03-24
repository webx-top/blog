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
	"encoding/gob"

	//"github.com/webx-top/blog/app/admin/lib"
	"github.com/webx-top/blog/app/base"
	"github.com/webx-top/blog/app/base/dbschema"
	X "github.com/webx-top/webx"
)

func init() {
	gob.Register(&dbschema.User{})
}

func New(c *X.Context) *Base {
	a := &Base{}
	a.Init(c)
	return a
}

type Base struct {
	*base.Controller
	*dbschema.User
}

func (a *Base) Init(c *X.Context) error {
	a.Controller = base.NewController(c)
	a.User = &dbschema.User{}
	return nil
}

func (a *Base) Before() error {
	ss := a.Session()
	if user, ok := ss.Get(`user`).(*dbschema.User); !ok || user == nil || user.Id < 1 {
		var errMsg = a.T(`请先登录`)
		a.SetNoAuth(errMsg)
		return a.Redirect(a.Url(`Public`, `Login`))
	} else {
		a.User = user
		user.Passwd = `[HIDE]`
		user.Salt = `[HIDE]`
		a.Assign(`User`, user)
	}
	return nil
}
