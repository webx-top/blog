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
package base

import (
	X "github.com/webx-top/webx"
	C "github.com/webx-top/webx/lib/captcha"
	"github.com/webx-top/webx/lib/i18n"
)

func NewController(c *X.Context) *Controller {
	a := &Controller{
		Controller: X.NewController(c),
	}
	return a
}

type Controller struct {
	*X.Controller
}

func (a *Controller) Init(c *X.Context) error {
	return nil
}

func (a *Controller) Before() error {
	return a.Controller.Before()
}

func (a *Controller) After() error {
	return a.Controller.After()
}

func (a *Controller) Lang() string {
	if a.Language == `` {
		a.Language = DefaultLang
	}
	return a.Language
}

func (a *Controller) T(key string, args ...interface{}) string {
	return i18n.T(a.Lang(), key, args...)
}

func (a *Controller) NotFoundData() error {
	return a.SetErr(a.T(`数据不存在`))
}

func (a *Controller) NotModified() error {
	return a.SetErr(a.T(`没有修改任何内容`))
}

func (a *Controller) Failed() error {
	return a.SetErr(a.T(`操作失败`))
}

func (a *Controller) Done() error {
	return a.SetSuc(a.T(`操作成功`))
}

// 验证码验证
func (a *Controller) VerifyCaptcha(captchaSolution string) bool {
	captchaId := a.Form("captchaId")
	if !C.VerifyString(captchaId, captchaSolution) {
		return false
	}
	return true
}
