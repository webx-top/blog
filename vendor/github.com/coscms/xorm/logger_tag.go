// Copyright 2015 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	"fmt"

	"github.com/coscms/xorm/core"
)

type CLogger struct {
	Name      string
	Disabled  bool
	Processor func(tag string, format string, args []interface{}) (string, []interface{})
	Logger    core.ILogger
}

func (c *CLogger) Error(v ...interface{}) {
	if c.Disabled {
		return
	}
	if c.Processor != nil {
		_, v = c.Processor(c.Name, ``, v)
		if v == nil {
			return
		}
	}
	c.Logger.Error(v...)
}

func (c *CLogger) Errorf(format string, v ...interface{}) {
	if c.Disabled {
		return
	}
	if c.Processor != nil {
		format, v = c.Processor(c.Name, format, v)
		if v == nil {
			return
		}
	}
	c.Logger.Errorf(format, v...)
}

func (c *CLogger) Debug(v ...interface{}) {
	if c.Disabled {
		return
	}
	if c.Processor != nil {
		_, v = c.Processor(c.Name, ``, v)
		if v == nil {
			return
		}
	}
	c.Logger.Debug(v...)
}

func (c *CLogger) Debugf(format string, v ...interface{}) {
	if c.Disabled {
		return
	}
	if c.Processor != nil {
		format, v = c.Processor(c.Name, format, v)
		if v == nil {
			return
		}
	}
	c.Logger.Debugf(format, v...)
}

func (c *CLogger) Info(v ...interface{}) {
	if c.Disabled {
		return
	}
	if c.Processor != nil {
		_, v = c.Processor(c.Name, ``, v)
		if v == nil {
			return
		}
	}
	c.Logger.Info(v...)
}

func (c *CLogger) Infof(format string, v ...interface{}) {
	if c.Disabled {
		return
	}
	if c.Processor != nil {
		format, v = c.Processor(c.Name, format, v)
		if v == nil {
			return
		}
	}
	c.Logger.Infof(format, v...)
}

func (c *CLogger) Warn(v ...interface{}) {
	if c.Disabled {
		return
	}
	if c.Processor != nil {
		_, v = c.Processor(c.Name, ``, v)
		if v == nil {
			return
		}
	}
	c.Logger.Warn(v...)
}

func (c *CLogger) Warnf(format string, v ...interface{}) {
	if c.Disabled {
		return
	}
	if c.Processor != nil {
		format, v = c.Processor(c.Name, format, v)
		if v == nil {
			return
		}
	}
	c.Logger.Warnf(format, v...)
}

type TLogger struct {
	SQL   *CLogger
	Event *CLogger
	Cache *CLogger
	ETime *CLogger
	Base  *CLogger
	Other *CLogger
}

func (t *TLogger) Open(tags ...string) {
	if len(tags) == 0 {
		t.SQL.Disabled = false
		t.Event.Disabled = false
		t.Cache.Disabled = false
		t.ETime.Disabled = false
		t.Base.Disabled = false
		t.Other.Disabled = false
		return
	}
	for _, tag := range tags {
		t.SetStatusByName(tag, false)
	}
}

func (t *TLogger) Close(tags ...string) {
	if len(tags) == 0 {
		t.SQL.Disabled = true
		t.Event.Disabled = true
		t.Cache.Disabled = true
		t.ETime.Disabled = true
		t.Base.Disabled = true
		t.Other.Disabled = true
		return
	}
	for _, tag := range tags {
		t.SetStatusByName(tag, true)
	}
}

func (t *TLogger) SetStatusByName(tag string, status bool) {
	switch tag {
	case "sql":
		t.SQL.Disabled = status
	case "event":
		t.Event.Disabled = status
	case "cache":
		t.Cache.Disabled = status
	case "etime":
		t.ETime.Disabled = status
	case "base":
		t.Base.Disabled = status
	case "other":
		t.Other.Disabled = status
	}
}

func (t *TLogger) SetLogger(logger core.ILogger) {
	t.SQL.Logger = logger
	t.Event.Logger = logger
	t.Cache.Logger = logger
	t.ETime.Logger = logger
	t.Base.Logger = logger
	t.Other.Logger = logger
}

func (t *TLogger) SetProcessor(processor func(tag string, format string, args []interface{}) (string, []interface{})) {
	t.SQL.Processor = processor
	t.Event.Processor = processor
	t.Cache.Processor = processor
	t.ETime.Processor = processor
	t.Base.Processor = processor
	t.Other.Processor = processor
}

func NewTLogger(logger core.ILogger) *TLogger {
	return &TLogger{
		SQL:   &CLogger{Name: "sql", Disabled: false, Processor: defaultLogProcessor, Logger: logger},
		Event: &CLogger{Name: "event", Disabled: false, Processor: defaultLogProcessor, Logger: logger},
		Cache: &CLogger{Name: "cache", Disabled: false, Processor: defaultLogProcessor, Logger: logger},
		ETime: &CLogger{Name: "etime", Disabled: false, Processor: defaultLogProcessor, Logger: logger},
		Base:  &CLogger{Name: "base", Disabled: false, Processor: defaultLogProcessor, Logger: logger},
		Other: &CLogger{Name: "other", Disabled: false, Processor: defaultLogProcessor, Logger: logger},
	}
}

var defaultLogProcessor = func(tag string, format string, args []interface{}) (string, []interface{}) {
	if format == "" {
		if len(args) > 0 {
			args[0] = "[" + tag + "] " + fmt.Sprintf("%v", args[0])
		}
		return format, args
	}
	format = "[" + tag + "] " + format
	return format, args
}
