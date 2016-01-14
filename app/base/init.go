package base

import (
	X "github.com/webx-top/webx"
	"github.com/webx-top/webx/lib/com"
	"github.com/webx-top/webx/lib/htmlcache"
	mw "github.com/webx-top/webx/lib/middleware"
	"github.com/webx-top/webx/lib/middleware/session"
)

var (
	RootDir     = com.SelfDir()
	Language    = mw.NewLanguage()
	SessionMW   = session.Middleware(`cookie`, `webx.top`)
	theme       = `default`
	templateDir = RootDir + `/data/theme/`
	Server      = X.Serv().InitTmpl(ThemePath()).Pprof().Debug(true)
	HtmlCache   = &htmlcache.Config{
		HtmlCacheDir:   RootDir + `/data/html`,
		HtmlCacheOn:    true,
		HtmlCacheRules: make(map[string]interface{}),
		HtmlCacheTime:  86400,
	}
	HtmlCacheMW = HtmlCache.Middleware(Server.TemplateEngine)
	BaseCtrl    = NewController()
)

func init() {

	// ======================
	// 初始化默认Server
	// ======================
	Language.Set(`zh-cn`, true, true)
	Language.Set(`en`, true)
	Server.Pprof()
	Server.Debug(true)
	Server.SetHook(Language.DetectURI)
}

func ThemePath(args ...string) string {
	if len(args) < 1 {
		return templateDir + theme
	}
	return templateDir + args[0]
}

func SetTheme(args ...string) {
	if len(args) > 1 && args[0] == `admin` {
		return
	}
	X.Serv().InitTmpl(ThemePath(args...))
}
