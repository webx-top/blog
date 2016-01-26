package controller

import (
	"github.com/webx-top/blog/app/admin/lib"
	"github.com/webx-top/blog/app/base"
	X "github.com/webx-top/webx"
)

func init() {
	lib.App.RC(&Public{}).Auto()
}

type Public struct {
	login X.Mapper
	*base.Controller
}

func (a *Public) Init(c *X.Context) {
	a.Controller = base.NewController(c)
}

func (a *Public) Login() error {
	return nil
}

func (a *Public) Logout(c *X.Context) error {
	return nil
}
