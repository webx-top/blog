package lib

import (
	"html/template"

	//X "github.com/webx-top/webx"
	"github.com/webx-top/blog/app/base"
	"github.com/webx-top/webx/lib/tplex"
	"github.com/webx-top/webx/lib/tplfunc"
)

var (
	Name       = `admin`
	App        = base.Server.NewApp(Name, base.Language.Store(), base.SessionMW, base.Xsrf.Middleware() /*, base.Jwt.Validate()*/)
	FuncMap    = base.Server.FuncMap()
	StaticPath = `/assets`
	Static     *tplfunc.Static
)

func init() {
	tp := base.ThemePath(`admin`)
	te := tplex.New(tp)
	te.InitMgr(true, true)
	Static = base.Server.Static(`/`+Name+StaticPath, tp+StaticPath, &FuncMap)
	FuncMap["AppUrl"] = func(p ...string) string {
		if len(p) > 0 {
			return App.Url + p[0]
		}
		return App.Url
	}
	te.FuncMapFn = func() template.FuncMap {
		return FuncMap
	}
	te.FileChangeEvent = Static.OnUpdate(te.TemplateDir)
	x := App.Webx()
	x.SetRenderer(te)
	x.Static(StaticPath, tp+StaticPath)
}
