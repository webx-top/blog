package main

import (
	_ "github.com/webx-top/blog/app/admin"
	_ "github.com/webx-top/blog/app/blog"
	X "github.com/webx-top/webx"
)

func main() {
	X.Serv().Run("127.0.0.1", "8080")
}
