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
	X "github.com/webx-top/webx"
	"github.com/webx-top/com"
)

type SecCookieStorage struct {
	HighSafe bool
}

func (c *SecCookieStorage) Get(key string, ctx echo.Context) string {
	return X.X(ctx).GetSecCookie(key)
}

func (c *SecCookieStorage) Set(key, val string, ctx echo.Context) {
	X.X(ctx).SetSecCookie(key, c.safelyVal(val, ctx))
}

func (c *SecCookieStorage) safelyVal(val string, ctx echo.Context) string {
	if c.HighSafe {
		context := X.X(ctx)
		val = val + `|` + context.IP() + `|` + context.Request().UserAgent()
	}
	return val
}

func (c *SecCookieStorage) Valid(key, val string, ctx echo.Context) bool {
	return c.Get(key, ctx) == c.safelyVal(val, ctx)
}

type CookieStorage struct {
	AuthKey  string
	HighSafe bool
}

func (c *CookieStorage) Get(key string, ctx echo.Context) string {
	return X.X(ctx).GetCookie(key)
}

func (c *CookieStorage) Set(key, val string, ctx echo.Context) {
	X.X(ctx).SetCookie(key, c.safelyVal(val, ctx))
}

func (c *CookieStorage) safelyVal(val string, ctx echo.Context) string {
	if c.AuthKey != `` {
		if c.HighSafe {
			context := X.X(ctx)
			val = com.Md5(val + `|` + context.IP() + `|` + context.Request().UserAgent() + `|` + c.AuthKey)
		} else {
			val = com.Md5(val + `|` + c.AuthKey)
		}
	}
	return val
}

func (c *CookieStorage) Valid(key, val string, ctx echo.Context) bool {
	return c.Get(key, ctx) == c.safelyVal(val, ctx)
}
