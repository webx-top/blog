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
