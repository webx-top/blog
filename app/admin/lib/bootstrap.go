package lib

import (
	"html/template"
	"path"
	"strings"

	//X "github.com/webx-top/webx"
	"github.com/webx-top/blog/app/base"
	"github.com/webx-top/webx/lib/com"
	"github.com/webx-top/webx/lib/tplex"
	"github.com/webx-top/webx/lib/tplfunc"
)

var (
	Name       = `admin`
	App        = base.Server.NewApp(Name, base.Language.Store(), base.SessionMW)
	FuncMap    = tplfunc.TplFuncMap
	StaticPath = `/assets`
	Static     *tplfunc.Static
)

func init() {
	tp := base.ThemePath(`admin`)
	te := tplex.New(tp)
	te.InitMgr(true, true)
	Static = tplfunc.NewStatic(`/`+Name+StaticPath, tp+StaticPath)
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
	te.FileChangeEvent = func(name string) {
		name = path.Join(te.TemplateDir, name)
		name = strings.TrimPrefix(com.FixDirSeparator(name), com.FixDirSeparator(Static.RootPath)+`/`)
		//base.Server.Echo.Logger().Info(`file change: %v`, name)
		Static.DeleteCombined(name)
	}
	x := App.Webx()
	x.SetRenderer(te)
	x.Static(StaticPath, tp+StaticPath)
}
