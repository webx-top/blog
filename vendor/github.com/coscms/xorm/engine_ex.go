package xorm

import (
	"strings"

	"github.com/coscms/xorm/core"
)

// =====================================
// 增加Engine结构体中的方法
// @author Admpub <swh@admpub.com>
// =====================================

func (engine *Engine) Init() {
	engine.RelTagIdentifier = `rel`
	engine.AliasTagIdentifier = `alias`
	engine.TLogger = NewTLogger(engine.logger)
}

func (engine *Engine) SetTblMapper(mapper core.IMapper) {
	if prefixMapper, ok := mapper.(core.PrefixMapper); ok {
		engine.TablePrefix = prefixMapper.Prefix
	} else if suffixMapper, ok := mapper.(core.SuffixMapper); ok {
		engine.TableSuffix = suffixMapper.Suffix
	}
	engine.TableMapper = mapper
}

func (engine *Engine) OpenLog(types ...string) {
	engine.TLogger.Open(types...)
}

func (engine *Engine) CloseLog(types ...string) {
	engine.TLogger.Close(types...)
}

func (engine *Engine) QuoteWithDelim(s, d string) string {
	return engine.Quote(strings.Replace(s, d, engine.Quote(d), -1))
}

func (engine *Engine) ToSQL(sql string) core.SQL {
	return core.SQL(sql)
}
