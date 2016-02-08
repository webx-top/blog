package model

import (
	"errors"
	"strings"

	D "github.com/webx-top/blog/app/base/dbschema"
	X "github.com/webx-top/webx"
	"github.com/webx-top/webx/lib/com"
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
