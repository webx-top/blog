package base

import (
	"github.com/webx-top/echo"
	X "github.com/webx-top/webx"
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

//渲染模板
func (a *Controller) Render(c *X.Context) error {
	if ignore, _ := c.Get(`webx:ignoreRender`).(bool); ignore {
		return nil
	}
	return c.Display()
}

func (a *Controller) Before(c *X.Context) error {
	Xsrf.Register(c)
	c.SetFunc("Query", c.Query)
	c.SetFunc("Form", c.Form)
	c.SetFunc("Form", c.Path)
	c.Assign("Path", c.Path())
	return nil
}

func (a *Controller) After(c *X.Context) error {
	return a.Render(c)
}

func (a *Controller) Lang(c *X.Context) string {
	if c.Language == `` {
		c.Language = DefaultLang
	}
	return c.Language
}

//TODO: 移到echo.Context中
func (a *Controller) T(c *X.Context, key string, args ...interface{}) string {
	return i18n.T(a.Lang(c), key, args...)
}

func (a *Controller) X(c echo.Context) *X.Context {
	return X.X(c)
}
