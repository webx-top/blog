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
	"reflect"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/handler/mvc"
)

func HandlerWrapper(h interface{}) echo.Handler {
	if handle, y := h.(func(*Context) error); y {
		return echo.HandlerFunc(func(c echo.Context) error {
			return handle(X(c))
		})
	}
	return nil
}

func NewServer(name string) (s *Server) {
	s = &Server{
		Application: mvc.NewWithContext(name, func(e *echo.Echo) echo.Context {
			return NewContext(s, echo.NewContext(nil, nil, e))
		}),
	}
	s.Core.SetHandlerWrapper(HandlerWrapper)
	s.ContextInitial = func(ctx echo.Context, wrp *mvc.Wrapper, controller interface{}, actionName string) (err error, exit bool) {
		return ctx.(ContextInitial).Init(wrp, controller, actionName)
	}
	s.MapperCheck = func(t reflect.Type) bool {
		return t == mapperType
	}
	servs.Set(name, s)
	return
}

type Server struct {
	*mvc.Application
}
