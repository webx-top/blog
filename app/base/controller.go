package base

import (
	"net/http"

	"github.com/webx-top/echo"
)

func NewController() *Controller {
	return &Controller{
		tmplField: `Tmpl`,
		dataField: `Data`,
	}
}

type V map[string]interface{}

func (a V) Get(key string) interface{} {
	if v, ok := a[key]; ok {
		return v
	}
	return nil
}

func (a V) Set(key string, val interface{}) {
	a[key] = val
}

type Controller struct {
	tmplField string //模板名称字段名称
	dataField string //模板数据字段名称
}

//设置退出标记
func (a *Controller) Exit(c *echo.Context) error {
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
func (a *Controller) Tmpl(tmpl string, c *echo.Context) error {
	c.Set(a.tmplField, tmpl)
}

//模板数据赋值
func (a *Controller) Assign(data map[string]interface{}, c *echo.Context) {
	v, ok := c.Get(a.dataField).(V)
	if ok {
		for key, val := range data {
			v[key] = val
		}
	} else {
		v = V(data)
	}
	c.Set(a.dataField, v)
}

//渲染模板
func (a *Controller) Render(c *echo.Context) error {
	tmpl := c.Get(a.tmplField).(string)
	return c.Render(http.StatusOK, tmpl, c.Get(a.dataField))
}

func (a *Controller) Before(c *echo.Context) error {
	return nil
}

func (a *Controller) After(c *echo.Context) error {
	return a.Render(c)
}
