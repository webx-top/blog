package controller

import (
	"github.com/webx-top/blog/app/admin/lib"
	"github.com/webx-top/echo"
	//"github.com/webx-top/webx/lib/com"
)

var indexCtl = &Index{Base: New()}

func init() {
	c := lib.App.RC(indexCtl)
	c.R(`/`, indexCtl.Index)
}

type Index struct {
	*Base
}

func (a *Index) Before(c echo.Context) error {
	return a.Base.Before(c)
}

func (a *Index) Index(c echo.Context) error {
	a.Tmpl(`index`, c)
	return nil
}
