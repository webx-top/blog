package controller

import (
	"github.com/webx-top/blog/app/admin/lib"
	"github.com/webx-top/echo"
	"github.com/webx-top/webx/lib/middleware/session"
)

type Base struct {
	session.Session
	Uid int64
}

func (a *Base) Before(c *echo.Context) error {
	a.Session = session.Default(c)
	a.Uid, _ = a.Session.Get(`uid`).(int64)
	if a.Uid < 1 {
		c.Redirect(301, lib.App.Url+`login`)
		c.Set(`web:exit`, true)
	}
	return nil
}
