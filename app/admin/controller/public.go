package controller

import (
	"github.com/webx-top/blog/app/admin/lib"
	"github.com/webx-top/blog/app/base"
	"github.com/webx-top/echo"
)

var publicCtrl = &Public{Controller: base.BaseCtrl}

func init() {
	c := lib.App.RC(publicCtrl)
	c.R(`/login`, publicCtrl.Login)
}

type Public struct {
	*base.Controller
}

func (a *Public) Login(c *echo.Context) error {
	a.Tmpl(`login`, c)
	return nil
}
