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
	"github.com/webx-top/webx/lib/tplex/pongo2"
	"github.com/webx-top/webx/lib/tplfunc"
)

var (
	App        = base.Server.NewApp("", base.Language.Store(), base.SessionMW, base.HtmlCacheMW)
	FuncMap    = base.Server.FuncMap()
	StaticPath = `/assets`
	Static     *tplfunc.Static
)

func init() {
	tp := base.ThemePath(`default`)
	te := pongo2.New(tp)
	te.Init(true, true)
	Static = base.Server.Static(StaticPath, tp+StaticPath, &FuncMap)
	te.SetFuncMapFn(func() map[string]interface{} {
		return FuncMap
	})
	te.MonitorEvent(Static.OnUpdate(tp))
	x := App.Webx()
	x.SetRenderer(te)
	x.Static(StaticPath, tp+StaticPath)

	App.R(`/ping`, func(c *X.Context) error {
		return c.String(200, `pong`)
	})
}
