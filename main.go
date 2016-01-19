package main

import (
	_ "github.com/webx-top/blog/app/admin"
	_ "github.com/webx-top/blog/app/blog"

	"github.com/webx-top/blog/app/base"
)

func main() {
	base.Server.Run("127.0.0.1", "5000")
}
