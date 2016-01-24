package controller

import (
	"github.com/webx-top/blog/app/base"
	"github.com/webx-top/blog/app/blog/lib"
	X "github.com/webx-top/webx"
)

func New() *Base {
	return &Base{
		Controller: base.BaseCtl,
	}
}

type Base struct {
	*base.Controller
}

func (a *Base) Before(c *X.Context) error {
	if uid, ok := c.GetSession(`uid`).(int64); !ok || uid < 1 {
		c.Redirect(301, lib.App.Url+`login`)
		c.Exit = true
	}
	return a.Controller.Before(c)
}

func (a *Base) After(c *X.Context) error {
	return nil
}
