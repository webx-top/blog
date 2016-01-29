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
	Engine    string
	DbUser    string
	DbPass    string
	DbName    string
	DbHost    string
	DbPort    string
	DbCharset string
}

func (a *DB) Dsn() string {
	var dsn string
	switch a.Engine {
	case `mysql`:
		var host string
		if strings.HasPrefix(a.DbHost, `unix:`) {
			host = "unix(" + strings.TrimPrefix(a.DbHost, `unix:`) + ")"
		} else {
			if a.DbPort == `` {
				a.DbPort = "3306"
			}
			host = "tcp(" + a.DbHost + ":" + a.DbPort + ")"
		}
		dsn = com.UrlEncode(a.DbUser) + ":" + com.UrlEncode(a.DbPass) + "@" + host + "/" + a.DbName + "?charset=" + a.DbCharset
	case `mymysql`: //tcp:localhost:3306*gotest/root/root
		var host string
		if strings.HasPrefix(a.DbHost, `unix:`) {
			host = "unix:" + strings.TrimPrefix(a.DbHost, `unix:`)
		} else {
			if a.DbPort == `` {
				a.DbPort = "3306"
			}
			host = "tcp:" + a.DbHost + ":" + a.DbPort
		}
		dsn = host + "*" + a.DbName + "/" + com.UrlEncode(a.DbUser) + "/" + com.UrlEncode(a.DbPass)
	default:
		panic(a.Engine + ` is not supported.`)
	}
	return dsn
}

type Config struct {
	DB `json:"DB"`
}
