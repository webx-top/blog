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
package session

import (
	codec "github.com/gorilla/securecookie"
	"github.com/webx-top/echo"
	ss "github.com/webx-top/echo/middleware/session/engine"
	cookieStore "github.com/webx-top/echo/middleware/session/engine/cookie"
)

func NewSession(options *echo.SessionOptions, ctx echo.Context) echo.Sessioner {
	return ss.NewSession(options, ctx)
}

type Store interface {
	ss.Store
}

func NewMySession(store ss.Store, name string, ctx echo.Context) echo.Sessioner {
	return ss.NewMySession(store, name, ctx)
}

func StoreEngine(options *echo.SessionOptions) (store Store) {
	store = ss.StoreEngine(options)
	if store == nil {
		cs := cookieStore.New(&cookieStore.CookieOptions{
			KeyPairs: [][]byte{
				[]byte(codec.GenerateRandomKey(32)),
				[]byte(codec.GenerateRandomKey(32)),
			},
			SessionOptions: options,
		})
		cookieStore.Reg(cs)
		store = cs
	}
	return
}
