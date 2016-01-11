package main

import (
	X "bitbucket.org/admpub/webx"
	_ "github.com/webx-top/blog/app/admin"
	_ "github.com/webx-top/blog/app/blog"
)

func main() {
	X.Serv().Run("127.0.0.1", "8080")
}
