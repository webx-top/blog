package controller

import (
	"github.com/webx-top/blog/app/base"
	"github.com/webx-top/blog/app/blog/lib"
	"github.com/webx-top/echo"
)

func New() *Base {
	return &Base{
		Controller: base.BaseCtl,
	}
}

type Base struct {
	*base.Controller
}

func (a *Base) Before(c echo.Context) error {
	if uid, ok := a.X(c).GetSession(`uid`).(int64); !ok || uid < 1 {
		c.Redirect(301, lib.App.Url+`login`)
		a.Exit(c)
	}
	return a.Controller.Before(c)
}

func (a *Base) After(c echo.Context) error {
	return nil
}
