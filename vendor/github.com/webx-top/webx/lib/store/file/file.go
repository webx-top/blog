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
package file

import (
	"io"
)

type Result struct {
	//网址路径
	Path string //UrlPath

	//文件保存路径
	Save string //SavePath

	//文件类型(图片:image)
	Type string

	//引擎名称
	Engine string
}

type Storer interface {
	//打开引擎(连接或认证等)
	Open() error

	//保存文件
	Put(io.ReadCloser, string) (*Result, error)

	//获取文件
	Get(string) (io.ReadCloser, error)

	//删除文件
	Del(string) error

	//关闭引擎
	Close() error

	//引擎名称
	EName() string
}
