/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the `License`);
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an `AS IS` BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package webx

import (
	"encoding/gob"
	"reflect"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/handler/mvc"
)

func init() {
	gob.Register(&Output{})
}

func NewContext(s *Server, c echo.Context) *Context {
	return &Context{
		Context: mvc.NewContext(s.Application, c),
		Server:  s,
	}
}

type ContextInitial interface {
	Init(*mvc.Wrapper, interface{}, string) (error, bool)
}

type Context struct {
	*mvc.Context
	Server *Server
}

func (c *Context) Reset(req engine.Request, resp engine.Response) {
	c.Context.Reset(req, resp)
	c.SetData(&Output{
		context: c.Context,
		Status:  1,
		Message: ``,
		For:     ``,
		Data:    echo.H{},
	})
}

func (c *Context) AcceptFormat() string {
	return c.Format()
}

func (c *Context) Assign(key string, val interface{}) *Context {
	c.Context.Assign(key, val)
	return c
}

func (c *Context) Assignx(values *map[string]interface{}) *Context {
	c.Context.Assignx(values)
	return c
}

func (c *Context) Exit(args ...bool) *Context {
	c.Context.Exit(args...)
	return c
}

func (c *Context) AssignX(values *map[string]interface{}) *Context {
	c.Assignx(values)
	return c
}

// SetOutput 设置输出(status,message,for,data)
func (c *Context) SetOutput(code int, args ...interface{}) *Context {
	c.Data().Set(code, args...)
	return c
}

// SetSuc 设置响应类型为“操作成功”(message,for,data)
func (c *Context) SetSuc(args ...interface{}) *Context {
	c.Context.SetSuc(args...)
	return c
}

// SetSucData 设置成功返回的数据
func (c *Context) SetSucData(data interface{}) *Context {
	c.Context.SetSucData(data)
	return c
}

// SetErr 设置出错类型为“操作失败”(message,for,data)
func (c *Context) SetErr(args ...interface{}) *Context {
	c.Context.SetErr(args...)
	return c
}

// SetNoAuth 设置出错类型为“未登录”(info,zone,data)
func (c *Context) SetNoAuth(args ...interface{}) *Context {
	c.Context.SetNoAuth(args...)
	return c
}

// SetNoPerm 设置出错类型为“未授权”(message,for,data)
func (c *Context) SetNoPerm(args ...interface{}) *Context {
	c.Context.SetNoPerm(args...)
	return c
}

// SetTmpl 指定要渲染的模板路径
func (c *Context) SetTmpl(args ...string) *Context {
	c.Context.SetTmpl(args...)
	return c
}

// A 调用控制器方法
func (c *Context) A(ctl string, act string) (err error) {
	a := c.Module.Wrapper(`controller.` + ctl)
	if a == nil {
		return c.Atoe(`Controller "` + ctl + `" does not exist.`)
	}
	k := `webx.controller.reflect.type:` + ctl
	var e reflect.Type
	if t, ok := c.Get(k).(reflect.Type); ok {
		e = t
	} else {
		e = reflect.Indirect(reflect.ValueOf(a.Controller)).Type()
		c.Set(k, e)
	}
	return a.Exec(c, e, act)
}
