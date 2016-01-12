package tplfunc

import (
	"fmt"
	"html/template"
	"time"
)

var TplFuncMap template.FuncMap = template.FuncMap{
	"Now":         Now,
	"Eq":          Eq,
	"Add":         Add,
	"Sub":         Sub,
	"IsNil":       IsNil,
	"Html":        Html,
	"Js":          Js,
	"Css":         Css,
	"HtmlAttr":    HtmlAttr,
	"ToHtmlAttrs": ToHtmlAttrs,
}

func IsNil(a interface{}) bool {
	switch a.(type) {
	case nil:
		return true
	}
	return false
}

func Add(left interface{}, right interface{}) interface{} {
	var rleft, rright int64
	var fleft, fright float64
	var isInt bool = true
	switch left.(type) {
	case int:
		rleft = int64(left.(int))
	case int8:
		rleft = int64(left.(int8))
	case int16:
		rleft = int64(left.(int16))
	case int32:
		rleft = int64(left.(int32))
	case int64:
		rleft = left.(int64)
	case float32:
		fleft = float64(left.(float32))
		isInt = false
	case float64:
		fleft = left.(float64)
		isInt = false
	}

	switch right.(type) {
	case int:
		rright = int64(right.(int))
	case int8:
		rright = int64(right.(int8))
	case int16:
		rright = int64(right.(int16))
	case int32:
		rright = int64(right.(int32))
	case int64:
		rright = right.(int64)
	case float32:
		fright = float64(left.(float32))
		isInt = false
	case float64:
		fleft = left.(float64)
		isInt = false
	}

	var intSum int64 = rleft + rright

	if isInt {
		return intSum
	} else {
		return fleft + fright + float64(intSum)
	}
}

func Sub(left interface{}, right interface{}) interface{} {
	var rleft, rright int64
	var fleft, fright float64
	var isInt bool = true
	switch left.(type) {
	case int:
		rleft = int64(left.(int))
	case int8:
		rleft = int64(left.(int8))
	case int16:
		rleft = int64(left.(int16))
	case int32:
		rleft = int64(left.(int32))
	case int64:
		rleft = left.(int64)
	case float32:
		fleft = float64(left.(float32))
		isInt = false
	case float64:
		fleft = left.(float64)
		isInt = false
	}

	switch right.(type) {
	case int:
		rright = int64(right.(int))
	case int8:
		rright = int64(right.(int8))
	case int16:
		rright = int64(right.(int16))
	case int32:
		rright = int64(right.(int32))
	case int64:
		rright = right.(int64)
	case float32:
		fright = float64(left.(float32))
		isInt = false
	case float64:
		fleft = left.(float64)
		isInt = false
	}

	if isInt {
		return rleft - rright
	} else {
		return fleft + float64(rleft) - (fright + float64(rright))
	}
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

func Html(raw string) template.HTML {
	return template.HTML(raw)
}

func HtmlAttr(raw string) template.HTMLAttr {
	return template.HTMLAttr(raw)
}

func ToHtmlAttrs(raw map[string]interface{}) (r map[template.HTMLAttr]interface{}) {
	r = make(map[template.HTMLAttr]interface{})
	for k, v := range raw {
		r[HtmlAttr(k)] = v
	}
	return
}

func Js(raw string) template.JS {
	return template.JS(raw)
}

func Css(raw string) template.CSS {
	return template.CSS(raw)
}

func NewStatic(staticPath string) *Static {
	return &Static{Path: staticPath}
}

type Static struct {
	Path string
}

func (s *Static) StaticUrl(staticFile string) (r string) {
	r = s.Path + "/" + staticFile
	return
}

func (s *Static) JsUrl(staticFile string) (r string) {
	r = s.StaticUrl("js/" + staticFile)
	return
}

func (s *Static) CssUrl(staticFile string) (r string) {
	r = s.StaticUrl("css/" + staticFile)
	return r
}

func (s *Static) ImgUrl(staticFile string) (r string) {
	r = s.StaticUrl("img/" + staticFile)
	return r
}

func (s *Static) JsTag(staticFiles ...string) template.HTML {
	var r string
	for _, staticFile := range staticFiles {
		r += `<script type="text/javascript" src="` + s.JsUrl(staticFile) + `"></script>`
	}
	return template.HTML(r)
}

func (s *Static) CssTag(staticFiles ...string) template.HTML {
	var r string
	for _, staticFile := range staticFiles {
		r += `<link rel="stylesheet" href="` + s.CssUrl(staticFile) + `" />`
	}
	return template.HTML(r)
}

func (s *Static) ImgTag(staticFile string, attrs ...string) template.HTML {
	var attr string
	for i, l := 0, len(attrs); i+1 < l; i++ {
		var k, v string
		k = attrs[i]
		i++
		v = attrs[i]
		attr += ` ` + k + `="` + v + `"`
	}
	r := `<img src="` + s.ImgUrl(staticFile) + `"` + attr + ` />`
	return template.HTML(r)
}
