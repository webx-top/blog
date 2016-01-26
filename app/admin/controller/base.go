package controller

import (
	//"github.com/webx-top/blog/app/admin/lib"
	"github.com/webx-top/blog/app/base"
	X "github.com/webx-top/webx"
)

func New(c *X.Context) *Base {
	return &Base{
		Controller: base.NewController(c),
	}
}

type Base struct {
	*base.Controller
}

func (a *Base) Before() error {
	if uid, ok := a.GetSession(`uid`).(int64); !ok || uid < 1 {
		a.Redirect(301, a.App.Url+`public/login`)
		a.Exit = true
	}
	return a.Controller.Before()
}
