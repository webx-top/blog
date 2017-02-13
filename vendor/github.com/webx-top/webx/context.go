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
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/handler/mvc"
)

func init() {
	gob.Register(&Output{})
}

func NewContext(s *Server, c echo.Context) *Context {
	return &Context{
		Context: c,
		Server:  s,
	}
}

const (
	NO_PERM = -2 //无权限
	NO_AUTH = -1 //未登录
	FAILURE = 0  //操作失败
	SUCCESS = 1  //操作成功
)

type IniterFunc func(interface{}) error

type Output struct {
	Status  int
	Message interface{}
	For     interface{}
	Data    interface{}
}

type ContextInitial interface {
	Init(*mvc.Wrapper, interface{}, string) (error, bool)
}

type Context struct {
	//public
	echo.Context
	Server         *Server
	Module         *mvc.Module
	Output         *Output
	C              interface{}
	ControllerName string
	ActionName     string
	Tmpl           string

	//private
	middleware []func(IniterFunc) IniterFunc
	exit       bool
	body       []byte
}

func (c *Context) Reset(req engine.Request, resp engine.Response) {
	c.Context.Reset(req, resp)
	c.Context.SetSessionOptions(c.Server.SessionOptions)

	c.ControllerName = ``
	c.Module = nil
	c.ActionName = ``
	c.exit = false
	c.Output = &Output{1, ``, ``, echo.H{}}
	c.Tmpl = ``
	c.Format()
	c.middleware = nil
	c.body = nil
	c.C = nil
}

func (c *Context) AcceptFormat() string {
	return c.Format()
}

func (c *Context) Error(err error) {
	code := http.StatusInternalServerError
	msg := http.StatusText(code)
	if he, ok := err.(*echo.HTTPError); ok {
		c.SetCode(he.Code)
		msg = he.Message
	}
	if c.Object().Echo().Debug() {
		msg = err.Error()
	}
	if !c.Response().Committed() {
		if c.Request().Method() == echo.HEAD {
			c.NoContent(code)
		} else {
			c.SetErr(msg)
			c.Display(``)
		}
	}
	c.Echo().Logger().Debug(err)
}

func (c *Context) Init(wrp *mvc.Wrapper, controller interface{}, actName string) (error, bool) {
	c.Module = wrp.Module
	c.C = controller
	if c.Module.URLRecovery != nil {
		c.ControllerName = c.Module.URLRecovery(wrp.ControllerName)
		c.ActionName = c.Module.URLRecovery(actName)
	} else {
		c.ControllerName = wrp.ControllerName
		c.ActionName = actName
	}
	c.Tmpl = c.Module.Name + `/` + c.ControllerName + `/` + c.ActionName
	c.Context.SetRenderer(c.Module.Renderer)
	c.Context.SetFunc(`URLFor`, c.URLFor)
	c.Context.SetFunc(`BuildURL`, c.BuildURL)
	c.Context.SetFunc(`ModuleURLFor`, c.ModuleURLFor)
	c.Context.SetFunc(`ModuleURL`, c.ModuleURL)
	c.Context.SetFunc(`ControllerName`, func() string {
		return c.ControllerName
	})
	c.Context.SetFunc(`ActionName`, func() string {
		return c.ActionName
	})
	c.Context.SetFunc(`ModuleName`, func() interface{} {
		return c.Module.Name
	})
	c.Context.SetFunc(`ModuleRoot`, func() string {
		return c.Module.URL
	})
	c.Context.SetFunc(`ModuleDomain`, func() string {
		return c.Module.Domain
	})
	c.Context.SetFunc(`C`, func() interface{} {
		return c.C
	})
	return c.execMW(c.C), false
}

func (c *Context) execMW(ctl interface{}) error {
	var h IniterFunc = func(c interface{}) error {
		return nil
	}
	for i := len(c.middleware) - 1; i >= 0; i-- {
		h = c.middleware[i](h)
	}
	return h(ctl)
}

func (c *Context) Use(i ...func(IniterFunc) IniterFunc) {
	for _, v := range i {
		c.middleware = append(c.middleware, v)
	}
}

func (c *Context) SetSecCookie(key string, value interface{}) {
	if c.Server.Codec == nil {
		val, _ := value.(string)
		c.SetCookie(key, val)
		return
	}
	encoded, err := c.Server.Codec.Encode(key, value)
	if err != nil {
		c.Server.Core.Logger().Error(err)
	} else {
		c.SetCookie(key, encoded)
	}
}

func (c *Context) SecCookie(key string, value interface{}) {
	cookieValue := c.GetCookie(key)
	if cookieValue == `` {
		return
	}
	if c.Server.Codec != nil {
		err := c.Server.Codec.Decode(key, cookieValue, value)
		if err != nil {
			c.Server.Core.Logger().Error(err)
		}
		return
	}
	if v, ok := value.(*string); ok {
		*v = cookieValue
	}
}

func (c *Context) GetSecCookie(key string) (value string) {
	c.SecCookie(key, &value)
	return
}

func (c *Context) Body() ([]byte, error) {
	if c.body != nil {
		return c.body, nil
	}
	b := c.Request().Body()
	defer b.Close()
	body, err := ioutil.ReadAll(b)
	if err != nil {
		return nil, err
	}
	c.body = body
	return body, nil
}

func (c *Context) IP() string {
	return c.Request().RealIP()
}

func (c *Context) IsPjax() bool {
	return c.Header(`X-PJAX`) == `true`
}

func (c *Context) PjaxContainer() string {
	return c.Header(`X-PJAX-Container`)
}

func (c *Context) OnlyAjax() bool {
	return c.IsAjax() && !c.IsPjax()
}

// Refer returns http referer header.
func (c *Context) Refer() string {
	return c.Referer()
}

// SubDomain returns sub domain string.
// if aa.bb.domain.com, returns aa.bb .
func (c *Context) SubDomain() string {
	parts := strings.Split(c.Host(), `.`)
	if len(parts) >= 3 {
		return strings.Join(parts[:len(parts)-2], `.`)
	}
	return ``
}

func (c *Context) Assign(key string, val interface{}) *Context {
	data, _ := c.Output.Data.(echo.H)
	if data == nil {
		data = echo.H{}
	}
	data[key] = val
	c.Output.Data = data
	return c
}

func (c *Context) AssignX(values *map[string]interface{}) *Context {
	if values == nil {
		return c
	}
	data, _ := c.Output.Data.(echo.H)
	if data == nil {
		data = echo.H{}
	}
	for key, val := range *values {
		data[key] = val
	}
	c.Output.Data = data
	return c
}

func (c *Context) Exit(args ...bool) *Context {
	exit := true
	if len(args) > 0 {
		exit = args[0]
	}
	c.exit = exit
	return c
}

func (c *Context) IsExit() bool {
	return c.exit
}

func (c *Context) tmplPath(tpath string) string {
	if c.Module == nil {
		return tpath
	}
	if len(tpath) == 0 {
		return ``
	}
	if tpath[0] == '/' {
		tpath = c.Module.Name + tpath
	} else if !strings.Contains(tpath, `/`) {
		tpath = c.Module.Name + `/` + c.ControllerName + `/` + tpath
	}
	return tpath
}

func (c *Context) Display(args ...interface{}) error {
	if c.Response().Committed() {
		return nil
	}
	switch len(args) {
	case 2:
		if v, ok := args[0].(string); ok {
			c.Tmpl = c.tmplPath(v)
		}
		if v, ok := args[1].(int); ok && v > 0 {
			c.SetCode(v)
		}
	case 1:
		if v, ok := args[0].(int); ok {
			if v > 0 {
				c.SetCode(v)
			}
		} else if v, ok := args[0].(string); ok {
			c.Tmpl = c.tmplPath(v)
		}
	}
	if c.Code() <= 0 {
		c.SetCode(http.StatusOK)
	}
	if ignore, _ := c.Get(`webx:ignoreRender`).(bool); ignore {
		return nil
	}

	flash, ok := c.GetFlash()
	var err error
	switch c.Format() {
	case `xml`:
		err = c.XML(c.Output, c.Code())
	case `json`:
		if callback := c.Query(`callback`); callback != `` {
			err = c.JSONP(callback, c.Output, c.Code())
		} else {
			err = c.JSON(c.Output, c.Code())
		}
	default:
		if len(c.Tmpl) == 0 {
			err = c.HTML(fmt.Sprintf(`<pre code="%v" for="%v">%v</pre>`, c.Output.Status, c.Output.For, c.Output.Message), c.Code())
		} else {
			c.SetTmplDefaultFuncs(flash, ok)
			err = c.Render(c.Tmpl, c.Output.Data, c.Code())
		}
	}
	if err != nil {
		c.Server.Core.DefaultHTTPErrorHandler(c.ErrorWithCode(http.StatusInternalServerError, err.Error()), c)
		c.Exit()
	}
	return nil
}

func (c *Context) GetFlash() (flash *Output, ok bool) {
	flash, ok = c.Session().Get(`webx:flash`).(*Output)
	if ok {
		c.Session().Delete(`webx:flash`).Save()
	}
	return
}

// SetTmplDefaultFuncs 设置模板默认函数
func (c *Context) SetTmplDefaultFuncs(flash *Output, ok bool) {
	c.Context.SetFunc(`Status`, func() int {
		if ok {
			return flash.Status
		}
		return c.Output.Status
	})
	c.Context.SetFunc(`Message`, func() interface{} {
		if ok {
			return flash.Message
		}
		return c.Output.Message
	})
	c.Context.SetFunc(`For`, func() interface{} {
		if ok {
			return flash.For
		}
		return c.Output.For
	})
}

// MapForm 映射表单数据到结构体
// ParseStruct mapping forms' name and values to struct's field
// For example:
//		<form>
//			<input name=`user.id`/>
//			<input name=`user.name`/>
//			<input name=`user.age`/>
//		</form>
//
//		type User struct {
//			Id int64
//			Name string
//			Age string
//		}
//
//		var user User
//		err := c.MapForm(&user,`user`)
//
func (c *Context) MapForm(i interface{}, names ...string) error {
	var name string
	if len(names) > 0 {
		name = names[0]
	}
	return echo.NamedStructMap(c.Server.Core, i, c.Request().Form().All(), name)
}

// MapData 映射数据到结构体
func (c *Context) MapData(i interface{}, data map[string][]string, names ...string) error {
	var name string
	if len(names) > 0 {
		name = names[0]
	}
	return echo.NamedStructMap(c.Server.Core, i, data, name)
}

// ErrorWithCode 生成HTTPError
func (c *Context) ErrorWithCode(code int, args ...string) *echo.HTTPError {
	return echo.NewHTTPError(code, args...)
}

// SetOutput 设置输出(code,message,for,data)
func (c *Context) SetOutput(code int, args ...interface{}) *Context {
	c.Output.Status = code
	var hasData bool
	switch len(args) {
	case 3:
		c.Output.Data = args[2]
		hasData = true
		fallthrough
	case 2:
		c.Output.For = args[1]
		fallthrough
	case 1:
		c.Output.Message = args[0]
		if !hasData {
			flash := &Output{
				Status:  c.Output.Status,
				Message: c.Output.Message,
				For:     c.Output.For,
				Data:    nil,
			}
			c.Session().Set(`webx:flash`, flash).Save()
		}
	}
	return c
}

// SetSuc 设置响应类型为“操作成功”(message,for,data)
func (c *Context) SetSuc(args ...interface{}) *Context {
	return c.SetOutput(SUCCESS, args...)
}

// SetSucData 设置成功返回的数据
func (c *Context) SetSucData(data interface{}) *Context {
	return c.SetOutput(SUCCESS, ``, ``, data)
}

// SetErr 设置出错类型为“操作失败”(message,for,data)
func (c *Context) SetErr(args ...interface{}) *Context {
	return c.SetOutput(FAILURE, args...)
}

// SetNoAuth 设置出错类型为“未登录”(message,for,data)
func (c *Context) SetNoAuth(args ...interface{}) *Context {
	return c.SetOutput(NO_AUTH, args...)
}

// SetNoPerm 设置出错类型为“未授权”(message,for,data)
func (c *Context) SetNoPerm(args ...interface{}) *Context {
	return c.SetOutput(NO_PERM, args...)
}

// ModuleURLFor 生成指定Module网址
func (c *Context) ModuleURLFor(ppath string, args ...map[string]interface{}) string {
	return c.Server.URLs.BuildFromPath(ppath, args...)
}

// ModuleURL 生成指定Module网址
func (c *Context) ModuleURL(mod string, ctl string, act string, args ...interface{}) string {
	return c.Server.URLs.Build(mod, ctl, act, args...)
}

// URLFor 生成当前Module网址
func (c *Context) URLFor(ppath string, args ...map[string]interface{}) string {
	if len(ppath) == 0 {
		if len(c.ControllerName) > 0 {
			ppath = c.ControllerName + `/`
		}
		ppath += c.ActionName
		return c.Server.URLs.BuildFromPath(c.Module.Name+`/`+ppath, args...)
	}
	ppath = strings.TrimLeft(ppath, `/`)
	return c.Server.URLs.BuildFromPath(c.Module.Name+`/`+ppath, args...)
}

// BuildURL 生成当前Module网址
func (c *Context) BuildURL(ctl string, act string, args ...interface{}) string {
	return c.Server.URLs.Build(c.Module.Name, ctl, act, args...)
}

// TmplPath 生成模板路径 args: ActionName,ControllerName,ModuleName
func (c *Context) TmplPath(args ...string) string {
	var mod, ctl, act = c.Module.Name, c.ControllerName, c.ActionName
	switch len(args) {
	case 3:
		mod = args[2]
		fallthrough
	case 2:
		ctl = args[1]
		fallthrough
	case 1:
		act = args[0]
	}
	return mod + `/` + ctl + `/` + act
}

// SetTmpl 指定要渲染的模板路径
func (c *Context) SetTmpl(args ...string) *Context {
	c.Tmpl = c.TmplPath(args...)
	return c
}

// Atoe 字符串转error
func (c *Context) Atoe(v string) error {
	return errors.New(v)
}

// NextURL 获取下一步网址
func (c *Context) NextURL(defaultURL ...string) string {
	next := c.GetNextURL()
	if len(next) == 0 {
		next = c.U(defaultURL...)
	}
	return next
}

// GotoNext 跳转到下一步
func (c *Context) GotoNext(defaultURL ...string) error {
	return c.Redir(c.NextURL(defaultURL...))
}

// GetNextURL 自动获取下一步网址
func (c *Context) GetNextURL() string {
	next := c.Form(`next`)
	if len(next) > 0 {
		return c.ParseNextURL(next)
	}
	next = c.Referer()
	if len(next) > 0 {
		if strings.HasSuffix(next, c.Request().URI()) {
			next = ``
		}
	}
	return next
}

// ParseNextURL 解析下一步网址
func (c *Context) ParseNextURL(next string) string {
	if len(next) == 0 {
		return next
	}
	if next[0] == '!' {
		next = next[1:]
		next = strings.Replace(next, `-`, `/`, -1)
		next = strings.Replace(next, ` `, `+`, -1)
		for strings.HasSuffix(next, `_`) {
			next = strings.TrimSuffix(next, `_`) + `=`
		}
		var err error
		next, err = com.Base64Decode(next)
		if err != nil {
			c.Server.Core.Logger().Error(err)
		}
	}
	return next
}

// GenNextURL 生成安全编码后的下一步网址
func (c *Context) GenNextURL(u string) string {
	if len(u) == 0 {
		return ``
	}
	if u[0] == '!' {
		return u
	}
	u = com.Base64Encode(u)
	for strings.HasSuffix(u, `=`) {
		u = strings.TrimSuffix(u, `=`) + `_`
	}
	u = strings.Replace(u, `/`, `-`, -1)
	return `!` + u
}

// U 网址生成
func (c *Context) U(args ...string) (s string) {
	var p string
	switch len(args) {
	case 3:
		if args[2][0] != '?' {
			return c.ModuleURL(args[0], args[1], args[2])
		}
		p = args[2]
		fallthrough
	case 2:
		if len(p) > 0 || args[1][0] != '?' {
			return c.BuildURL(args[0], args[1]) + p
		}
		p = args[1]
		fallthrough
	case 1:
		size := len(args[0])
		if len(p) > 0 || (size > 0 && args[0][0] != '?') {
			if size > 0 {
				switch args[0][0] {
				case '/': //usage: /webx/index => {module}/webx/index
					s = c.URLFor(args[0])
					return s + p
				case ':': //usage: :http://webx.top => http://webx.top
					s = args[0][1:]
					return s + p
				}
			}
			if strings.Contains(args[0], `/`) {
				s = c.ModuleURLFor(args[0])
			} else {
				s = c.ModuleURL(c.Module.Name, c.ControllerName, args[0])
			}
			return s + p
		}
		p = args[0]
		fallthrough
	case 0:
		s = c.ModuleURL(c.Module.Name, c.ControllerName, c.ActionName) + p
	}
	return
}

// Redir 页面跳转
func (c *Context) Redir(url string, args ...interface{}) error {
	var code = http.StatusFound //302. 307:http.StatusTemporaryRedirect
	if len(args) > 0 {
		if v, ok := args[0].(bool); ok && v {
			code = http.StatusMovedPermanently
		} else if v, ok := args[0].(int); ok {
			code = v
		}
	}
	c.Exit()
	if c.Format() != `html` {
		c.Set(`webx:ignoreRender`, false)
		c.Assign(`Location`, url)
		return c.Display()
	}
	return c.Context.Redirect(url, code)
}

// Goto 页面跳转(根据路由生成网址后跳转)
func (c *Context) Goto(goURL string, args ...interface{}) error {
	goURL = c.U(goURL)
	return c.Redir(goURL, args...)
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
