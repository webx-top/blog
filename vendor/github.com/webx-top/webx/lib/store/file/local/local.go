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
package local

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	up "github.com/webx-top/webx/lib/client/upload"
	. "github.com/webx-top/webx/lib/store/file"
	"github.com/webx-top/webx/lib/uuid"
)

func New(conf map[string]string) *Local {
	//return &LocalFs{SavePath: "./coscms/app/base/static/upload/", PublicUrl: "/upload/"}
	return &Local{
		SavePath:  conf["SavePath"],
		PublicUrl: conf["PublicUrl"],
		RootPath:  conf["RootPath"] + "/",
		Engine:    "local",
	}
}

func RegStore(s *Local) {
	Reg(s.Engine, func() Storer {
		return s
	})
}

type Local struct {
	SavePath  string
	PublicUrl string
	RootPath  string
	Engine    string
}

func (this *Local) Open() error {
	return nil
}

func (this *Local) Close() error {
	return nil
}

//生成文件保存位置和网址
func (this *Local) New(obj *Result, fileName string) (savePath string, publicUrl string) {
	now := time.Now()
	year, month, day := now.Date()

	fileExtn := filepath.Ext(fileName)
	fileExtn = strings.ToLower(fileExtn)
	dirName := up.GetType(fileExtn)
	obj.Type = dirName
	if len(dirName) > 0 {
		dirName += "/"
	}
	savePath = fmt.Sprintf("%s%s%d/%d/%d/", this.SavePath, dirName, year, month, day)
	if err := os.MkdirAll(savePath, 0777); err != nil {
		fmt.Println(err)
	}
	//strings.Replace(uuid.NewUUID().String(), "-", "", -1)
	fileName = uuid.NewUUID().String() + fileExtn
	savePath += fileName //文件存储路径

	publicUrl = fmt.Sprintf("%s%s%d/%d/%d/%s", this.PublicUrl, dirName, year, month, day, fileName) //文件网址
	return
}

func (this *Local) Put(body io.ReadCloser, fileName string) (obj *Result, err error) {
	obj = &Result{Engine: this.Engine}
	savePath, publicUrl := this.New(obj, fileName)
	var f *os.File
	f, err = os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return
	}
	defer f.Close()
	_, err = io.Copy(f, body)
	if err != nil {
		return
	}
	obj.Save = strings.TrimPrefix(savePath, this.RootPath)
	obj.Path = publicUrl
	return
}

//获取文件内容
func (this *Local) Get(key string) (io.ReadCloser, error) {
	return os.Open(this.RootPath + key)
}

//删除文件
func (this *Local) Del(key string) error {
	return os.Remove(this.RootPath + key)
}

func (this *Local) EName() string {
	return this.Engine
}
