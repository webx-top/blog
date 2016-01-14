package controller

import (
	"github.com/webx-top/blog/app/base"
	"github.com/webx-top/blog/app/blog/lib"
	"github.com/webx-top/echo"
	"github.com/webx-top/webx/lib/middleware/session"
)

func New() *Base {
	return &Base{
		Controller: base.BaseCtrl,
	}
}

type Base struct {
	session.Session
	Uid int64
	*base.Controller
}

func (a *Base) Before(c *echo.Context) error {
	a.Session = session.Default(c)
	a.Uid, _ = a.Session.Get(`uid`).(int64)
	if a.Uid < 1 {
		c.Redirect(301, lib.App.Url+`login`)
		a.Exit(c)
	}
	return a.Controller.Before(c)
}
