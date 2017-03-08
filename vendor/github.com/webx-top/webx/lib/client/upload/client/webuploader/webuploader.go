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
package webuploader

import (
	"io"

	"github.com/webx-top/echo"
	uploadClient "github.com/webx-top/webx/lib/client/upload"
)

func init() {
	uploadClient.Reg(`webuploader`, func() uploadClient.Client {
		return New()
	})
}

func New() uploadClient.Client {
	return &Webuploader{}
}

type Webuploader struct {
	result *uploadClient.Result
	echo.Context
}

func (a *Webuploader) Init(ctx echo.Context, res *uploadClient.Result) {
	a.Context = ctx
	a.result = res
}

func (a *Webuploader) Name() string {
	return "file"
}

func (a *Webuploader) Body() (file io.ReadCloser, err error) {
	file, a.result.FileName, err = uploadClient.Receive(a.Name(), a.Context)
	if err != nil {
		return
	}
	return
}

func (a *Webuploader) Result(errMsg string) (r string) {
	cid := a.Context.Form("id")
	if cid == "" {
		form := a.Context.Request().MultipartForm()
		if form != nil && form.Value != nil {
			if v, ok := form.Value["id"]; ok && len(v) > 0 {
				cid = v[0]
			}
		}
	}
	if errMsg == "" {
		r = `{"jsonrpc":"2.0","result":{"id":"` + a.result.FileIdString() + `","containerid":"` + cid + `"},"error":null}`
		return
	}
	code := "100"
	r = `{"jsonrpc":"2.0","result":{"id":"` + a.result.FileIdString() + `","containerid":"` + cid + `"},"error":{"code":"` + code + `","message":"` + errMsg + `"}}`

	return
}
