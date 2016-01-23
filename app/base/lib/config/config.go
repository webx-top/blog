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
