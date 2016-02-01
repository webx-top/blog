package model

import (
	"errors"

	D "github.com/webx-top/blog/app/base/dbschema"
	"github.com/webx-top/webx/lib/com"
)

func NewUser(lang string) *User {
	return &User{M: NewM(lang)}
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
	} else if !com.CheckPassword(passwd, u.Passwd, u.Salt) {
		err = errors.New(a.T("用户名或密码不正确"))
	}
	return
}
