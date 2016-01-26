package controller

import (
	"github.com/webx-top/blog/app/admin/lib"
	X "github.com/webx-top/webx"
	//"github.com/webx-top/webx/lib/com"
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

func (a *Index) Before() error {
	return a.Base.Before()
}

func (a *Index) Index() error {
	return nil
}
