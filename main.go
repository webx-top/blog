package main

import (
	"flag"

	_ "github.com/webx-top/blog/app/admin"
	_ "github.com/webx-top/blog/app/blog"

	"github.com/webx-top/blog/app/base"
)

func main() {
	port := flag.String("p", "5000", "port of your app.")
	flag.Parse()

	base.Server.Run("127.0.0.1", *port)
}
