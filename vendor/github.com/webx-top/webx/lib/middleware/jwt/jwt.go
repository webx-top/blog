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
package jwt

import (
	"errors"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	jwtRequest "github.com/dgrijalva/jwt-go/request"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
)

func New(secret string) *JWT {
	return &JWT{
		Secret: secret,
		CondFn: func(c echo.Context) bool {
			ignore, _ := c.Get(`webx:ignoreJwt`).(bool)
			return !ignore
		},
		HeaderKey: "Authorization",
		URLKey:    "access_token",
	}
}

type JWT struct {
	Secret    string
	CondFn    func(echo.Context) bool
	HeaderKey string
	URLKey    string
}

func (j *JWT) Validate() echo.MiddlewareFunc {
	return echo.MiddlewareFunc(func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			if j.CondFn != nil && j.CondFn(c) == false {
				return h.Handle(c)
			}
			/*//Test
			tokenString, err := j.Response(map[string]interface{}{"uid": "1", "username": "admin"})
			if err == nil {
				println("jwt token:", tokenString)
			}
			//*/
			token, err := j.ParseFromRequest(c.Request(), func(token *jwt.Token) (interface{}, error) {
				b := ([]byte(j.Secret))
				return b, nil
			})
			if err != nil {
				return err
			}
			if !token.Valid {
				return errors.New(`Incorrect signature.`)
			}
			c.Set(`webx:jwtClaims`, token.Claims.(jwt.MapClaims))
			return h.Handle(c)
		})
	})
}

func (j *JWT) Claims(c echo.Context) map[string]interface{} {
	r, _ := c.Get(`webx:jwtClaims`).(map[string]interface{})
	return r
}

func (j *JWT) Ignore(on bool, c echo.Context) {
	c.Set(`webx:ignoreJwt`, on)
}

/*
本函数所生成结果的用法
用法一：写入header头的属性“Authorization”中，值设为：前缀BEARER加tokenString的值
用法二：发送post或get参数“access_token”，值设为：tokenString的值
*/
func (j *JWT) Response(values map[string]interface{}) (tokenString string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	// Headers
	token.Header["alg"] = "HS256"
	token.Header["typ"] = "JWT"
	// Claims
	values["exp"] = time.Now().Add(time.Hour * 72).Unix()
	token.Claims = jwt.MapClaims(values)
	tokenString, err = token.SignedString([]byte(j.Secret))
	return
}

func (j *JWT) ParseFromRequest(req engine.Request, keyFunc jwt.Keyfunc) (token *jwt.Token, err error) {

	// Look for an Authorization header
	if ah := req.Header().Get(j.HeaderKey); ah != "" {
		// Should be a bearer token
		if len(ah) > 6 && strings.ToUpper(ah[0:6]) == "BEARER" {
			return jwt.Parse(ah[7:], keyFunc)
		}
	}

	// Look for "access_token" parameter
	if tokStr := req.FormValue(j.URLKey); tokStr != "" {
		return jwt.Parse(tokStr, keyFunc)
	}

	return nil, jwtRequest.ErrNoTokenInRequest

}
