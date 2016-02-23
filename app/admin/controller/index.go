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
	//"github.com/webx-top/blog/app/admin/lib"
	X "github.com/webx-top/webx"
	//"github.com/webx-top/webx/lib/com"
)

type Index struct {
	index X.Mapper
	*Base
}

func (a *Index) Init(c *X.Context) error {
	a.Base = New(c)
	return nil
}

func (a *Index) Before() error {
	return a.Base.Before()
}

func (a *Index) Index() error {
	return a.Display()
}
