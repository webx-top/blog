package webx

import (
	"fmt"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/handler/mvc"
)

var _ = mvc.Data(&Output{})

type Output struct {
	context echo.Context
	Status  int
	Message interface{}
	For     interface{}
	Data    interface{}
}

func (c *Output) Render(tmpl string, code ...int) error {
	return c.context.Render(tmpl, c.Data, code...)
}

func (c *Output) String() string {
	return fmt.Sprintf(`%v`, c.Message)
}

func (c *Output) Assign(key string, val interface{}) {
	data, _ := c.Data.(echo.H)
	if data == nil {
		data = echo.H{}
	}
	data[key] = val
	c.Data = data
}

func (c *Output) Assignx(values *map[string]interface{}) {
	if values == nil {
		return
	}
	data, _ := c.Data.(echo.H)
	if data == nil {
		data = echo.H{}
	}
	for key, val := range *values {
		data[key] = val
	}
	c.Data = data
}

func (c *Output) SetTmplFuncs() {
	flash, ok := c.context.Session().Get(`webx:flash`).(*Output)
	if ok {
		c.context.Session().Delete(`webx:flash`).Save()
		c.context.SetFunc(`Status`, func() int {
			return flash.Status
		})
		c.context.SetFunc(`Message`, func() interface{} {
			return flash.Message
		})
		c.context.SetFunc(`For`, func() interface{} {
			return flash.For
		})
	} else {
		c.context.SetFunc(`Status`, func() int {
			return c.Status
		})
		c.context.SetFunc(`Message`, func() interface{} {
			return c.Message
		})
		c.context.SetFunc(`For`, func() interface{} {
			return c.For
		})
	}
	return
}

// Set 设置输出(code,message,for,data)
func (c *Output) Set(code int, args ...interface{}) {
	c.Status = code
	var hasData bool
	switch len(args) {
	case 3:
		c.Data = args[2]
		hasData = true
		fallthrough
	case 2:
		c.For = args[1]
		fallthrough
	case 1:
		c.Message = args[0]
		if !hasData {
			flash := &Output{
				Status:  c.Status,
				Message: c.Message,
				For:     c.For,
				Data:    nil,
			}
			c.context.Session().Set(`webx:flash`, flash).Save()
		}
	}
}
