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
package database

import (
	"log"
	"regexp"
	"strings"
	"time"

	. "github.com/coscms/xorm"
	"github.com/coscms/xorm/core"
	_ "github.com/go-sql-driver/mysql"
	"github.com/webx-top/com"
	cachestore "github.com/webx-top/webx/lib/store/cache"
	//_ "github.com/ziutek/mymysql/godrv"
	xlog "github.com/admpub/log"
	"github.com/webx-top/webx/lib/config"
	"github.com/webx-top/webx/lib/database/balancer"
)

func NewOrm(engine string, cfgMaster *config.DB, cfgSlaves ...*config.DB) (db *Orm, err error) {
	db = &Orm{}
	dsn := cfgMaster.Dsn()
	for _, cfg := range cfgSlaves {
		dsn += `;` + cfg.Dsn()
	}
	db.Balancer, err = balancer.New(engine, dsn)
	if err != nil {
		log.Println("The database connection failed:", err)
		return
	}
	err = db.Ping()
	if err != nil {
		log.Println("The database ping failed:", err)
		return
	}
	db.SetToAllDbs(func(x *Engine) {
		x.OpenLog()
	})
	db.SetLogger(xlog.New(`orm`))
	return
}

type Orm struct {
	*balancer.Balancer
	CacheStore   interface{}
	PrefixMapper core.PrefixMapper
}

func (this *Orm) SetToAllDbs(setter func(*Engine)) *Orm {
	for _, db := range this.GetAllDbs() {
		setter(db)
	}
	return this
}

func (this *Orm) SetTimezone(timezone string) *Orm {
	var location *time.Location
	switch timezone {
	case "UTC", "U":
		location = time.UTC
	case "Local", "L", "":
		location = time.Local
	default:
		var err error
		location, err = time.LoadLocation(timezone)
		if err != nil {
			log.Println(err)
		}
	}

	this.SetToAllDbs(func(x *Engine) {
		x.TZLocation = location
	})
	return this
}

func (this *Orm) SetPrefix(prefix string) *Orm {
	this.PrefixMapper = core.NewPrefixMapper(core.SnakeMapper{}, prefix)

	this.SetToAllDbs(func(x *Engine) {
		x.SetTblMapper(this.PrefixMapper)
	})
	return this
}

// TblName 取得完整的表名
func (this *Orm) TblName(noPrefixTableName string) string {
	return this.Engine.TableName(noPrefixTableName)
}

func (this *Orm) JoinAlias(table string, alias string) []string {
	return []string{this.Engine.TableName(table), alias}
}

func (this *Orm) JoinAliasObj(table interface{}, alias string) []interface{} {
	return []interface{}{table, alias}
}

func (this *Orm) SetLogger(logger *xlog.Logger) *Orm {
	this.Balancer.TraceOn(``, logger)
	return this
}

func (this *Orm) SetCacher(cs core.CacheStore) *Orm {
	this.CacheStore = cs
	if this.CacheStore != nil {
		var (
			cacher     *LRUCacher
			lifeTime   int32 = 86400
			maxEleSize       = 999999999 //max element size
		)
		//NewLRUCacher(store core.CacheStore, maxElementSize int)
		cacher = NewLRUCacher(this.CacheStore.(core.CacheStore), maxEleSize)
		cacher.Expired = time.Duration(lifeTime) * time.Second
		this.SetToAllDbs(func(x *Engine) {
			x.SetDefaultCacher(cacher)
		})
	}
	return this
}

func (this *Orm) Close() {
	//重置数据库连接
	this.SetToAllDbs(func(x *Engine) {
		if x != nil {
			if x.Cacher != nil {
				x.Cacher = nil
			}
			_ = x.Close()
			x = nil
		}
	})

	//重置缓存对象
	if closer, ok := this.CacheStore.(cachestore.Closer); ok {
		closer.Close()
	}
	this.CacheStore = nil
}

var (
	searchMultiKwRule   = regexp.MustCompile(`[\s]+`)                        //多个关键词
	splitMultiIdRule    = regexp.MustCompile(`[^\d-]+`)                      //多个Id
	searchCompareRule   = regexp.MustCompile(`^[><!][=]?[\d]+(?:\.[\d]+)?$`) //多个Id
	searchIdRule        = regexp.MustCompile(`^[\s\d-,]+$`)                  //多个Id
	searchParagraphRule = regexp.MustCompile(`"[^"]+"`)                      //段落
)

func (this *Orm) CompareField(idField string, keywords string) string {
	if len(keywords) == 0 || len(idField) == 0 {
		return ""
	}
	var op string
	if keywords[1] == '=' {
		op = keywords[0:2]
	} else {
		op = keywords[0:1]
	}
	field := this.Quote(idField)
	sql := "(" + field + op + keywords[2:] + ")"
	return sql
}

func IsCompareField(keywords string) bool {
	return len(searchCompareRule.FindString(keywords)) > 0
}

func IsRangeField(keywords string) bool {
	return len(searchIdRule.FindString(keywords)) > 0
}

func (this *Orm) SearchFields(fields []string, keywords string, idFields ...string) string {
	if len(keywords) == 0 || len(fields) == 0 {
		return ""
	}
	idField := ""
	if len(idFields) > 0 {
		idField = idFields[0]
	}
	keywords = strings.TrimSpace(keywords)
	sql := ""
	if len(idField) > 0 {
		switch {
		case IsCompareField(keywords):
			return this.CompareField(idField, keywords)
		case IsRangeField(keywords):
			return this.RangeField(idField, keywords)
		}
	}
	var paragraphs []string = make([]string, 0)
	keywords = searchParagraphRule.ReplaceAllStringFunc(keywords, func(paragraph string) string {
		paragraph = strings.Trim(paragraph, `"`)
		paragraphs = append(paragraphs, paragraph)
		return ""
	})
	kws := searchMultiKwRule.Split(keywords, -1)
	kws = append(kws, paragraphs...)
	ds := make([]string, len(fields))
	for _, v := range kws {
		v = strings.TrimSpace(v)
		if len(v) == 0 {
			continue
		}
		cond := ""
		if strings.Contains(v, "||") {
			vals := strings.Split(v, "||")
			for _, val := range vals {
				val = AddSlashes(val, '_', '%')
				cond += " OR 'FIELD' LIKE '%" + val + "%'"
			}
			cond = cond[3:]
		} else {
			v = AddSlashes(v, '_', '%')
		}
		for k, f := range fields {
			f = this.Quote(f)
			if len(cond) > 0 {
				ds[k] += " AND (" + strings.Replace(cond, "'FIELD'", f, -1) + ")"
				continue
			}
			ds[k] += " AND " + f + " LIKE '%" + v + "%'"
		}
	}
	for _, v := range ds {
		if len(v) > 0 {
			sql += " OR (" + v[4:] + ")"
		}
	}
	if len(sql) > 0 {
		sql = "(" + sql[3:] + ")"
	}
	return sql
}

/**
 * 搜索某个字段
 * @param field 字段名。支持搜索多个字段，各个字段之间用半角逗号“,”隔开
 * @param keywords 关键词
 * @param idFields 如要搜索id字段需要提供id字段名
 * @return sql
 * @author swh <swh@admpub.com>
 */
func (this *Orm) SearchField(field string, keywords string, idFields ...string) string {
	if len(keywords) == 0 || len(field) == 0 {
		return ""
	}
	idField := ""
	if len(idFields) > 0 {
		idField = idFields[0]
	}
	keywords = strings.TrimSpace(keywords)
	sql := ""
	if len(idField) > 0 {
		switch {
		case IsCompareField(keywords):
			return this.CompareField(idField, keywords)
		case IsRangeField(keywords):
			return this.RangeField(idField, keywords)
		}
	}
	var paragraphs []string = make([]string, 0)
	keywords = searchParagraphRule.ReplaceAllStringFunc(keywords, func(paragraph string) string {
		paragraph = strings.Trim(paragraph, `"`)
		paragraphs = append(paragraphs, paragraph)
		return ""
	})
	kws := searchMultiKwRule.Split(keywords, -1)
	kws = append(kws, paragraphs...)
	if strings.Contains(field, ",") {
		fs := strings.Split(field, ",")
		ds := make([]string, len(fs))
		for _, v := range kws {
			v = strings.TrimSpace(v)
			if len(v) == 0 {
				continue
			}
			cond := ""
			if strings.Contains(v, "||") {
				vals := strings.Split(v, "||")
				for _, val := range vals {
					val = AddSlashes(val, '_', '%')
					cond += " OR 'FIELD' LIKE '%" + val + "%'"
				}
				cond = cond[3:]
			} else {
				v = AddSlashes(v, '_', '%')
			}
			for k, f := range fs {
				f = this.Quote(f)
				if len(cond) > 0 {
					ds[k] += " AND (" + strings.Replace(cond, "'FIELD'", f, -1) + ")"
					continue
				}
				ds[k] += " AND " + f + " LIKE '%" + v + "%'"
			}
		}
		for _, v := range ds {
			if len(v) > 0 {
				sql += " OR (" + v[4:] + ")"
			}
		}
		if len(sql) > 0 {
			sql = "(" + sql[3:] + ")"
		}
	} else {
		field = this.Quote(field)
		for _, v := range kws {
			v = strings.TrimSpace(v)
			if len(v) == 0 {
				continue
			}
			if strings.Contains(v, "||") {
				vals := strings.Split(v, "||")
				cond := ""
				for _, val := range vals {
					val = AddSlashes(val, '_', '%')
					cond += " OR " + field + " LIKE '%" + val + "%'"
				}
				sql += " AND (" + cond[3:] + ")"
				continue
			}
			v = AddSlashes(v, '_', '%')
			sql += " AND " + field + " LIKE '%" + v + "%'"
		}
		if len(sql) > 0 {
			sql = "(" + sql[4:] + ")"
		}
	}
	return sql
}

func (this *Orm) RangeField(idField string, keywords string) string {
	if len(keywords) == 0 || len(idField) == 0 {
		return ""
	}
	var sql string
	var logic string
	keywords = strings.TrimSpace(keywords)
	kws := splitMultiIdRule.Split(keywords, -1)

	field := this.Quote(idField)
	for _, v := range kws {
		length := len(v)
		if length < 1 {
			continue
		}
		if strings.Contains(v, "-") {
			if length < 2 {
				continue
			}
			if v[0] == '-' {
				v = strings.Trim(v, "-")
				if v == "" {
					continue
				}
				sql += logic + field + "<='" + v + "'"
				logic = " OR "
				continue
			}
			if v[length-1] == '-' {
				v = strings.Trim(v, "-")
				if v == "" {
					continue
				}
				sql += logic + field + ">='" + v + "'"
				logic = " OR "
				continue
			}

			v = strings.Trim(v, "-")
			if v == "" {
				continue
			}
			vs := strings.SplitN(v, "-", 2)
			sql += logic + field + " BETWEEN '" + vs[0] + "' AND '" + vs[1] + "'"
			logic = " OR "
		} else {
			sql += logic + field + "='" + v + "'"
			logic = " OR "
		}
	}
	if len(sql) > 0 {
		sql = "(" + sql + ")"
	}
	return sql
}

func (this *Orm) EqField(field string, keywords string) string {
	if len(keywords) == 0 || len(field) == 0 {
		return ""
	}
	keywords = strings.TrimSpace(keywords)
	return this.Quote(field) + "='" + AddSlashes(keywords) + "'"
}

/**
 * GenDateRangeSql :
 * 生成日期范围SQL语句
 * @param cond 已有条件sql
 * @param field 字段名。支持搜索多个字段，各个字段之间用半角逗号“,”隔开
 * @param keywords 关键词
 * @return sql
 */
func (this *Orm) GenDateRangeSql(cond *string, field string, keywords string) {
	if len(keywords) == 0 || len(field) == 0 {
		return
	}
	var skwd, skwdExt string
	dataRange := strings.Split(keywords, ` - `)
	skwd = dataRange[0]
	if len(dataRange) > 1 {
		skwdExt = dataRange[1]
	}
	//日期范围
	dateBegin := com.StrToTime(skwd + ` 00:00:00`)
	if len(*cond) > 0 {
		*cond += ` AND `
	}
	field = this.Quote(field)
	*cond += field + `>=` + com.Str(dateBegin)
	if len(skwdExt) > 0 {
		dateEnd := com.StrToTime(skwd + ` 23:59:59`)
		*cond += ` AND ` + field + `<=` + com.Str(dateEnd)
	}
}
