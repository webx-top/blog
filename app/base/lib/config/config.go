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
package config

import (
	"strings"

	"github.com/webx-top/webx/lib/com"
)

type DB struct {
	Engine  string
	User    string
	Pass    string
	Name    string
	Host    string
	Port    string
	Charset string
	Prefix  string
}

func (a *DB) Dsn() string {
	var dsn string
	switch a.Engine {
	case `mysql`:
		var host string
		if strings.HasPrefix(a.Host, `unix:`) {
			host = "unix(" + strings.TrimPrefix(a.Host, `unix:`) + ")"
		} else {
			if a.Port == `` {
				a.Port = "3306"
			}
			host = "tcp(" + a.Host + ":" + a.Port + ")"
		}
		dsn = com.UrlEncode(a.User) + ":" + com.UrlEncode(a.Pass) + "@" + host + "/" + a.Name + "?charset=" + a.Charset
	case `mymysql`: //tcp:localhost:3306*gotest/root/root
		var host string
		if strings.HasPrefix(a.Host, `unix:`) {
			host = a.Host
		} else {
			if a.Port == `` {
				a.Port = "3306"
			}
			host = "tcp:" + a.Host + ":" + a.Port
		}
		dsn = host + "*" + a.Name + "/" + com.UrlEncode(a.User) + "/" + com.UrlEncode(a.Pass)
	default:
		panic(a.Engine + ` is not supported.`)
	}
	return dsn
}

type Cookie struct {
	Prefix   string
	HttpOnly bool
	AuthKey  string
	BlockKey string
	Expires  int64
	Domain   string
}

type Session struct {
	StoreEngine string
	StoreConfig interface{}
}

type Language struct {
	Default string
	AllList []string
}

type Template struct {
	Theme  string
	Engine string
	Style  string
}

type Config struct {
	DB               `json:"DB"`
	Cookie           `json:"Cookie"`
	Session          `json:"Session"`
	Language         `json:"Language"`
	FrontendTemplate Template `json:"FrontendTemplate"`
	BackendTemplate  Template `json:"BackendTemplate"`
}
