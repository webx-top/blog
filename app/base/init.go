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
	"github.com/webx-top/echo"
	mw "github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/session"
	boltStore "github.com/webx-top/echo/middleware/session/engine/bolt"
	cookieStore "github.com/webx-top/echo/middleware/session/engine/cookie"
	X "github.com/webx-top/webx"
	"github.com/webx-top/webx/lib/database"
	"github.com/webx-top/webx/lib/i18n"
	"github.com/webx-top/webx/lib/middleware/jwt"
	"github.com/webx-top/webx/lib/middleware/language"
	"github.com/webx-top/webx/lib/static/htmlcache"
	"github.com/webx-top/webx/lib/xsrf"

	_ "github.com/webx-top/webx/lib/client/list/datatable"
	localStore "github.com/webx-top/webx/lib/store/file/local"
	_ "github.com/webx-top/webx/lib/tplex/pongo2"

	"github.com/webx-top/webx/lib/config"
)

var (
	Server    *X.Server
	SessionMW echo.MiddlewareFuncd
	HtmlCache *htmlcache.Config
	I18n      *i18n.I18n
	Xsrf      *xsrf.Xsrf
	Jwt       *jwt.JWT
	DB        *database.Orm

	Project  = `blog`
	Language = language.New()
	Config   = &config.Config{}
	Version  = `1.0.0`
)

func init() {

	// ======================
	// 初始化默认Server
	// ======================
	Server = X.Serv(Project)
	Server.RootAppName = `blog`
	err := Server.LoadConfig(Server.RootDir()+`/data/config/config.yaml`, Config)
	if err != nil {
		panic(err)
	}
	if Config.FrontendTemplate.Theme == `` {
		Config.FrontendTemplate.Theme = `default`
	}
	if Config.BackendTemplate.Theme == `` {
		Config.FrontendTemplate.Theme = `admin`
	}
	Server.TemplateDir = Server.RootDir() + `/data/theme/`
	Server.SetTheme(&Config.FrontendTemplate)
	Server.InitStatic()
	Server.Pprof().Debug(true)
	Server.Core.PreUse(Language.Middleware())

	Server.Session = &Config.Session
	Server.Cookie = &Config.Cookie
	Server.InitCodec([]byte(Server.Cookie.AuthKey), []byte(Server.Cookie.BlockKey))

	// ======================
	// 设置Session中间件
	// ======================
	sessionOptions := &echo.SessionOptions{
		Engine: Server.Session.StoreEngine,
		Name:   `SESSIONID`,
		CookieOptions: &echo.CookieOptions{
			Path:     `/`,
			Domain:   Server.Cookie.Domain,
			MaxAge:   Server.Cookie.MaxAge,
			Secure:   false,
			HttpOnly: Server.Cookie.HttpOnly,
		},
	}
	cookieStore.RegWithOptions(&cookieStore.CookieOptions{
		KeyPairs: [][]byte{
			[]byte(Server.Session.StoreConfig.(string)),
		},
		SessionOptions: sessionOptions,
	})
	boltStore.RegWithOptions(&boltStore.BoltOptions{
		File: Server.RootDir() + `/data/bolt/session`,
		KeyPairs: [][]byte{
			[]byte(Server.Session.StoreConfig.(string)),
		},
		BucketName:     `session`,
		SessionOptions: sessionOptions,
	})
	SessionMW = session.Middleware(sessionOptions)
	/*
		map[string]string{
			"file": RootDir + `/data/bolt/session.db`,
			"key":  Server.CookieAuthKey,
			"name": Server.Name,
		})
	*/

	// ======================
	// 设置静态页缓存
	// ======================
	HtmlCache = &htmlcache.Config{
		HtmlCacheDir:   Server.RootDir() + `/data/html`,
		HtmlCacheOn:    true,
		HtmlCacheRules: make(map[string]interface{}),
		HtmlCacheTime:  86400,
	}

	// ======================
	// 设置其它常用功能组件
	// ======================
	I18n = i18n.New(&Config.Language)
	Xsrf = xsrf.New()
	Jwt = jwt.New(Server.Cookie.AuthKey)
	Language.Init(&Config.Language)

	// ======================
	// 连接数据库
	// ======================
	DB, err = database.NewOrm(Config.DB.Engine, Config.DB.Dsn())
	if err == nil {
		DB.SetPrefix(Config.DB.Prefix)
	}

	store := localStore.New(map[string]string{
		"SavePath":  Server.RootDir() + `/data/upload/`,
		"PublicUrl": `/upload/`,
		"RootPath":  Server.RootDir(),
	})
	localStore.RegStore(store)
	Server.Core.Use(mw.Static(&mw.StaticOptions{Path: `/upload`, Root: Server.RootDir() + `/data/upload/`}))
}
