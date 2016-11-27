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
	"errors"
	"strings"

	D "github.com/webx-top/blog/app/base/dbschema"
	"github.com/webx-top/com"
	X "github.com/webx-top/webx"
)

func NewUser(ctx *X.Context) *User {
	return &User{M: NewM(ctx)}
}

type User struct {
	*M
}

func (a *User) Login(uname string, passwd string) (u *D.User, err error) {
	u = &D.User{}
	var has bool
	has, err = a.M.DB.Where(`uname=?`, uname).Get(u)
	if err != nil {
	} else if !has {
		err = errors.New(a.T("用户不存在"))
	} else if u.Active == "N" {
		err = errors.New(a.T("您的账号还没有激活，无法登录"))
	} else if !com.CheckPassword(passwd, u.Passwd, u.Salt) {
		err = errors.New(a.T("用户名或密码不正确"))
	}
	return
}

func (a *User) Register(uname string, passwd string, active bool) (u *D.User, err error) {
	uname = strings.TrimSpace(uname)
	passwd = strings.TrimSpace(passwd)
	if len(passwd) < 6 {
		err = errors.New(a.T("密码不能少于6个字符"))
		return
	}
	for _, v := range []string{`@`, `,`, `*`, `:`, `|`, `/`, `=`, `&`, `?`} {
		if strings.Contains(uname, v) {
			err = errors.New(a.T("用户名不能包含“" + v + "”"))
			return
		}
	}
	u = &D.User{}
	var has bool
	has, err = a.M.DB.Where(`uname=?`, uname).Get(u)
	if has {
		err = errors.New(a.T("用户已经存在"))
		return
	}
	u.Salt = com.Salt()
	u.Passwd = com.MakePassword(passwd, u.Salt)
	u.Uname = uname
	if active {
		u.Active = "Y"
	} else {
		u.Active = "N"
	}
	_, err = a.DB.Insert(u)

	return
}

func (a *User) Delete(id int) (affected int64, err error) {
	affected, err = a.DB.Id(id).Delete(&D.User{})
	return
}
