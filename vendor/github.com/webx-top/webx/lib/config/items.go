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

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/webx/lib/get"
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
			host = "tcp(" + a.Host
			if len(a.Port) > 0 {
				host += ":" + a.Port
			}
			host += ")"
		}
		dsn = com.UrlEncode(a.User) + ":" + com.UrlEncode(a.Pass) + "@" + host + "/" + a.Name + "?charset=" + a.Charset
	case `mymysql`: //tcp:localhost:3306*gotest/root/root
		var host string
		if strings.HasPrefix(a.Host, `unix:`) {
			host = a.Host
		} else {
			host = "tcp:" + a.Host
			if len(a.Port) > 0 {
				host += ":" + a.Port
			}
		}
		dsn = host + "*" + a.Name + "/" + com.UrlEncode(a.User) + "/" + com.UrlEncode(a.Pass)
	default:
		panic(a.Engine + ` is not supported.`)
	}
	return dsn
}

type Session struct {
	echo.SessionOptions
	AuthKey  string
	BlockKey string
}

type Items map[string]string

func (i Items) Get(key string) string {
	if v, y := i[key]; y {
		return v
	}
	return ``
}

func (i Items) GetInt(key string) int {
	if v, y := i[key]; y {
		return get.ParseString(v, `int`).(int)
	}
	return 0
}

func (i Items) GetInt64(key string) int64 {
	if v, y := i[key]; y {
		return get.ParseString(v, `int64`).(int64)
	}
	return 0
}

func (i Items) GetFloat32(key string) float32 {
	if v, y := i[key]; y {
		return get.ParseString(v, `float32`).(float32)
	}
	return 0
}

func (i Items) GetFloat64(key string) float64 {
	if v, y := i[key]; y {
		return get.ParseString(v, `float64`).(float64)
	}
	return 0
}

func (i Items) GetBool(key string) bool {
	if v, y := i[key]; y {
		return get.ParseString(v, `bool`).(bool)
	}
	return false
}

type IItems map[string]interface{}

func (i IItems) Get(key string) interface{} {
	if v, y := i[key]; y {
		return v
	}
	return nil
}

func (i IItems) GetString(key string) string {
	if val, ok := i.Get(key).(string); ok {
		return val
	}
	return ``
}

func (i IItems) GetInt(key string) int {
	v := i.Get(key)
	if val, ok := v.(int); ok {
		return val
	} else if val, ok := v.(int64); ok {
		return int(val)
	}
	return 0
}

func (i IItems) GetInt64(key string) int64 {
	v := i.Get(key)
	if val, ok := v.(int64); ok {
		return val
	} else if val, ok := v.(int); ok {
		return int64(val)
	}
	return 0
}

func (i IItems) GetFloat32(key string) float32 {
	v := i.Get(key)
	if val, ok := v.(float32); ok {
		return val
	}
	if val, ok := v.(int64); ok {
		return float32(val)
	}
	return 0
}

func (i IItems) GetFloat64(key string) float64 {
	v := i.Get(key)
	if val, ok := v.(float64); ok {
		return val
	}
	if val, ok := v.(int64); ok {
		return float64(val)
	}
	return 0
}

func (i IItems) GetSlice(key string) []interface{} {
	v := i.Get(key)
	if val, ok := v.([]interface{}); ok {
		return val
	}
	return []interface{}{}
}

func (i IItems) GetStringSlice(key string) []string {
	v := i.GetSlice(key)
	if len(v) > 0 {
		val := []string{}
		for _, vv := range v {
			if vvv, ok := vv.(string); ok {
				val = append(val, vvv)
			}
		}
		return val
	}
	return []string{}
}

func (i IItems) GetBool(key string) bool {
	if val, ok := i.Get(key).(bool); ok {
		return val
	}
	return false
}
