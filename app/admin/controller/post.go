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
	//"fmt"
	//"strings"

	"github.com/webx-top/blog/app/admin/lib"
	D "github.com/webx-top/blog/app/base/dbschema"
	"github.com/webx-top/blog/app/base/model"
	X "github.com/webx-top/webx"
	"github.com/webx-top/webx/lib/com"
)

func init() {
	lib.App.Reg(&Post{}).Auto()
}

type Post struct {
	index  X.Mapper
	edit   X.Mapper
	delete X.Mapper
	view   X.Mapper
	*Base
	postM *model.Post
}

func (a *Post) Init(c *X.Context) error {
	a.Base = New(c)
	a.postM = model.NewPost(c)
	return nil
}

func (a *Post) Index() error {
	if a.Format != `html` {
		sel := a.postM.NewSelect(&D.Post{})
		sel.Condition = `uid=?`
		sel.AddP(a.User.Id).FromClient(true, "title")
		countFn, data, _ := a.postM.List(sel)
		sel.Client.SetCount(countFn).Data(data)
	}
	return nil
}

func (a *Post) Edit() error {
	id := com.Int(a.Form(`id`))
	m, err := a.postM.Get(id)
	if err != nil {
		return err
	}
	a.Assign(`Detail`, m)
	return nil
}

func (a *Post) Delete() error {
	return nil
}

func (a *Post) View() error {
	return nil
}
