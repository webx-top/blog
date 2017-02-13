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
	"net/http"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/logger"
)

func NewController(c echo.Context) *Controller {
	a := &Controller{}
	a.Init(c)
	return a
}

type Controller struct {
	*Context
	logger.Logger
}

func (a *Controller) Init(c echo.Context) error {
	a.Context = a.X(c)
	a.Logger = c.Echo().Logger()
	a.SetFunc("Queryx", a.Queryx)
	a.SetFunc("Formx", a.Formx)
	a.SetFunc("Path", a.Path)
	return nil
}

func (a *Controller) X(c echo.Context) *Context {
	return X(c)
}

func (a *Controller) NotFound(args ...string) error {
	code := http.StatusNotFound
	text := http.StatusText(code)
	a.Context.Exit()
	if len(args) > 0 {
		text = args[0]
	}
	return a.ErrorWithCode(code, text)
}
