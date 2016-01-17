package lib

import (
	"github.com/webx-top/blog/app/base"
)

var App = base.Server.NewApp("", base.Language.Store(), base.SessionMW, base.HtmlCacheMW)
