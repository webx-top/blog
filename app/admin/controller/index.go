package controller

import (
	"github.com/webx-top/blog/app/admin/lib"
	"github.com/webx-top/echo"
)

var indexCtrl = &Index{Base: &Base{}}

func init() {
	c := lib.App.RC(indexCtrl)
	c.R(`/`, indexCtrl.Index).R(`/login`, indexCtrl.Login)
}

type Index struct {
	*Base
}

func (a *Index) Before(c *echo.Context) error {
	return a.Base.Before(c)
}

func (a *Index) Index(c *echo.Context) error {
	return c.Render(200, `index`, `test`)
}

func (a *Index) Login(c *echo.Context) error {
	return c.Render(200, `login`, `test`)
}
