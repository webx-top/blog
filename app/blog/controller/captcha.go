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
	"strings"

	//"github.com/webx-top/blog/app/blog/lib"
	C "github.com/webx-top/captcha"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/handler/captcha"
	X "github.com/webx-top/webx"
)

type Captcha struct {
	*Base
	show   X.Mapper `webx:":name"`
	reload X.Mapper `webx:".JSON"`
}

func (a *Captcha) Init(c echo.Context) error {
	a.Base = New(c)
	return nil
}

func (this *Captcha) Show() error {
	defer func() {
		this.Exit()
	}()
	if _, ok := this.checkRefer(func() error { return nil }); !ok {
		return nil
	}
	return captcha.Captcha()(this.Context)
}

func (this *Captcha) Reload() (err error) {
	err, _ = this.checkRefer(func() error {
		d := struct{ Id string }{Id: C.New()}
		this.SetSucData(d)
		return nil
	})
	if err != nil {
		return err
	}
	return this.Display()
}

func (this *Captcha) checkAllowedDomain(refer string, allowedUrl string) bool {
	if allowedUrl == "" {
		return false
	}
	return strings.HasPrefix(refer, allowedUrl)
}

func (this *Captcha) checkRefer(f func() error) (err error, ret bool) {
	app := this.Query("app")
	allowedDomain := ""
	if app != "" {
		if domain := this.Server.Module(app).Domain; domain != `` {
			allowedDomain = "http://" + domain
			this.Response().Header().Set("Access-Control-Allow-Origin", allowedDomain)
		}
	}
	//this.SetHeader("Access-Control-Allow-Origin", "*")
	r := this.Refer()
	logger := this.Server.Core.Logger()
	//println("[Refer]", r, this.IsAjax(), this.Request().Host())
	if r == "" || (strings.Contains(r, "://") && !this.checkAllowedDomain(r, allowedDomain) && !strings.Contains(r, "://"+this.Request().Host()+"/") && !strings.Contains(r, "://"+this.Request().Host()+":")) {
		logger.Errorf("[IP:%s]Update captcha from [%s] => Denied!", this.IP(), r)
		err = this.NotFound()
	} else {
		err = f()
		logger.Infof("[IP:%s]Update captcha from [%s] => Allowed.", this.IP(), r)
		ret = true
	}
	return
}
