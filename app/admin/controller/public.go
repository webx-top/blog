package controller

import (
	"github.com/webx-top/blog/app/admin/lib"
	"github.com/webx-top/blog/app/base"
	"github.com/webx-top/echo"
)

var publicCtl = &Public{Controller: base.BaseCtl}

func init() {
	c := lib.App.RC(publicCtl)
	c.R(`/login`, publicCtl.Login, `GET`, `POST`)
	c.R(`/logout/:next`, publicCtl.Logout, `GET`, `POST`)
}

type Public struct {
	*base.Controller
}

func (a *Public) Login(c echo.Context) error {
	a.Tmpl(`login`, c)
	return nil
}

func (a *Public) Logout(c echo.Context) error {
	a.Tmpl(`login`, c)
	return nil
}
