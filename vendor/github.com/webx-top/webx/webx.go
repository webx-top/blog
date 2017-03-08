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
package webx

import (
	"reflect"

	"github.com/webx-top/echo/handler/mvc"
)

// VERSION 框架版本号
const VERSION = `1.1.0`

var (
	// 默认服务名称
	defaultServName = "webx"

	// 默认服务实例
	serv = NewServer(defaultServName)

	// 服务实例集合
	servs = Servers{}

	mapperType = reflect.TypeOf(Mapper{})
)

// 结构体中定义路由的字段类型
type Mapper struct{}

// Servers 服务集合
type Servers map[string]*Server

// Get 根据服务名称获取服务实例
func (s Servers) Get(name string) (sv *Server) {
	sv, _ = s[name]
	return
}

// Set 根据名称设置服务实例
func (s Servers) Set(name string, sv *Server) {
	s[name] = sv
}

// Serv 获取服务实例
func Serv(args ...string) *Server {
	if len(args) > 0 {
		if sv, ok := servs[args[0]]; ok {
			return sv
		}
		return NewServer(args[0])
	}
	return serv
}

// Module 获取模块实例，不传递参数时返回根模块
func Module(args ...string) *mvc.Module {
	return serv.Module(args...)
}

// Register 注册路由到根模块
func Register(p string, h interface{}, methods ...string) *mvc.Module {
	return serv.Module().Register(p, h, methods...)
}

// Use 注册控制器到根App
func Use(args ...interface{}) *mvc.Module {
	return serv.Module().Use(args...)
}

// Run 启动服务
func Run(args ...interface{}) {
	serv.Run(args...)
}
