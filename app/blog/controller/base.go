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
package controller

import (
	"github.com/webx-top/blog/app/base"
	//"github.com/webx-top/blog/app/blog/lib"
	"github.com/webx-top/echo"
)

func New(c echo.Context) *Base {
	return &Base{
		Controller: base.NewController(c),
	}
}

type Base struct {
	*base.Controller
}

func (a *Base) Before() error {
	/*
		if uid, ok := a.GetSession(`uid`).(int64); !ok || uid < 1 {
			a.Redir(a.App.Url+`login`)
			a.Exit()
		}
	*/
	return nil
}
