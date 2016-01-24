package controller

import (
	"github.com/webx-top/blog/app/base"
	"github.com/webx-top/blog/app/blog/lib"
	X "github.com/webx-top/webx"
)

var indexCtl = &Index{Controller: base.NewController()}

func init() {
	c := lib.App.RC(indexCtl)
	c.R(`/`, indexCtl.Index)
	c2 := &Test{}
	lib.App.RC(c2).AutoRoute()
}

type Index struct {
	*base.Controller
}

func (a *Index) Index(c *X.Context) error {
	a.Tmpl(`index`, c)
	return nil
}

func (a *Index) After(c *X.Context) error {
	return a.Controller.After(c)
}

type Test struct {
	*X.Context
	*X.App
}

func (a *Test) Init(c *X.Context, app *X.App) {
	a.Context = c
	a.App = app
}

func (a *Test) Index_GET_POST() error {
	a.Assign(`path`, a.Context.Path())
	return a.Display()
}
