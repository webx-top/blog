package lib

import (
	"html/template"

	X "bitbucket.org/admpub/webx"
	"bitbucket.org/admpub/webx/lib/tplex"
	"github.com/webx-top/blog/app/base"
	"github.com/webx-top/blog/app/base/lib/tplfunc"
)

var (
	Name       = `admin`
	App        = X.Serv().NewApp(Name, base.Language.Store(), base.SessionMW)
	FuncMap    = tplfunc.TplFuncMap
	StaticPath = `/assets`
	Static     = tplfunc.NewStatic(`/` + Name + StaticPath)
)

func init() {
	tp := base.ThemePath(`admin`)
	te := tplex.New(tp)
	te.InitMgr(true, true)

	FuncMap["StaticUrl"] = Static.StaticUrl
	FuncMap["JsUrl"] = Static.JsUrl
	FuncMap["CssUrl"] = Static.CssUrl
	FuncMap["ImgUrl"] = Static.ImgUrl
	FuncMap["JsTag"] = Static.JsTag
	FuncMap["CssTag"] = Static.CssTag
	FuncMap["ImgTag"] = Static.ImgTag
	te.FuncMapFn = func() template.FuncMap {
		return FuncMap
	}
	x := App.Webx()
	x.SetRenderer(te)
	x.Static(StaticPath, tp+StaticPath)
}
