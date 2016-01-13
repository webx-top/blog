package controller

import (
	"github.com/webx-top/blog/app/admin/lib"
	"github.com/webx-top/echo"
)

var publicCtrl = &Public{}

func init() {
	c := lib.App.RC(publicCtrl)
	c.R(`/login`, publicCtrl.Login)
}

type Public struct {
}

func (a *Public) Login(c *echo.Context) error {
	return c.Render(200, `login`, `test`)
}
