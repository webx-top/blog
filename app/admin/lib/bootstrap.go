package lib

import (
	X "bitbucket.org/admpub/webx"
	"bitbucket.org/admpub/webx/lib/tplex"
	"github.com/webx-top/blog/app/base"
)

var App = X.Serv().NewApp("admin", base.Language.Store(), base.SessionMW)

func init() {
	te := tplex.New(base.ThemePath(`admin`))
	te.InitMgr(true, true)
	App.Webx().SetRenderer(te)
}
