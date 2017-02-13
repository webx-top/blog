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
package simple

import (
	"io"

	X "github.com/webx-top/webx"
	uploadClient "github.com/webx-top/webx/lib/client/upload"
)

func init() {
	uploadClient.Reg(`simple`, func() uploadClient.Client {
		return New()
	})
}

func New() uploadClient.Client {
	return &Simple{}
}

type Simple struct {
	result *uploadClient.Result
	*X.Context
}

func (a *Simple) Init(ctx *X.Context, res *uploadClient.Result) {
	a.Context = ctx
	a.result = res
}

func (a *Simple) Name() string {
	return "filedata"
}

func (a *Simple) Body() (file io.ReadCloser, err error) {
	file, a.result.FileName, err = uploadClient.Receive(a.Name(), a.Context)
	if err != nil {
		return
	}
	return
}

func (a *Simple) Result(errMsg string) (r string) {
	status := "1"
	if errMsg != "" {
		status = "0"
	}
	r = `{"Status":` + status + `,"Message":"` + errMsg + `","Data":{"Url":"` + a.result.FileUrl + `","Id":"` + a.result.FileIdString() + `"}}`
	return
}
