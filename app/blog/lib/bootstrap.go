package lib

import (
	X "bitbucket.org/admpub/webx"
	"github.com/webx-top/blog/app/base"
)

var App = X.Serv().NewApp("", base.Language.Store(), base.SessionMW, base.HtmlCacheMW)
