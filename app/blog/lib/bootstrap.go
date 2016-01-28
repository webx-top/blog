package lib

import (
	"github.com/webx-top/blog/app/base"
	X "github.com/webx-top/webx"
)

var App = base.Server.NewApp("", base.Language.Store(), base.SessionMW, base.HtmlCacheMW)

func init() {
	App.R(`/ping`, func(c *X.Context) error {
		return c.String(200, `pong`)
	})
}
