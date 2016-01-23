package base

import (
	//"fmt"
	"io/ioutil"
	"strings"

	"github.com/webx-top/blog/app/base/lib/database"
	X "github.com/webx-top/webx"
	"github.com/webx-top/webx/lib/com"
	"github.com/webx-top/webx/lib/htmlcache"
	"github.com/webx-top/webx/lib/i18n"
	"github.com/webx-top/webx/lib/middleware/jwt"
	"github.com/webx-top/webx/lib/middleware/language"
	"github.com/webx-top/webx/lib/middleware/session"
	"github.com/webx-top/webx/lib/xsrf"

	"github.com/webx-top/blog/app/base/lib/config"

	"github.com/admpub/confl"
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
	HtmlCacheMW = HtmlCache.Middleware()
	BaseCtl     = NewController()
	I18n        = i18n.New(RootDir+`/data/lang/rules`, RootDir+`/data/lang/messages`, DefaultLang, DefaultLang)
	Xsrf        = xsrf.New()
	Jwt         = jwt.New(SecretKey)
	Config      = &config.Config{}
	configFile  = RootDir + `/data/config/config.yaml`
	DB          *database.Orm
)

func init() {
	LoadConfig(configFile)

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

	// ======================
	// 连接数据库
	// ======================
	DB, _ = database.NewOrm(Config.DB.Engine, Config.DB.Dsn())
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

func LoadConfig(file string) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	} else {
		err = confl.Unmarshal(content, Config)
		if err != nil {
			panic(err)
		}
	}
}
