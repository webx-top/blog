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
	"strconv"
	"strings"

	"github.com/webx-top/blog/app/blog/lib"
	E "github.com/webx-top/echo"
	X "github.com/webx-top/webx"
	C "github.com/webx-top/webx/lib/captcha"
)

func init() {
	lib.App.Reg(&Captcha{}).Auto()
}

type Captcha struct {
	*Base
	show   X.Mapper `webx:":name"`
	reload X.Mapper `webx:".JSON"`
}

func (a *Captcha) Init(c *X.Context) error {
	a.Base = New(c)
	return nil
}

func (this *Captcha) Show() error {
	defer func() {
		this.Exit = true
	}()
	if _, ok := this.checkRefer(func() error { return nil }); !ok {
		return nil
	}
	param := this.P(0)
	var id, ext string
	if p := strings.LastIndex(param, `.`); p > 0 {
		id = param[0:p]
		ext = param[p:]
	}
	if ext == "" || id == "" {
		return this.NotFound()
	}
	if this.Query("reload") != "" {
		C.Reload(id)
	}
	header := this.Response().Header()
	download := this.Query("download") != ""
	if download {
		header.Set(E.ContentType, "application/octet-stream")
	}
	switch ext {
	case ".png":
		if !download {
			header.Set(E.ContentType, "image/png")
		}
		return C.WriteImage(this.Response().Writer(), id, C.StdWidth, C.StdHeight)
	case ".wav":
		lang := strings.ToLower(this.Query("lang"))
		if lang != `en` && lang != `ru` && lang != `zh` {
			lang = `zh`
		}
		au, err := C.GetAudio(id, lang)
		if err != nil {
			return err
		}
		if !download {
			header.Set(E.ContentType, "audio/x-wav")
		}
		header.Set("Content-Length", strconv.Itoa(au.EncodedLen()))
		_, err = au.WriteTo(this.Response().Writer())
		return err
	}
	return nil
}

func (this *Captcha) Reload() (err error) {
	err, _ = this.checkRefer(func() error {
		d := struct{ Id string }{Id: C.New()}
		this.Output.Data = d
		return nil
	})
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
		if domain := this.Server.App(app).Domain; domain != `` {
			allowedDomain = "http://" + domain
			this.Response().Header().Set("Access-Control-Allow-Origin", allowedDomain)
		}
	}
	//this.SetHeader("Access-Control-Allow-Origin", "*")
	r := this.Refer()
	logger := this.Server.Core.Logger()
	//println("[Refer]", r, this.IsAjax(), this.Request.Host)
	if r == "" || (strings.Contains(r, "://") && !this.checkAllowedDomain(r, allowedDomain) && !strings.Contains(r, "://"+this.Request().Host+"/")) {
		logger.Errorf("[IP:%s]Update captcha from [%s] => Denied!", this.IP(), r)
		err = this.NotFound()
	} else {
		err = f()
		logger.Infof("[IP:%s]Update captcha from [%s] => Allowed.", this.IP(), r)
		ret = true
	}
	return
}
