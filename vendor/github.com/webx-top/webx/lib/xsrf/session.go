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
package xsrf

import (
	"github.com/webx-top/echo"
)

type SessionStorage struct {
}

func (c *SessionStorage) Get(key string, ctx echo.Context) string {
	s := ctx.Session()
	if s == nil {
		return ""
	}
	val, _ := s.Get(key).(string)
	return val
}

func (c *SessionStorage) Set(key, val string, ctx echo.Context) {
	s := ctx.Session()
	if s == nil {
		return
	}
	s.Set(key, val)
	s.Save()
}

func (c *SessionStorage) Valid(key, val string, ctx echo.Context) bool {
	return c.Get(key, ctx) == val
}
