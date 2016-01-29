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
	"github.com/webx-top/blog/app/admin/lib"
	"github.com/webx-top/blog/app/base"
	X "github.com/webx-top/webx"
)

func init() {
	lib.App.RC(&Public{}).Auto()
}

type Public struct {
	login X.Mapper
	*base.Controller
}

func (a *Public) Init(c *X.Context) {
	a.Controller = base.NewController(c)
}

func (a *Public) Login() error {
	return nil
}

func (a *Public) Logout(c *X.Context) error {
	return nil
}
