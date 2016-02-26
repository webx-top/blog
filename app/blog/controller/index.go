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
package controller

import (
	//"github.com/webx-top/blog/app/base"
	//"github.com/webx-top/blog/app/blog/lib"
	X "github.com/webx-top/webx"
	uploadClient "github.com/webx-top/webx/lib/client/upload"
	_ "github.com/webx-top/webx/lib/client/upload/client/markdown"
	_ "github.com/webx-top/webx/lib/client/upload/client/simple"
	_ "github.com/webx-top/webx/lib/client/upload/client/webuploader"
	_ "github.com/webx-top/webx/lib/client/upload/client/xheditor"
	//"github.com/webx-top/webx/lib/com"
	fileStore "github.com/webx-top/webx/lib/store/file"
)

type Index struct {
	index  X.Mapper
	upload X.Mapper `.JSON`
	*Base
}

func (a *Index) Init(c *X.Context) error {
	a.Base = New(c)
	return nil
}

func (a *Index) Index() error {
	a.AssignX(&map[string]interface{}{
		`name`: `webx`,
		"test": "times---",
		"r":    []string{"one", "two", "three"},
	})
	return a.Display()
}

func (a *Index) Upload() error {
	client := a.Form(`client`)
	uc := uploadClient.Get(client)
	uc.Init(a.Context)
	rs := &uploadClient.Result{}
	f, err := uc.Body(rs)
	if err != nil {
		return err
	}
	defer f.Close()
	store := fileStore.Get("local")
	if r, err := store.Put(f, rs.FileName); err != nil {
		return err
	} else {
		rs.FileUrl = r.Path
	}
	res := uc.Result(``)
	return a.JSONBlob(200, []byte(res))
}
