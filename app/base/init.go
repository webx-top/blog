package base

import (
	"strings"

	X "github.com/webx-top/webx"
	"github.com/webx-top/webx/lib/com"
	"github.com/webx-top/webx/lib/htmlcache"
	"github.com/webx-top/webx/lib/i18n"
	"github.com/webx-top/webx/lib/middleware/jwt"
	"github.com/webx-top/webx/lib/middleware/language"
	"github.com/webx-top/webx/lib/middleware/session"
	"github.com/webx-top/webx/lib/xsrf"
)

var (
	SecretKey   = `webx.top`
	DefaultLang = `zh-cn`
	Project     = `blog`
	RootDir     = com.SelfDir()
	Language    = language.NewLanguage()
	SessionMW   = session.Middleware(`cookie`, SecretKey)
	theme       = `default`
	templateDir = RootDir + `/data/theme/`
	Server      = X.Serv(Project).InitTmpl(ThemePath())
	HtmlCache   = &htmlcache.Config{
		HtmlCacheDir:   RootDir + `/data/html`,
		HtmlCacheOn:    true,
		HtmlCacheRules: make(map[string]interface{}),
		HtmlCacheTime:  86400,
	}
	HtmlCacheMW = HtmlCache.Middleware(Server.TemplateEngine)
	BaseCtl     = NewController()
	I18n        = i18n.New(RootDir+`/data/lang/rules`, RootDir+`/data/lang/messages`, DefaultLang, DefaultLang)
	Xsrf        = xsrf.New()
	Jwt         = jwt.New(SecretKey)
)

func init() {

	// ======================
	// 初始化默认Server
	// ======================
	Language.Set(DefaultLang, true, true)
	Language.Set(`en`, true)
	Server.Pprof()
	Server.Debug(true)
	Server.SetHook(Language.DetectURI)

	// ======================
	// 监控语言文件更改
	// ======================
	moniterLanguageResource()
}

func moniterLanguageResource() {
	var callback = com.MoniterEventFunc{
		Modify: func(file string) {
			Server.Echo.Logger().Info("reload language: %v", file)
			I18n.Reload(file)
		},
	}
	go com.Moniter(RootDir+`/data/lang/messages`, callback, func(f string) bool {
		return strings.HasSuffix(f, `.yaml`)
	})
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
	Server.InitTmpl(ThemePath(args...))
}
