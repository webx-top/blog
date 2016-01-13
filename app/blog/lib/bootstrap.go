package lib

import (
	"github.com/webx-top/blog/app/base"
	X "github.com/webx-top/webx"
)

var App = X.Serv().NewApp("", base.Language.Store(), base.SessionMW, base.HtmlCacheMW)
