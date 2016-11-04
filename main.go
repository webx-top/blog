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
package main

import (
	"flag"

	_ "github.com/webx-top/blog/app/admin"
	_ "github.com/webx-top/blog/app/blog"

	"github.com/admpub/overseer"
	"github.com/webx-top/blog/app/base"
	"github.com/webx-top/echo/engine"
)

var Version = `1.0.0`

func main() {
	base.Version = Version
	port := flag.String("p", "8080", "port of your blog.")
	ngin := flag.String("e", "", "http engine")
	flag.Parse()

	addr := ":" + *port

	overseer.Run(overseer.Config{
		Program: func(state overseer.State) {
			conf := &engine.Config{
				Address:  addr,
				Listener: state.Listener,
			}
			base.Server.Run(conf, *ngin)
		},
		Address: addr,
		Debug:   true,
	})
}
