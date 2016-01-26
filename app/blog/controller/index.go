package controller

import (
	//"github.com/webx-top/blog/app/base"
	"github.com/webx-top/blog/app/blog/lib"
	X "github.com/webx-top/webx"
)

func init() {
	lib.App.RC(&Index{}).Auto()
}

type Index struct {
	index X.Mapper
	*Base
}

func (a *Index) Init(c *X.Context) {
	a.Base = New(c)
}

func (a *Index) Index() error {
	return nil
}

func (a *Index) After() error {
	return a.Controller.After()
}
