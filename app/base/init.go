/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package base

import (
	//"fmt"
	"io/ioutil"

	"github.com/webx-top/echo"
	X "github.com/webx-top/webx"
	"github.com/webx-top/webx/lib/com"
	"github.com/webx-top/webx/lib/database"
	"github.com/webx-top/webx/lib/htmlcache"
	"github.com/webx-top/webx/lib/i18n"
	"github.com/webx-top/webx/lib/middleware/jwt"
	"github.com/webx-top/webx/lib/middleware/language"
	"github.com/webx-top/webx/lib/middleware/session"
	"github.com/webx-top/webx/lib/session/ssi"
	"github.com/webx-top/webx/lib/tplfunc"
	"github.com/webx-top/webx/lib/xsrf"

	_ "github.com/webx-top/webx/lib/client/datatable"
	_ "github.com/webx-top/webx/lib/tplex/pongo2"

	"github.com/webx-top/webx/lib/config"

	"github.com/admpub/confl"
)

var (
	Server        *X.Server
	SessionMW     echo.MiddlewareFunc
	HtmlCache     *htmlcache.Config
	I18n          *i18n.I18n
	Xsrf          *xsrf.Xsrf
	Jwt           *jwt.JWT
	DB            *database.Orm
	Static        *tplfunc.Static
	FuncMap       map[string]interface{}
	AbsThemePath  string
	AbsStaticPath string

	DefaultLang = `zh-cn`
	Project     = `blog`
	RootDir     = com.SelfDir()
	Language    = language.New()
	theme       = `default`
	templateDir = RootDir + `/data/theme/`
	Config      = &config.Config{}
	configFile  = RootDir + `/data/config/config.yaml`
	StaticPath  = `/assets`
)

func init() {
	LoadConfig(configFile)

	theme = Config.FrontendTemplate.Theme

	AbsThemePath = ThemePath()
	AbsStaticPath = AbsThemePath + StaticPath

	// ======================
	// 初始化默认Server
	// ======================
	Server = X.Serv(Project).ResetTmpl(AbsThemePath, Config.FrontendTemplate.Engine)
	Server.Core.PreUse(Language.Middleware())
	ApplyConfig()
	SessionMW = session.Middleware(&ssi.Options{
		Engine:   Server.Session.StoreEngine,
		Path:     `/`,
		Domain:   Server.Cookie.Domain,
		MaxAge:   int(Server.Cookie.Expires),
		Secure:   false,
		HttpOnly: Server.Cookie.HttpOnly,
	}, Server.Session.StoreConfig)
	/*
		map[string]string{
			"file": RootDir + `/data/bolt/session.db`,
			"key":  Server.CookieAuthKey,
			"name": Server.Name,
		})
	*/

	HtmlCache = &htmlcache.Config{
		HtmlCacheDir:   RootDir + `/data/html`,
		HtmlCacheOn:    true,
		HtmlCacheRules: make(map[string]interface{}),
		HtmlCacheTime:  86400,
	}
	I18n = i18n.New(Config.Language)
	Xsrf = xsrf.New()
	Jwt = jwt.New(Server.Cookie.AuthKey)

	Server.Pprof()
	Server.Debug(true)

	FuncMap = Server.FuncMap()
	Server.TemplateEngine.SetFuncMapFn(func() map[string]interface{} {
		return FuncMap
	})
	Static = Server.Static(StaticPath, AbsStaticPath, &FuncMap)
	Server.TemplateEngine.MonitorEvent(Static.OnUpdate(AbsThemePath))

	// ======================
	// 连接数据库
	// ======================
	var err error
	DB, err = database.NewOrm(Config.DB.Engine, Config.DB.Dsn())
	if err == nil {
		DB.SetPrefix(Config.DB.Prefix)
	}
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
	Server.ResetTmpl(ThemePath(args...))
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
	if Config.FrontendTemplate.Theme == `` {
		Config.FrontendTemplate.Theme = `default`
	}
	if Config.BackendTemplate.Theme == `` {
		Config.BackendTemplate.Theme = `admin`
	}
	//fmt.Printf("%#v\n", Config)
}

func ApplyConfig() {
	Server.Session = &Config.Session
	Server.Cookie = &Config.Cookie
	Server.InitCodec([]byte(Server.Cookie.AuthKey), []byte(Server.Cookie.BlockKey))
	Language.Init(Config.Language)
}
