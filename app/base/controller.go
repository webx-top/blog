package base

import (
	"net/http"

	"github.com/webx-top/echo"
	"github.com/webx-top/webx/lib/htmlcache"
	"github.com/webx-top/webx/lib/i18n"
)

func NewController() *Controller {
	return &Controller{
		tmplField: `Tmpl`,
		dataField: `Data`,
	}
}

type Controller struct {
	tmplField string //模板名称字段名称
	dataField string //模板数据字段名称
}

//设置退出标记
func (a *Controller) Exit(c echo.Context) {
	c.Set(`web:exit`, true)
}

//获取模板名称字段名
func (a *Controller) TmplField() string {
	return a.tmplField
}

//获取模板数据字段名
func (a *Controller) DataField() string {
	return a.dataField
}

//指定模板
func (a *Controller) Tmpl(tmpl string, c echo.Context) {
	c.Set(a.tmplField, tmpl)
}

//模板数据赋值
func (a *Controller) Assign(data interface{}, c echo.Context) {
	switch data.(type) {
	case map[string]interface{}:
		d := data.(map[string]interface{})
		v, ok := c.Get(a.dataField).(map[string]interface{})
		if ok {
			for key, val := range d {
				v[key] = val
			}
		} else {
			v = d
		}
		c.Set(a.dataField, v)
	default:
		c.Set(a.dataField, data)
	}
}

//渲染模板
func (a *Controller) Render(c echo.Context) error {
	if ignore, _ := c.Get(`webx:ignoreRender`).(bool); ignore {
		return nil
	}
	return htmlcache.Render(c, http.StatusOK)
}

func (a *Controller) Before(c echo.Context) error {
	Xsrf.Register(c)
	c.SetFunc("Query", c.Query)
	c.SetFunc("Form", c.Form)
	a.Assign(map[string]interface{}{"Status": 1, "Message": "", "Path": c.Path()}, c)
	return nil
}

func (a *Controller) After(c echo.Context) error {
	return a.Render(c)
}

func (a *Controller) Lang(c echo.Context) string {
	lang, _ := c.Get(`webx:language`).(string)
	if lang == `` {
		lang = DefaultLang
	}
	return lang
}

//TODO: 移到echo.Context中
func (a *Controller) T(c echo.Context, key string, args ...interface{}) string {
	return i18n.T(a.Lang(c), key, args...)
}
