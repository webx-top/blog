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
package htmlcache

import (
	"fmt"
	"github.com/webx-top/echo"
	"net/http"
	"time"
)

type Rule struct {
	SaveFile   string                                              //保存名称
	SaveFunc   func(saveFile string, c echo.Context) string        //自定义保存名称
	ExpireTime int                                                 //过期时间(秒)
	ExpireFunc func(saveFile string, c echo.Context) (int64, bool) //判断缓存是否过期
}

func HttpCache(ctx echo.Context, eTag interface{}, etagValidator func(oldEtag, newEtag string) bool) bool {
	var etag string
	if eTag == nil {
		etag = fmt.Sprintf(`%v`, time.Now().UTC().Unix())
	} else {
		etag = fmt.Sprintf(`%v`, eTag)
	}
	resp := ctx.Response()
	//resp.Header().Set(`Connection`, `keep-alive`)
	resp.Header().Set(`X-Cache`, `HIT from Webx-Page-Cache`)
	if inm := ctx.Request().Header().Get("If-None-Match"); inm != `` {
		var valid bool
		if etagValidator != nil {
			valid = etagValidator(inm, etag)
		} else {
			valid = inm == etag
		}
		if valid {
			resp.Header().Del(`Content-Type`)
			resp.Header().Del(`Content-Length`)
			resp.WriteHeader(http.StatusNotModified)
			ctx.Object().Echo().Logger().Debugf(`%v is not modified.`, ctx.Path())
			return true
		}
	}
	resp.Header().Set(`Etag`, etag)
	resp.Header().Set(`Cache-Control`, `public,max-age=1`)
	return false
}
