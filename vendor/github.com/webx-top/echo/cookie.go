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
package echo

import (
	"net/http"
	"net/url"
	"time"
)

type CookieOptions struct {
	Prefix string

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'.
	// MaxAge>0 means Max-Age attribute present and given in seconds.
	MaxAge int

	Path     string
	Domain   string
	Secure   bool
	HttpOnly bool
}

type Cookier interface {
	Get(key string) string
	Set(key string, val string, args ...interface{})
}

func NewCookier(ctx Context) Cookier {
	return &cookie{
		context: ctx,
	}
}

func NewCookie(name string, value string, opts ...*CookieOptions) *Cookie {
	opt := &CookieOptions{}
	if len(opts) > 0 {
		opt = opts[0]
	}
	if len(opt.Path) == 0 {
		opt.Path = `/`
	}
	cookie := &Cookie{
		cookie: &http.Cookie{
			Name:     opt.Prefix + name,
			Value:    value,
			Path:     opt.Path,
			Domain:   opt.Domain,
			MaxAge:   opt.MaxAge,
			Secure:   opt.Secure,
			HttpOnly: opt.HttpOnly,
		},
	}
	return cookie
}

type Cookie struct {
	cookie *http.Cookie
}

func (c *Cookie) Path(p string) *Cookie {
	c.cookie.Path = p
	return c
}

func (c *Cookie) Domain(p string) *Cookie {
	c.cookie.Domain = p
	return c
}

func (c *Cookie) MaxAge(p int) *Cookie {
	c.cookie.MaxAge = p
	return c
}

func (c *Cookie) Expires(p int64) *Cookie {
	if p > 0 {
		c.cookie.Expires = time.Unix(time.Now().Unix()+p, 0)
	} else if p < 0 {
		c.cookie.Expires = time.Unix(1, 0)
	}
	return c
}

func (c *Cookie) Secure(p bool) *Cookie {
	c.cookie.Secure = p
	return c
}

func (c *Cookie) HttpOnly(p bool) *Cookie {
	c.cookie.HttpOnly = p
	return c
}

func (c *Cookie) Send(ctx Context) {
	ctx.Response().Header().Set(HeaderSetCookie, c.cookie.String())
}

type cookie struct {
	context Context
}

func (c *cookie) Get(key string) string {
	var val string
	if v := c.context.Request().Cookie(c.context.CookieOptions().Prefix + key); v != `` {
		val, _ = url.QueryUnescape(v)
	}
	return val
}

func (c *cookie) Set(key string, val string, args ...interface{}) {
	val = url.QueryEscape(val)
	cookie := NewCookie(key, val, c.context.CookieOptions())
	switch len(args) {
	case 5:
		httpOnly, _ := args[4].(bool)
		cookie.HttpOnly(httpOnly)
		fallthrough
	case 4:
		secure, _ := args[3].(bool)
		cookie.Secure(secure)
		fallthrough
	case 3:
		domain, _ := args[2].(string)
		cookie.Domain(domain)
		fallthrough
	case 2:
		ppath, _ := args[1].(string)
		cookie.Path(ppath)
		fallthrough
	case 1:
		var liftTime int64
		switch args[0].(type) {
		case int:
			liftTime = int64(args[0].(int))
		case int64:
			liftTime = args[0].(int64)
		case time.Duration:
			liftTime = int64(args[0].(time.Duration).Seconds())
		}
		cookie.Expires(liftTime)
	}
	cookie.Send(c.context)
}
