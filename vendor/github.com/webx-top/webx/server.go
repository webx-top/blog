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
package webx

import (
	"strings"

	codec "github.com/gorilla/securecookie"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/handler/mvc"
	"github.com/webx-top/echo/handler/pprof"
	"github.com/webx-top/echo/middleware/tplfunc"
)

func HandlerWrapper(h interface{}) echo.Handler {
	if handle, y := h.(func(*Context) error); y {
		return echo.HandlerFunc(func(c echo.Context) error {
			return handle(X(c))
		})
	}
	return nil
}

func NewServer(name string, middlewares ...interface{}) (s *Server) {
	s = &Server{
		MVC: mvc.NewWithContext(name, func(e *echo.Echo) echo.Context {
			return NewContext(s, echo.NewContext(nil, nil, e))
		}, middlewares...),
	}
	s.MVC.Core.SetHandlerWrapper(HandlerWrapper)
	s.FuncMap = s.DefaultFuncMap()
	s.ContextInitial = func(ctx echo.Context, wrp *mvc.Wrapper, controller interface{}, actionName string) (err error, exit bool) {
		return ctx.(ContextInitial).Init(wrp, controller, actionName)
	}
	servs.Set(name, s)
	return
}

type Server struct {
	*mvc.MVC
	Codec codec.Codec `json:"-" xml:"-"`
}

// InitCodec 初始化 加密/解密 接口
func (s *Server) InitCodec(hashKey []byte, blockKey []byte) {
	s.Codec = codec.New(hashKey, blockKey)
}

// Pprof 启用pprof
func (s *Server) Pprof() *Server {
	pprof.Wrapper(s.Core)
	return s
}

func (s *Server) DefaultFuncMap() (r map[string]interface{}) {
	r = tplfunc.New()
	r["RootURL"] = func(p ...string) string {
		if len(p) > 0 {
			return s.URL + p[0]
		}
		return s.URL
	}
	return
}

// Tree  module -> controller -> action -> HTTP-METHODS
func (s *Server) Tree(args ...*echo.Echo) (r map[string]map[string]map[string]map[string]string) {
	core := s.Core
	if len(args) > 0 {
		core = args[0]
	}
	nrs := core.NamedRoutes()
	rs := core.Routes()
	r = map[string]map[string]map[string]map[string]string{}
	for name, indexes := range nrs {
		p := strings.LastIndex(name, `/`)
		s := strings.Split(name[p+1:], `.`)
		var appName, ctlName, actName string
		switch len(s) {
		case 3:
			if !strings.HasSuffix(s[2], `-fm`) {
				continue
			}
			actName = strings.TrimSuffix(s[2], `-fm`)
			ctlName = strings.TrimPrefix(s[1], `(`)
			ctlName = strings.TrimPrefix(ctlName, `*`)
			ctlName = strings.TrimSuffix(ctlName, `)`)
			p2 := strings.LastIndex(name[0:p], `/`)
			appName = name[p2+1 : p]
		default:
			continue
		}
		if _, ok := r[appName]; !ok {
			r[appName] = map[string]map[string]map[string]string{}
		}
		if _, ok := r[appName][ctlName]; !ok {
			r[appName][ctlName] = map[string]map[string]string{}
		}
		if _, ok := r[appName][ctlName][actName]; !ok {
			r[appName][ctlName][actName] = map[string]string{}
			for _, index := range indexes {
				route := rs[index]
				if _, ok := r[appName][ctlName][actName][route.Method]; !ok {
					r[appName][ctlName][actName][route.Method] = route.Method
				}
			}
		}
	}
	return
}
