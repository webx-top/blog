package base

import (
	X "github.com/webx-top/webx"
	"github.com/webx-top/webx/lib/i18n"
)

func NewController(c *X.Context) *Controller {
	a := &Controller{
		Controller: X.NewController(c),
	}
	return a
}

type Controller struct {
	*X.Controller
}

func (a *Controller) Init(c *X.Context) {
}

func (a *Controller) Before() error {
	return a.Controller.Before()
}

func (a *Controller) After() error {
	return a.Controller.After()
}

func (a *Controller) Lang() string {
	if a.Language == `` {
		a.Language = DefaultLang
	}
	return a.Language
}

func (a *Controller) T(key string, args ...interface{}) string {
	return i18n.T(a.Lang(), key, args...)
}
