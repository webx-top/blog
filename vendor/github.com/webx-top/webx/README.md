# webx
webx 是一个golang版本的web通用开发框架。

webx 基于***echo框架***，是对 ***echo框架***(双引擎可切换加强版：https://github.com/webx-top/echo) 的再次封装，具备了开发一个网站所需要的所有组件。

# 特性
1. 同时支持函数式和面向对象式编程；
2. 支持多模块架构；
3. 支持多语言；
4. 支持响应多种格式(JSON/HTML/XML)；
5. 完善的[模板引擎](https://github.com/webx-top/webx/tree/master/lib/tplex)(也可以选择pongo2模板引擎)
6. 支持多种缓存引擎(memcache/redis/boltDB/levelDB ...)
7. 文件支持多种存储方式(本地/FTP/其它云存储 ...)
8. 支持多种前端组件上传文件(xhEditor/webuploader/Editor.md/其它 ...)
9. 更多...略

## 1. 函数式编程

```Go
package main

import "github.com/webx-top/webx"

func main(){

    webx.R("/ping",func(c *webx.Context) error {
	    return c.String(200, "pong")
    })

    webx.Run(":8080")

}
```


## 2. 面向对象式编程

```Go
package main

import (
	"fmt"

	"github.com/webx-top/webx"
)

// ==Index===================================
type Index struct {
	index webx.Mapper
	*webx.Controller
}

func (h *Index) Init(c *webx.Context) error {
	h.Controller = webx.NewController(c)
	return nil
}

// URL: /index/index or /
func (h *Index) Index() error {
	return h.String(200, "Hello world.[Controller:Index]")
}

// ==Home===================================
type Home struct {
	ping  webx.Mapper
	index webx.Mapper
	*webx.Controller
}

func (h *Home) Init(c *webx.Context) error {
	h.Controller = webx.NewController(c)
	return nil
}

// URL: /home/ or /home/index
func (h *Home) Index() error {
	return h.String(200, "Hello world.[Controller:Home]")
}

// URL: /home/ping
func (h *Home) Ping() error {
	return h.String(200, "pong")
}

// 前置操作(可选)
func (h *Home) Before() error {
	fmt.Println(`Before.`)
	return nil
}

// 后置操作(可选)
func (h *Home) After() error {
	fmt.Println(`After.`)
	return nil
}

func main() {

	webx.Use(&Index{}, &Home{})

	webx.Run(":8080")
}
```

# 样例工程

[博客系统](https://github.com/webx-top/blog) -- 欢迎参与 :grinning: