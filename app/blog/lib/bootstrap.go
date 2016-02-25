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
package lib

import (
	"github.com/webx-top/blog/app/base"
	X "github.com/webx-top/webx"
)

var (
	App = base.Server.NewApp("", base.SessionMW, base.HtmlCache.Middleware())
)

func init() {
	App.RealName = "blog"
	App.R(`/ping`, func(c *X.Context) error {
		return c.String(200, `pong`)
	})
}
