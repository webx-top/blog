package controller

import (
	"github.com/webx-top/blog/app/admin/lib"
	"github.com/webx-top/blog/app/base"
	"github.com/webx-top/echo"
	"github.com/webx-top/webx/lib/middleware/session"
)

func New() *Base {
	return &Base{
		Controller: base.BaseCtl,
	}
}

type Base struct {
	session.Session
	*base.Controller
}

func (a *Base) Before(c echo.Context) error {
	a.Session = session.Default(c)
	if uid, ok := a.Session.Get(`uid`).(int64); !ok || uid < 1 {
		c.Redirect(301, lib.App.Url+`login`)
		a.Exit(c)
	}
	return a.Controller.Before(c)
}
