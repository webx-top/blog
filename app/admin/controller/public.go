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
	"github.com/webx-top/blog/app/base"
	"github.com/webx-top/blog/app/base/model"
	"github.com/webx-top/echo"
	X "github.com/webx-top/webx"
)

type Public struct {
	login    X.Mapper
	register X.Mapper
	logout   X.Mapper
	*base.Controller
	user *model.User
}

func (a *Public) Init(c echo.Context) error {
	a.Controller = base.NewController(c)
	a.user = model.NewUser(a.Context)
	return nil
}

func (a *Public) Login() error {
	ss := a.Session()
	if a.IsPost() {
		var uname, passwd = a.Form(`uname`), a.Form(`passwd`)
		if !a.VerifyCaptcha(a.Form(`captcha`)) {
			return a.SetErr(a.T(`验证码错误`)).Display()
		}
		u, err := a.user.Login(uname, passwd)
		if err != nil {
			return a.SetErr(err).Display()
		}
		ss.Set(`user`, u).Save()
		return a.GotoNext(`Index`, `Index`)
	}
	return a.Display()
}

func (a *Public) Register() error {
	if a.IsPost() {
		var uname, passwd = a.Form(`uname`), a.Form(`passwd`)
		if !a.VerifyCaptcha(a.Form(`captcha`)) {
			return a.SetErr(a.T(`验证码错误`)).Display()
		}
		var active = false
		u, err := a.user.Register(uname, passwd, active)
		if err != nil {
			return a.SetErr(err).Display()
		}
		if active {
			a.Session().Set(`user`, u).Save()
		}
		return a.GotoNext(`Index`, `Index`)
	}
	return a.Display()
}

func (a *Public) Logout() error {
	a.Session().Delete(`user`).Save()
	return a.GotoNext(`Index`, `Index`)
}
