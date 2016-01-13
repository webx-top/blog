package controller

import (
	"github.com/webx-top/blog/app/blog/lib"
	"github.com/webx-top/echo"
)

var indexCtrl = &Index{}

func init() {
	c := lib.App.RC(indexCtrl)
	c.R(`/`, indexCtrl.Index)
}

type Index struct {
}

func (a *Index) Index(c *echo.Context) error {
	c.Set(`Tmpl`, `index`)
	c.Set(`Data`, `test`)
	return nil
}
