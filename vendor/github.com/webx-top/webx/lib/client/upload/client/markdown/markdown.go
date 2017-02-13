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
package markdown

import (
	"io"

	"net/url"
	"time"

	X "github.com/webx-top/webx"
	uploadClient "github.com/webx-top/webx/lib/client/upload"
)

func init() {
	uploadClient.Reg(`markdown`, func() uploadClient.Client {
		return New()
	})
}

func New() uploadClient.Client {
	return &Markdown{}
}

type Markdown struct {
	result *uploadClient.Result
	*X.Context
}

func (a *Markdown) Init(ctx *X.Context, res *uploadClient.Result) {
	a.Context = ctx
	a.result = res
}

func (a *Markdown) Name() string {
	return "editormd-image-file"
}

func (a *Markdown) Body() (file io.ReadCloser, err error) {
	file, a.result.FileName, err = uploadClient.Receive(a.Name(), a.Context)
	if err != nil {
		return
	}
	return
}

func (a *Markdown) Result(errMsg string) (r string) {
	var succed string = "0" // 0 表示上传失败，1 表示上传成功
	if errMsg == "" {
		succed = "1"
	}
	callback := a.Context.Form(`callback`)
	dialogId := a.Context.Form(`dialog_id`)
	if callback != `` && dialogId != `` {
		//跨域上传返回操作
		nextUrl := callback + "?dialog_id=" + dialogId + "&temp=" + time.Now().String() + "&success=" + succed + "&message=" + url.QueryEscape(errMsg) + "&url=" + a.result.FileUrl
		a.Context.Redirect(nextUrl)
	} else {
		r = `{"success":` + succed + `,"message":"` + errMsg + `","url":"` + a.result.FileUrl + `","id":"` + a.result.FileIdString() + `"}`
	}
	return
}
