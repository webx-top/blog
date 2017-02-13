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
package upload

import (
	"fmt"
	"io"

	X "github.com/webx-top/webx"
)

type Result struct {
	FileId   int64
	FileName string
	FileUrl  string
	FileType FileType
	Addon    interface{}
}

func (r *Result) FileIdString() string {
	return fmt.Sprintf(`%d`, r.FileId)
}

func New() Client {
	return &defaultClient{}
}

type defaultClient struct {
	result *Result
	*X.Context
}

func (a *defaultClient) Init(ctx *X.Context, res *Result) {
	a.Context = ctx
	a.result = res
}

func (a *defaultClient) Name() string {
	return "filedata"
}

func (a *defaultClient) Body() (file io.ReadCloser, err error) {
	file, a.result.FileName, err = Receive(a.Name(), a.Context)
	if err != nil {
		return
	}
	return
}

func (a *defaultClient) Result(errMsg string) (r string) {
	status := "1"
	if errMsg != "" {
		status = "0"
	}
	r = `{"Status":` + status + `,"Message":"` + errMsg + `","Data":{"Url":"` + a.result.FileUrl + `","Id":"` + a.result.FileIdString() + `"}}`
	return
}

type Client interface {
	//初始化
	Init(*X.Context, *Result)

	//file表单域name属性值
	Name() string

	//文件内容
	Body() (io.ReadCloser, error)

	//返回结果
	Result(string) string
}
