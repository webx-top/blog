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
package config

import (
	"strings"

	"github.com/admpub/log"
)

type Logger struct {
	Targets         []IItems  `json:"Targets"`
	MaxLevel        log.Level `json:"MaxLevel"`
	LevelName       string    `json:"LevelName"`
	Category        string    `json:"Category"`
	BufferSize      int       `json:"BufferSize"`
	CallStackDepth  int       `json:"CallStackDepth"`
	CallStackFilter string    `json:"CallStackFilter"`
}

func (l *Logger) Apply(logger *log.Logger) {
	if l.LevelName != `` {
		logger.SetLevel(l.LevelName)
	} else {
		logger.MaxLevel = l.MaxLevel
	}
	logger.Category = l.Category
	if l.BufferSize != 0 {
		logger.BufferSize = l.BufferSize
	}
	logger.CallStackDepth = l.CallStackDepth
	logger.CallStackFilter = l.CallStackFilter
	logger.Close()
	logger.Targets = make([]log.Target, 0)
	for _, c := range l.Targets {
		if c.GetBool(`Disabled`) {
			continue
		}
		typ := c.GetString(`Type`)
		var tg log.Target
		switch typ {
		case `Console`:
			t := log.NewConsoleTarget()
			if _, ok := c[`ColorMode`]; ok {
				t.ColorMode = c.GetBool(`ColorMode`)
			}
			if level, ok := log.GetLevel(c.GetString(`LevelName`)); ok {
				t.MaxLevel = level
			}
			t.Categories = c.GetStringSlice(`Categories`)
			tg = t
		case `File`:
			t := log.NewFileTarget()
			t.FileName = c.GetString(`FileName`)
			if _, ok := c[`Rotate`]; ok {
				t.Rotate = c.GetBool(`Rotate`)
			}
			if _, ok := c[`BackupCount`]; ok {
				t.BackupCount = c.GetInt(`BackupCount`)
			}
			if _, ok := c[`MaxBytes`]; ok {
				t.MaxBytes = c.GetInt64(`MaxBytes`)
			}
			if level, ok := log.GetLevel(c.GetString(`LevelName`)); ok {
				t.MaxLevel = level
			}
			t.Categories = c.GetStringSlice(`Categories`)
			tg = t
		case `Mail`:
			t := log.NewMailTarget()
			t.Host = c.GetString(`Host`)
			t.Username = c.GetString(`Username`)
			t.Password = c.GetString(`Password`)
			t.Subject = c.GetString(`Subject`)
			t.Sender = c.GetString(`Sender`)
			recipients := c.GetString(`Recipients`)
			t.Recipients = strings.Split(recipients, `;`)
			if _, ok := c[`BufferSize`]; ok {
				t.BufferSize = c.GetInt(`BufferSize`)
			}
			if level, ok := log.GetLevel(c.GetString(`LevelName`)); ok {
				t.MaxLevel = level
			}
			t.Categories = c.GetStringSlice(`Categories`)
			tg = t
		case `Network`:
			t := log.NewNetworkTarget()
			t.Network = c.GetString(`Network`)
			t.Address = c.GetString(`Address`)
			if _, ok := c[`Persistent`]; ok {
				t.Persistent = c.GetBool(`Persistent`)
			}
			if _, ok := c[`BufferSize`]; ok {
				t.BufferSize = c.GetInt(`BufferSize`)
			}
			if level, ok := log.GetLevel(c.GetString(`LevelName`)); ok {
				t.MaxLevel = level
			}
			t.Categories = c.GetStringSlice(`Categories`)
			tg = t
		}
		if tg != nil {
			logger.Targets = append(logger.Targets, tg)
		}
	}
	logger.Open()
}
