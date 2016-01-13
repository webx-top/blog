package controller

import (
	"github.com/webx-top/blog/app/admin/lib"
	"github.com/webx-top/echo"
)

var indexCtrl = &Index{Base: &Base{}}

func init() {
	c := lib.App.RC(indexCtrl)
	c.R(`/`, indexCtrl.Index)
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
