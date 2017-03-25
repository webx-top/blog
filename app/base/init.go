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
	"github.com/admpub/log"
	"github.com/webx-top/echo"
	mw "github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/language"
	"github.com/webx-top/echo/middleware/session"
	boltStore "github.com/webx-top/echo/middleware/session/engine/bolt"
	cookieStore "github.com/webx-top/echo/middleware/session/engine/cookie"
	X "github.com/webx-top/webx"
	"github.com/webx-top/webx/lib/database"
	"github.com/webx-top/webx/lib/middleware/jwt"
	"github.com/webx-top/webx/lib/static/htmlcache"
	"github.com/webx-top/webx/lib/xsrf"

	_ "github.com/webx-top/echo/middleware/render/pongo2"
	_ "github.com/webx-top/webx/lib/client/list/datatable"
	localStore "github.com/webx-top/webx/lib/store/file/local"

	"github.com/webx-top/webx/lib/config"
)

var (
	Server    *X.Server
	SessionMW echo.MiddlewareFuncd
	HtmlCache *htmlcache.Config
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
	//Server.Core.Use(mw.Gzip())
	Server.RootModuleName = `blog`
	err := Server.LoadConfig(Server.RootDir()+`/data/config/config.yaml`, Config)
	if err != nil {
		panic(err)
	}
	if Config.FrontendTemplate.Theme == `` {
		Config.FrontendTemplate.Theme = `default`
	}
	if Config.BackendTemplate.Theme == `` {
		Config.BackendTemplate.Theme = `admin`
	}
	Server.TemplateDir = Server.RootDir() + `/data/theme/`
	Server.SetTheme(&Config.FrontendTemplate)
	Server.InitStatic()
	Server.Pprof().Debug(true)
	Server.Core.PreUse(Language.Middleware())
	Server.SetSessionOptions(&Config.Session.SessionOptions)
	Server.InitCodec([]byte(Config.Session.AuthKey), []byte(Config.Session.BlockKey))

	// ======================
	// 设置Session中间件
	// ======================
	cookieStore.RegWithOptions(&cookieStore.CookieOptions{
		KeyPairs: [][]byte{
			[]byte(Config.Session.AuthKey),
			[]byte(Config.Session.BlockKey),
		},
		SessionOptions: Server.SessionOptions,
	})

	boltStore.RegWithOptions(&boltStore.BoltOptions{
		File: Server.RootDir() + `/data/bolt/session.db`,
		KeyPairs: [][]byte{
			[]byte(Config.Session.AuthKey),
			[]byte(Config.Session.BlockKey),
		},
		BucketName:     `session`,
		SessionOptions: Server.SessionOptions,
	})
	SessionMW = session.Middleware(Server.SessionOptions)

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
	Xsrf = xsrf.New()
	Jwt = jwt.New(Config.Session.AuthKey)
	Language.Init(&Config.Language)

	// ======================
	// 连接数据库
	// ======================
	Config.DBSlaves = []*config.DB{&Config.DB} //testing
	DB, err = database.NewOrm(Config.DB.Engine, &Config.DB, Config.DBSlaves...)
	if err != nil {
		log.Error(err)
	}

	store := localStore.New(map[string]string{
		"SavePath":  Server.RootDir() + `/data/upload/`,
		"PublicUrl": `/upload/`,
		"RootPath":  Server.RootDir(),
	})
	localStore.RegStore(store)
	Server.Core.Use(mw.Static(&mw.StaticOptions{Path: `/upload`, Root: Server.RootDir() + `/data/upload/`}))
}
