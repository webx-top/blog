package lib

import (
	"html/template"

	//X "github.com/webx-top/webx"
	"github.com/webx-top/blog/app/base"
	"github.com/webx-top/blog/app/base/lib/tplfunc"
	"github.com/webx-top/webx/lib/tplex"
)

var (
	Name       = `admin`
	App        = base.Server.NewApp(Name, base.Language.Store(), base.SessionMW)
	FuncMap    = tplfunc.TplFuncMap
	StaticPath = `/assets`
	Static     = tplfunc.NewStatic(`/` + Name + StaticPath)
)

func init() {
	tp := base.ThemePath(`admin`)
	te := tplex.New(tp)
	te.InitMgr(true, true)
	FuncMap = Static.Register(FuncMap)
	FuncMap["Lang"] = func() string {
		return `zh-cn`
	}
	FuncMap["AppUrl"] = func(p ...string) string {
		if len(p) > 0 {
			return App.Url + p[0]
		}
		return App.Url
	}
	FuncMap["RootUrl"] = func(p ...string) string {
		if len(p) > 0 {
			return base.Server.Url + p[0]
		}
		return base.Server.Url
	}
	te.FuncMapFn = func() template.FuncMap {
		return FuncMap
	}
	x := App.Webx()
	x.SetRenderer(te)
	x.Static(StaticPath, tp+StaticPath)
}
