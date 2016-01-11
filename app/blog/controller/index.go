package controller

import (
	"github.com/admpub/echo"
	"github.com/webx-top/blog/app/blog/lib"
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
