// Copyright 2015 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	"github.com/admpub/log"
	"github.com/coscms/xorm/core"
)

// AdmpubLogger is the default implment of core.ILogger
type AdmpubLogger struct {
	*log.Logger
	level   core.LogLevel
	showSQL bool
}

var _ core.ILogger = &AdmpubLogger{}

func NewAdmpubLogger(args ...*log.Logger) *AdmpubLogger {
	l := &AdmpubLogger{
		level: DEFAULT_LOG_LEVEL,
	}
	if len(args) > 0 {
		l.Logger = args[0]
	} else {
		l.Logger = log.GetLogger(`xorm`)
	}
	return l
}

// Level implement core.ILogger
func (s *AdmpubLogger) Level() core.LogLevel {
	return s.level
}

// SetLevel implement core.ILogger
func (s *AdmpubLogger) SetLevel(l core.LogLevel) {
	le := ``
	switch l {
	case 0:
		le = `Debug`
	case 1:
		le = `Info`
	case 2:
		le = `Warn`
	case 3:
		le = `Error`
	case 4:
		le = `Fatal`
	default:
		le = `Fatal`
	}
	s.level = l
	s.Logger.SetLevel(le)
	return
}

// ShowSQL implement core.ILogger
func (s *AdmpubLogger) ShowSQL(show ...bool) {
	if len(show) == 0 {
		s.showSQL = true
		return
	}
	s.showSQL = show[0]
}

// IsShowSQL implement core.ILogger
func (s *AdmpubLogger) IsShowSQL() bool {
	return s.showSQL
}
