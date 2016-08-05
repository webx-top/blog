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
	"github.com/webx-top/webx/lib/com"
	fileStore "github.com/webx-top/webx/lib/store/file"

	D "github.com/webx-top/blog/app/base/dbschema"
)

type PostCollection struct {
	Post     *D.Post     `xorm:"extends"`
	User     *D.User     `xorm:"extends"`
	Ocontent *D.Ocontent `xorm:"extends"`
}

type PostCollection2 struct {
	Post     *D.Post     `xorm:"extends" rel:""`
	User     *D.User     `xorm:"extends" rel:"LEFT:User.id=Post.uid"`
	Ocontent *D.Ocontent `xorm:"extends" rel:"LEFT:Ocontent.rc_id=Post.id AND Ocontent.rc_type='post'"`
}

type Post2 struct {
	*D.Post     `xorm:"extends"`
	*D.User     `xorm:"extends"`
	*D.Ocontent `xorm:"extends"`
	Id          int64
}

type Post3 struct {
	*D.Post     `xorm:"extends" alias:"a"`
	*D.User     `xorm:"extends" alias:"b"`
	*D.Ocontent `xorm:"extends" alias:"c"`
	//Id       int64
}

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
	m := &D.Post{}
	a.DB.Id(1).Get(m)
	ms := []*PostCollection{}
	a.Logger.Debug(`测试extends查询：`)
	a.DB.Where(`Post.id=1`).Join(`LEFT`, `webx_user`, `User.id=Post.uid`).
		Join(`LEFT`, `webx_ocontent`, `Ocontent.rc_id=Post.id AND Ocontent.rc_type=?`, `post`).Find(&ms)
	ms2 := []*Post2{}

	a.Logger.Debug(`测试extends带普通字段的查询：`)
	a.DB.Where(`Post.id=1`).Join(`LEFT`, `webx_user`, `User.id=Post.uid`).
		Join(`LEFT`, `webx_ocontent`, `Ocontent.rc_id=Post.id AND Ocontent.rc_type=?`, `post`).Find(&ms2)

	a.Logger.Debug(`测试忽略Post.content字段的查询：`)
	//自动加表前缀
	a.DB.Omit(`Post.content`).Where(`Post.id=1`).Join(`LEFT`, `~user`, `User.id=Post.uid`).
		Join(`LEFT`, `~ocontent`, `Ocontent.rc_id=Post.id AND Ocontent.rc_type=?`, `post`).Find(&ms2)

	a.Logger.Debug(`测试通过rel标签指定关联条件的查询：`)
	//使用rel标签
	ms3 := []*PostCollection2{}
	a.DB.Omit(`User.passwd`).Where(`Post.id=2`).Find(&ms3)
	com.Dump(ms3)

	a.Logger.Debug(`测试通过alias标签指定表别名的查询：`)
	msAlias := []*Post3{}
	a.DB.Where(`a.id=1`).Join(`LEFT`, `webx_user`, `a.id=b.uid`).
		Join(`LEFT`, `webx_ocontent`, `c.rc_id=a.id AND c.rc_type=?`, `post`).Find(&msAlias)
	return a.Display()
}

func (a *Index) Upload() error {
	client := a.Form(`client`)
	uc := uploadClient.Get(client)
	rs := &uploadClient.Result{}
	uc.Init(a.Context, rs)
	f, err := uc.Body()
	if err != nil {
		return err
	}
	defer f.Close()
	store := fileStore.Get("local")
	store.Open()
	defer store.Close()
	if r, err := store.Put(f, rs.FileName); err != nil {
		return err
	} else {
		rs.FileUrl = r.Path
	}
	res := uc.Result(``)
	return a.JSONBlob(200, []byte(res))
}
