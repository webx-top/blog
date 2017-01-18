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
	C "github.com/webx-top/captcha"
	"github.com/webx-top/echo"
	"github.com/webx-top/validation"
	X "github.com/webx-top/webx"
	"github.com/webx-top/webx/lib/database"
)

func NewController(c echo.Context) *Controller {
	a := &Controller{
		Controller: X.NewController(c),
		DB:         DB,
	}
	return a
}

type Controller struct {
	*X.Controller
	DB *database.Orm
}

func (a *Controller) Init(c echo.Context) error {
	return nil
}

func (a *Controller) NotFoundData() *X.Context {
	return a.SetErr(a.T(`数据不存在`))
}

func (a *Controller) NotModified() *X.Context {
	return a.SetErr(a.T(`没有修改任何内容`))
}

func (a *Controller) Failed() *X.Context {
	return a.SetErr(a.T(`操作失败`))
}

func (a *Controller) Done() *X.Context {
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

func (a *Controller) Valid() (v *validation.Validation) {
	v = &validation.Validation{
		SendError: func(e *validation.ValidationError) {
			a.SetErr(e.Message, e.Field)
		},
	}
	return
}

func (a *Controller) ValidOk(m interface{}, args ...string) bool {
	v := a.Valid()
	ok := true
	if m != nil {
		ok, _ = v.Valid(m, args...)
	}
	return ok
}
