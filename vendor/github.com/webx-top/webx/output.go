package webx

import (
	"fmt"

	"github.com/webx-top/echo"
)

var _ = echo.Data(&Output{})

type Output struct {
	context echo.Context
	Status  echo.State
	Summary string `json:",omitempty" xml:",omitempty"`
	Message interface{}
	For     interface{} `json:",omitempty" xml:",omitempty"`
	Data    interface{} `json:",omitempty" xml:",omitempty"`
}

func (c *Output) Render(tmpl string, code ...int) error {
	return c.context.Render(tmpl, c.Data, code...)
}

func (c *Output) Reset() echo.Data {
	c.Status = echo.State(0)
	c.Summary = ``
	c.Message = nil
	c.For = nil
	c.Data = nil
	return c
}

func (c *Output) Gets() (echo.State, interface{}, interface{}, interface{}) {
	return c.Status, c.Message, c.For, c.Data
}

func (c *Output) GetCode() echo.State {
	return c.Status
}

func (c *Output) GetInfo() interface{} {
	return c.Message
}

func (c *Output) GetZone() interface{} {
	return c.For
}

func (c *Output) GetData() interface{} {
	return c.Data
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
		c.context.SetFunc(`Status`, func() echo.State {
			return flash.Status
		})
		c.context.SetFunc(`Message`, func() interface{} {
			return flash.Message
		})
		c.context.SetFunc(`For`, func() interface{} {
			return flash.For
		})
	} else {
		c.context.SetFunc(`Status`, func() echo.State {
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
func (c *Output) Set(code int, args ...interface{}) echo.Data {
	c.Status = echo.State(code)
	c.Summary = c.Status.String()
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
				Summary: c.Summary,
				Message: c.Message,
				For:     c.For,
				Data:    nil,
			}
			c.context.Session().Set(`webx:flash`, flash).Save()
		}
	}
	return c
}

func (c *Output) SetContext(ctx echo.Context) echo.Data {
	c.context = ctx
	return c
}

//SetError 设置错误
func (c *Output) SetError(err error, args ...int) echo.Data {
	if err != nil {
		if len(args) > 0 {
			c.SetCode(args[0])
		} else {
			c.SetCode(0)
		}
		c.Message = err.Error()
	} else {
		c.SetCode(1)
	}
	return c
}

//SetCode 设置状态码
func (c *Output) SetCode(code int) echo.Data {
	c.Status = echo.State(code)
	c.Summary = c.Status.String()
	return c
}

//SetInfo 设置提示信息
func (c *Output) SetInfo(info interface{}, args ...int) echo.Data {
	c.Message = info
	if len(args) > 0 {
		c.SetCode(args[0])
	}
	return c
}

//SetZone 设置提示区域
func (c *Output) SetZone(zone interface{}) echo.Data {
	c.For = zone
	return c
}

//SetData 设置正常数据
func (c *Output) SetData(data interface{}, args ...int) echo.Data {
	c.Data = data
	if len(args) > 0 {
		c.SetCode(args[0])
	} else {
		c.SetCode(1)
	}
	return c
}
