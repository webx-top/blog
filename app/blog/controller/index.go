package controller

import (
	"github.com/webx-top/blog/app/base"
	"github.com/webx-top/blog/app/blog/lib"
	"github.com/webx-top/echo"
)

var indexCtl = &Index{Controller: base.NewController()}

func init() {
	c := lib.App.RC(indexCtl)
	c.R(`/`, indexCtl.Index)
}

type Index struct {
	*base.Controller
}

func (a *Index) Index(c *echo.Context) error {
	a.Tmpl(`index`, c)
	return nil
}

func (a *Index) After(c *echo.Context) error {
	return nil
}
