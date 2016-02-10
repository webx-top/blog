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
	"strings"

	"github.com/webx-top/blog/app/base/lib/database"
	"github.com/webx-top/echo"
	X "github.com/webx-top/webx"
	"github.com/webx-top/webx/lib/com"
	"github.com/webx-top/webx/lib/htmlcache"
	"github.com/webx-top/webx/lib/i18n"
	"github.com/webx-top/webx/lib/middleware/jwt"
	"github.com/webx-top/webx/lib/middleware/language"
	"github.com/webx-top/webx/lib/middleware/session"
	"github.com/webx-top/webx/lib/session/ssi"
	"github.com/webx-top/webx/lib/xsrf"

	"github.com/webx-top/blog/app/base/lib/config"

	"github.com/admpub/confl"
)

var (
	Server      *X.Server
	SessionMW   echo.MiddlewareFunc
	HtmlCache   *htmlcache.Config
	HtmlCacheMW echo.MiddlewareFunc
	I18n        *i18n.I18n
	Xsrf        *xsrf.Xsrf
	Jwt         *jwt.JWT
	DB          *database.Orm

	DefaultLang = `zh-cn`
	Project     = `blog`
	RootDir     = com.SelfDir()
	Language    = language.NewLanguage()
	theme       = `default`
	templateDir = RootDir + `/data/theme/`
	Config      = &config.Config{}
	configFile  = RootDir + `/data/config/config.yaml`
)

func init() {
	LoadConfig(configFile)

	// ======================
	// 初始化默认Server
	// ======================
	Server = X.Serv(Project).InitTmpl(ThemePath())
	ApplyConfig()
	SessionMW = session.Middleware(&ssi.Options{
		Engine:   Server.SessionStoreEngine,
		Path:     `/`,
		Domain:   Server.CookieDomain,
		MaxAge:   int(Server.CookieExpires),
		Secure:   false,
		HttpOnly: Server.CookieHttpOnly,
	}, Server.SessionStoreConfig)
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
	HtmlCacheMW = HtmlCache.Middleware()
	I18n = i18n.New(RootDir+`/data/lang/rules`, RootDir+`/data/lang/messages`, DefaultLang, DefaultLang)
	Xsrf = xsrf.New()
	Jwt = jwt.New(Server.CookieAuthKey)

	Server.Pprof()
	Server.Debug(true)
	Server.SetHook(Language.DetectURI)

	// ======================
	// 监控语言文件更改
	// ======================
	monitorLanguageResource()

	// ======================
	// 连接数据库
	// ======================
	var err error
	DB, err = database.NewOrm(Config.DB.Engine, Config.DB.Dsn())
	if err == nil {
		DB.SetPrefix(Config.DB.Prefix)
	}
}

func monitorLanguageResource() {
	var callback = com.MonitorEventFunc{
		Modify: func(file string) {
			Server.Core.Logger().Info("reload language: %v", file)
			I18n.Reload(file)
		},
	}
	go com.Monitor(RootDir+`/data/lang/messages`, callback, func(f string) bool {
		Server.Core.Logger().Info("changed: %v", f)
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

func ApplyConfig() {

	Server.SessionStoreEngine = Config.Session.StoreEngine
	Server.SessionStoreConfig = Config.Session.StoreConfig

	Server.CookieDomain = Config.Cookie.Domain
	Server.CookieExpires = Config.Cookie.Expires
	Server.CookieHttpOnly = Config.Cookie.HttpOnly
	Server.CookiePrefix = Config.Cookie.Prefix
	Server.CookieAuthKey = Config.Cookie.AuthKey
	Server.CookieBlockKey = Config.Cookie.BlockKey

	Server.InitCodec([]byte(Server.CookieAuthKey), []byte(Server.CookieBlockKey))

	if Config.Language.AllList != nil {
		for _, lang := range Config.Language.AllList {
			Language.Set(lang, true, lang == Config.Language.Default)
		}
	} else {
		Language.Set(Config.Language.Default, true, true)
		Language.Set(`en`, true)
	}

}
