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
package tplfunc

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/webx-top/captcha"
	"github.com/webx-top/com"
)

func New() (r template.FuncMap) {
	r = template.FuncMap{}
	for name, function := range TplFuncMap {
		r[name] = function
	}
	return
}

var TplFuncMap template.FuncMap = template.FuncMap{
	// ======================
	// time
	// ======================
	"Now":             Now,
	"ElapsedMemory":   com.ElapsedMemory, //内存消耗
	"TotalRunTime":    com.TotalRunTime,  //运行时长(从启动服务时算起)
	"CaptchaForm":     CaptchaForm,       //验证码图片
	"FormatByte":      com.FormatByte,    //字节转为适合理解的格式
	"FriendlyTime":    FriendlyTime,
	"FormatPastTime":  com.FormatPastTime, //以前距离现在多长时间
	"DateFormat":      com.DateFormat,
	"DateFormatShort": com.DateFormatShort,
	"Ts2time":         TsToTime, // 时间戳数字转time.Time
	"Ts2date":         TsToDate, // 时间戳数字转日期字符串

	// ======================
	// compare
	// ======================
	"Eq":       Eq,
	"Add":      Add,
	"Sub":      Sub,
	"Div":      Div,
	"Mul":      Mul,
	"IsNil":    IsNil,
	"IsEmpty":  IsEmpty,
	"NotEmpty": NotEmpty,
	"IsNaN":    IsNaN,
	"IsInf":    IsInf,

	// ======================
	// conversion type
	// ======================
	"Html":         ToHTML,
	"Js":           ToJS,
	"Css":          ToCSS,
	"ToJS":         ToJS,
	"ToCSS":        ToCSS,
	"ToURL":        ToURL,
	"ToHTML":       ToHTML,
	"ToHTMLAttr":   ToHTMLAttr,
	"ToHTMLAttrs":  ToHTMLAttrs,
	"ToStrSlice":   ToStrSlice,
	"Str":          com.Str,
	"Int":          com.Int,
	"Int32":        com.Int32,
	"Int64":        com.Int64,
	"Uint":         com.Uint,
	"Uint32":       com.Uint32,
	"Uint64":       com.Uint64,
	"Float32":      com.Float32,
	"Float64":      com.Float64,
	"ToFloat64":    ToFloat64,
	"ToFixed":      ToFixed,
	"Math":         Math,
	"NumberFormat": NumberFormat,

	// ======================
	// string
	// ======================
	"Contains":  strings.Contains,
	"HasPrefix": strings.HasPrefix,
	"HasSuffix": strings.HasSuffix,

	"ToLower":        strings.ToLower,
	"ToUpper":        strings.ToUpper,
	"LowerCaseFirst": com.LowerCaseFirst,
	"CamelCase":      com.CamelCase,
	"PascalCase":     com.PascalCase,
	"SnakeCase":      com.SnakeCase,
	"Reverse":        com.Reverse,
	"Ext":            filepath.Ext,
	"InExt":          InExt,

	"Concat":    Concat,
	"Replace":   strings.Replace, //strings.Replace(s, old, new, n)
	"Split":     strings.Split,
	"Join":      strings.Join,
	"Substr":    com.Substr,
	"StripTags": com.StripTags,
	"Nl2br":     NlToBr, // \n替换为<br>
	"AddSuffix": AddSuffix,

	// ======================
	// encode & decode
	// ======================
	"JsonEncode":   JsonEncode,
	"UrlEncode":    com.UrlEncode,
	"UrlDecode":    com.UrlDecode,
	"Base64Encode": com.Base64Encode,
	"Base64Decode": com.Base64Decode,

	// ======================
	// map & slice
	// ======================
	"InSlice":        com.InSlice,
	"InSlicex":       com.InSliceIface,
	"Set":            Set,
	"Append":         Append,
	"InStrSlice":     InStrSlice,
	"SearchStrSlice": SearchStrSlice,
	"URLValues":      URLValues,
	"ToSlice":        ToSlice,

	// ======================
	// regexp
	// ======================
	"Regexp":      regexp.MustCompile,
	"RegexpPOSIX": regexp.MustCompilePOSIX,

	// ======================
	// other
	// ======================
	"Ignore":  Ignore,
	"Default": Default,
}

func JsonEncode(s interface{}) string {
	r, _ := com.SetJSON(s)
	return r
}

func Ignore(_ interface{}) interface{} {
	return nil
}

func URLValues(values ...interface{}) url.Values {
	v := url.Values{}
	var k string
	for i, j := 0, len(values); i < j; i++ {
		if i%2 == 0 {
			k = fmt.Sprint(values[i])
			continue
		}
		v.Add(k, fmt.Sprint(values[i]))
		k = ``
	}
	if len(k) > 0 {
		v.Add(k, ``)
		k = ``
	}
	return v
}

func ToStrSlice(s ...string) []string {
	return s
}

func ToSlice(s ...interface{}) []interface{} {
	return s
}

func Concat(s ...string) string {
	return strings.Join(s, ``)
}

func InExt(fileName string, exts ...string) bool {
	ext := filepath.Ext(fileName)
	ext = strings.ToLower(ext)
	for _, _ext := range exts {
		if ext == strings.ToLower(_ext) {
			return true
		}
	}
	return false
}

func Default(defaultV interface{}, v interface{}) interface{} {
	switch val := v.(type) {
	case nil:
		return defaultV
	case string:
		if len(val) == 0 {
			return defaultV
		}
	case uint8, int8, uint, int, uint32, int32, int64, uint64:
		if val == 0 {
			return defaultV
		}
	case float32, float64:
		if val == 0.0 {
			return defaultV
		}
	default:
		if len(com.Str(v)) == 0 {
			return defaultV
		}
	}
	return v
}

func Set(renderArgs map[string]interface{}, key string, value interface{}) string {
	renderArgs[key] = value
	return ``
}

func Append(renderArgs map[string]interface{}, key string, value interface{}) string {
	if renderArgs[key] == nil {
		renderArgs[key] = []interface{}{value}
	} else {
		renderArgs[key] = append(renderArgs[key].([]interface{}), value)
	}
	return ``
}

//NlToBr Replaces newlines with <br />
func NlToBr(text string) template.HTML {
	return template.HTML(Nl2br(text))
}

//CaptchaForm 验证码表单域
func CaptchaForm(args ...string) template.HTML {
	id := "captcha"
	format := `<img id="%[2]sImage" src="/captcha/%[1]s.png" alt="Captcha image" onclick="this.src=this.src.split('?')[0]+'?reload='+Math.random();" /><input type="hidden" name="captchaId" id="%[2]sId" value="%[1]s" />`
	switch len(args) {
	case 2:
		format = args[1]
		fallthrough
	case 1:
		id = args[0]
	}
	cid := captcha.New()
	return template.HTML(fmt.Sprintf(format, cid, id))
}

//CaptchaVerify 验证码验证
func CaptchaVerify(captchaSolution string, idGet func(string) string) bool {
	//id := r.FormValue("captchaId")
	id := idGet("captchaId")
	if !captcha.VerifyString(id, captchaSolution) {
		return false
	}
	return true
}

//Nl2br 将换行符替换为<br />
func Nl2br(text string) string {
	return com.Nl2br(template.HTMLEscapeString(text))
}

func IsNil(a interface{}) bool {
	switch a.(type) {
	case nil:
		return true
	}
	return false
}

func interface2Int64(value interface{}) (int64, bool) {
	switch v := value.(type) {
	case uint:
		return int64(v), true
	case uint8:
		return int64(v), true
	case uint16:
		return int64(v), true
	case uint32:
		return int64(v), true
	case uint64:
		return int64(v), true
	case int:
		return int64(v), true
	case int8:
		return int64(v), true
	case int16:
		return int64(v), true
	case int32:
		return int64(v), true
	case int64:
		return v, true
	default:
		return 0, false
	}
}

func interface2Float64(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float32:
		return float64(v), true
	case float64:
		return v, true
	default:
		return 0, false
	}
}

func ToFloat64(value interface{}) float64 {
	if v, ok := interface2Int64(value); ok {
		return float64(v)
	}
	if v, ok := interface2Float64(value); ok {
		return v
	}
	return com.Float64(value)
}

func Add(left interface{}, right interface{}) interface{} {
	var rleft, rright int64
	var fleft, fright float64
	var isInt bool
	rleft, isInt = interface2Int64(left)
	if !isInt {
		fleft, _ = interface2Float64(left)
	}
	rright, isInt = interface2Int64(right)
	if !isInt {
		fright, _ = interface2Float64(right)
	}
	intSum := rleft + rright

	if isInt {
		return intSum
	}
	return fleft + fright + float64(intSum)
}

func Div(left interface{}, right interface{}) interface{} {
	return ToFloat64(left) / ToFloat64(right)
}

func Mul(left interface{}, right interface{}) interface{} {
	return ToFloat64(left) * ToFloat64(right)
}

func Math(op string, args ...interface{}) interface{} {
	length := len(args)
	if length < 1 {
		return float64(0)
	}
	switch op {
	case `mod`: //模
		if length < 2 {
			return float64(0)
		}
		return math.Mod(ToFloat64(args[0]), ToFloat64(args[1]))
	case `abs`:
		return math.Abs(ToFloat64(args[0]))
	case `acos`:
		return math.Acos(ToFloat64(args[0]))
	case `acosh`:
		return math.Acosh(ToFloat64(args[0]))
	case `asin`:
		return math.Asin(ToFloat64(args[0]))
	case `asinh`:
		return math.Asinh(ToFloat64(args[0]))
	case `atan`:
		return math.Atan(ToFloat64(args[0]))
	case `atan2`:
		if length < 2 {
			return float64(0)
		}
		return math.Atan2(ToFloat64(args[0]), ToFloat64(args[1]))
	case `atanh`:
		return math.Atanh(ToFloat64(args[0]))
	case `cbrt`:
		return math.Cbrt(ToFloat64(args[0]))
	case `ceil`:
		return math.Ceil(ToFloat64(args[0]))
	case `copysign`:
		if length < 2 {
			return float64(0)
		}
		return math.Copysign(ToFloat64(args[0]), ToFloat64(args[1]))
	case `cos`:
		return math.Cos(ToFloat64(args[0]))
	case `cosh`:
		return math.Cosh(ToFloat64(args[0]))
	case `dim`:
		if length < 2 {
			return float64(0)
		}
		return math.Dim(ToFloat64(args[0]), ToFloat64(args[1]))
	case `erf`:
		return math.Erf(ToFloat64(args[0]))
	case `erfc`:
		return math.Erfc(ToFloat64(args[0]))
	case `exp`:
		return math.Exp(ToFloat64(args[0]))
	case `exp2`:
		return math.Exp2(ToFloat64(args[0]))
	case `floor`:
		return math.Floor(ToFloat64(args[0]))
	case `max`:
		if length < 2 {
			return float64(0)
		}
		return math.Max(ToFloat64(args[0]), ToFloat64(args[1]))
	case `min`:
		if length < 2 {
			return float64(0)
		}
		return math.Min(ToFloat64(args[0]), ToFloat64(args[1]))
	case `pow`: //幂
		if length < 2 {
			return float64(0)
		}
		return math.Pow(ToFloat64(args[0]), ToFloat64(args[1]))
	case `sqrt`: //平方根
		return math.Sqrt(ToFloat64(args[0]))
	case `sin`:
		return math.Sin(ToFloat64(args[0]))
	case `log`:
		return math.Log(ToFloat64(args[0]))
	case `log2`:
		return math.Log2(ToFloat64(args[0]))
	case `log10`:
		return math.Log10(ToFloat64(args[0]))
	case `tan`:
		return math.Tan(ToFloat64(args[0]))
	case `tanh`:
		return math.Tanh(ToFloat64(args[0]))
	case `add`: //加
		if length < 2 {
			return float64(0)
		}
		return Add(ToFloat64(args[0]), ToFloat64(args[1]))
	case `sub`: //减
		if length < 2 {
			return float64(0)
		}
		return Sub(ToFloat64(args[0]), ToFloat64(args[1]))
	case `mul`: //乘
		if length < 2 {
			return float64(0)
		}
		return Mul(ToFloat64(args[0]), ToFloat64(args[1]))
	case `div`: //除
		if length < 2 {
			return float64(0)
		}
		return Div(ToFloat64(args[0]), ToFloat64(args[1]))
	}
	return nil
}

func IsNaN(v interface{}) bool {
	return math.IsNaN(ToFloat64(v))
}

func IsInf(v interface{}, s interface{}) bool {
	return math.IsInf(ToFloat64(v), com.Int(s))
}

func Sub(left interface{}, right interface{}) interface{} {
	var rleft, rright int64
	var fleft, fright float64
	var isInt bool
	rleft, isInt = interface2Int64(left)
	if !isInt {
		fleft, _ = interface2Float64(left)
	}
	rright, isInt = interface2Int64(right)
	if !isInt {
		fright, _ = interface2Float64(right)
	}
	if isInt {
		return rleft - rright
	}
	return fleft + float64(rleft) - (fright + float64(rright))
}

func ToFixed(value interface{}, precision interface{}) string {
	return fmt.Sprintf("%.*f", com.Int(precision), ToFloat64(value))
}

func Now() time.Time {
	return time.Now()
}

func Eq(left interface{}, right interface{}) bool {
	leftIsNil := (left == nil)
	rightIsNil := (right == nil)
	if leftIsNil || rightIsNil {
		if leftIsNil && rightIsNil {
			return true
		}
		return false
	}
	return fmt.Sprintf("%v", left) == fmt.Sprintf("%v", right)
}

func ToHTML(raw string) template.HTML {
	return template.HTML(raw)
}

func ToHTMLAttr(raw string) template.HTMLAttr {
	return template.HTMLAttr(raw)
}

func ToHTMLAttrs(raw map[string]interface{}) (r map[template.HTMLAttr]interface{}) {
	r = make(map[template.HTMLAttr]interface{})
	for k, v := range raw {
		r[ToHTMLAttr(k)] = v
	}
	return
}

func ToJS(raw string) template.JS {
	return template.JS(raw)
}

func ToCSS(raw string) template.CSS {
	return template.CSS(raw)
}

func ToURL(raw string) template.URL {
	return template.URL(raw)
}

func AddSuffix(s string, suffix string, args ...string) string {
	beforeChar := `.`
	if len(args) > 0 {
		beforeChar = args[0]
		if beforeChar == `` {
			return s + suffix
		}
	}
	p := strings.LastIndex(s, beforeChar)
	if p < 0 {
		return s
	}
	return s[0:p] + suffix + s[p:]
}

func IsEmpty(a interface{}) bool {
	switch v := a.(type) {
	case nil:
		return true
	case string:
		return len(v) == 0
	case []interface{}:
		return len(v) < 1
	default:
		switch fmt.Sprintf(`%v`, a) {
		case `<nil>`, ``, `[]`:
			return true
		}
	}
	return false
}

func NotEmpty(a interface{}) bool {
	return !IsEmpty(a)
}

func InStrSlice(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

func SearchStrSlice(values []string, value string) int {
	for i, v := range values {
		if v == value {
			return i
		}
	}
	return -1
}

func FriendlyTime(t interface{}, args ...string) string {
	var td time.Duration
	switch v := t.(type) {
	case time.Duration:
		td = v
	case int64:
		td = time.Duration(v)
	case int:
		td = time.Duration(v)
	case uint:
		td = time.Duration(v)
	case int32:
		td = time.Duration(v)
	case uint32:
		td = time.Duration(v)
	case uint64:
		td = time.Duration(v)
	default:
		td = time.Duration(com.Int64(t))
	}
	return com.FriendlyTime(td, args...)
}

func TsToTime(timestamp interface{}) time.Time {
	return TimestampToTime(timestamp)
}

func TsToDate(format string, timestamp interface{}) string {
	t := TimestampToTime(timestamp)
	if t.IsZero() {
		return ``
	}
	return t.Format(format)
}

func TimestampToTime(timestamp interface{}) time.Time {
	var ts int64
	switch v := timestamp.(type) {
	case int64:
		ts = v
	case uint:
		ts = int64(v)
	case int:
		ts = int64(v)
	case uint32:
		ts = int64(v)
	case int32:
		ts = int64(v)
	case uint64:
		ts = int64(v)
	default:
		i, e := strconv.ParseInt(fmt.Sprint(timestamp), 10, 64)
		if e != nil {
			log.Println(e)
		}
		ts = i
	}
	return time.Unix(ts, 0)
}

func NumberFormat(number interface{}, precision int, delim ...string) string {
	r := fmt.Sprintf(`%.*f`, precision, ToFloat64(number))
	d := `,`
	if len(delim) > 0 {
		d = delim[0]
		if len(d) == 0 {
			return r
		}
	}
	i := len(r) - 1 - precision
	j := int(math.Ceil(float64(i) / float64(3)))
	s := make([]string, j)
	v := r[i:]
	for i > 0 && j > 0 {
		j--
		start := i - 3
		if start < 0 {
			start = 0
		}
		s[j] = r[start:i]
		i = start
	}
	return strings.Join(s, d) + v
}
