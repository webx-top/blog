package controller

import (
	"github.com/admpub/echo"
	"github.com/webx-top/blog/app/admin/lib"
)

var indexCtrl = &Index{}

func init() {
	c := lib.App.RC(indexCtrl)
	c.R(`/`, indexCtrl.Index).R(`/login`, indexCtrl.Login)
}

type Index struct {
}

func (a *Index) Index(c *echo.Context) error {
	return c.Render(200, `index`, `test`)
}

func (a *Index) Login(c *echo.Context) error {
	return c.Render(200, `login`, `test`)
}
