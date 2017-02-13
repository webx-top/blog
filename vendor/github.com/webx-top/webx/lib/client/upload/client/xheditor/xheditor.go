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
package xheditor

import (
	"io"
	"net/url"

	X "github.com/webx-top/webx"
	uploadClient "github.com/webx-top/webx/lib/client/upload"
)

func init() {
	uploadClient.Reg(`xheditor`, func() uploadClient.Client {
		return New()
	})
}

func New() uploadClient.Client {
	return &XhEditor{}
}

type XhEditor struct {
	result *uploadClient.Result
	*X.Context
}

func (a *XhEditor) Init(ctx *X.Context, res *uploadClient.Result) {
	a.Context = ctx
	a.result = res
}

func (a *XhEditor) Name() string {
	return "filedata"
}

func (a *XhEditor) Body() (file io.ReadCloser, err error) {
	file, a.result.FileName, err = uploadClient.Receive(a.Name(), a.Context)
	if err != nil {
		return
	}
	return
}

func (a *XhEditor) Result(errMsg string) (r string) {
	var msg, publicUrl string
	if a.Context.Form("immediate") == "1" {
		publicUrl = "!" + a.result.FileUrl
	} else {
		publicUrl = a.result.FileUrl
	}
	switch a.result.FileType {
	case uploadClient.TypeImage, "":
		msg = `{"url":"` + publicUrl + `||||` + url.QueryEscape(a.result.FileName) + `","localname":"` + a.result.FileName + `","id":"` + a.result.FileIdString() + `"}`
	case uploadClient.TypeFlash, uploadClient.TypeMedia, "file":
		fallthrough
	default:
		msg = `{"url":"` + publicUrl + `","id":"` + a.result.FileIdString() + `"}`
	}
	if msg == "" {
		msg = "{}"
	}
	r = `{"err":"` + errMsg + `","msg":` + msg + `}`
	return
}
