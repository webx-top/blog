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
	//"fmt"
	"io"
	"log"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/coscms/xorm"
	"github.com/coscms/xorm/core"
	_ "github.com/go-sql-driver/mysql"
	"github.com/webx-top/webx/lib/cachestore"
	//_ "github.com/ziutek/mymysql/godrv"
)

var CacheDir string = "data/cache"

func NewOrm(engine string, dsn string) (db *Orm, err error) {
	db = &Orm{}
	db.Engine, err = xorm.NewEngine(engine, dsn)
	if err != nil {
		log.Println("The database connection failed:", err)
		return
	}
	err = db.Engine.Ping()
	if err != nil {
		log.Println("The database ping failed:", err)
		return
	}
	return
}

type Orm struct {
	*xorm.Engine
	CacheStore   interface{}
	TransSession *xorm.Session
	PrefixMapper core.PrefixMapper
}

func (this *Orm) SetTimezone(timezone string) *Orm {
	switch timezone {
	case "UTC", "U":
		this.TZLocation = time.UTC
	default:
		this.TZLocation = time.Local
	}
	return this
}

func (this *Orm) SetPrefix(prefix string) *Orm {
	this.PrefixMapper = core.NewPrefixMapper(core.SnakeMapper{}, prefix)
	this.SetTableMapper(this.PrefixMapper)
	return this
}

//取得完整的表名
func (this *Orm) T(noPrefixTableName string) string {
	return this.PrefixMapper.Prefix + noPrefixTableName
}

func (this *Orm) SetLogger(out io.Writer) *Orm {
	this.Logger = xorm.NewSimpleLogger(out)
	return this
}

func (this *Orm) SetCacher(cs core.CacheStore) *Orm {
	this.CacheStore = cs
	if this.CacheStore != nil {
		var (
			cacher     *xorm.LRUCacher
			lifeTime   int32 = 86400
			maxEleSize int   = 999999999 //max element size
		)
		cacher = xorm.NewLRUCacher(this.CacheStore.(core.CacheStore), maxEleSize) //NewLRUCacher(store core.CacheStore, maxElementSize int)
		cacher.Expired = time.Duration(lifeTime) * time.Second
		this.SetDefaultCacher(cacher)
	}
	return this
}

func (this *Orm) Close() {
	//重置数据库连接
	if this.Engine != nil {
		if this.Engine.Cacher != nil {
			this.Engine.Cacher = nil
		}
		_ = this.Engine.Close()
		this.Engine = nil
	}

	//重置缓存对象
	if closer, ok := this.CacheStore.(cachestore.Closer); ok {
		closer.Close()
	}
	this.CacheStore = nil
}

func (this *Orm) Begin() *xorm.Session {
	if this.TransSession != nil {
		this.TransSession.Close()
	}
	this.TransSession = this.NewSession()
	err := this.TransSession.Begin()
	if err != nil {
		log.Println(err)
	}
	return this.TransSession
}

//事务是否已经开始
func (this *Orm) HasBegun() bool {
	if this.TransSession != nil {
		return true
	}
	return false
}

func (this *Orm) TSess() *xorm.Session { // TransSession
	if this.HasBegun() == false {
		return this.Begin()
	}
	return this.TransSession
}

func (this *Orm) Trans(fn func() error) *Orm {
	begun := this.HasBegun()
	result := fn()
	if !begun {
		this.End(result == nil)
	}
	return this
}

func (this *Orm) Sess() *xorm.Session { // TransSession or Session
	if this.HasBegun() == false {
		var session *xorm.Session = this.NewSession()
		session.IsAutoClose = true
		return session
	}
	return this.TransSession
}

func (this *Orm) End(result bool) (err error) {
	if result {
		err = this.TransSession.Commit()
	} else {
		err = this.TransSession.Rollback()
	}
	if err != nil {
		log.Println(err)
	}
	this.TransSession.Close()
	this.TransSession = nil
	return
}

//验证字段
func (this *Orm) VerifyField(v interface{}, field string) string {
	if field == "" {
		return ""
	}
	val := reflect.ValueOf(v)
	typ := val.Type()
	if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct {
		typ = typ.Elem()
		val = val.Elem()
	}
	name := this.Engine.ColumnMapper.Table2Obj(field)
	if field == name {
		return ""
	}
	for i := 0; i < val.NumField(); i++ {
		vt := typ.Field(i)
		if name == vt.Name {
			return field
		}
	}
	return ""
}

//验证字段并返回有效的字段切片
func (this *Orm) VerifyFields(v interface{}, fn func(reflect.Type, reflect.Value, int) string) []string {
	ret := make([]string, 0)
	val := reflect.ValueOf(v)
	typ := val.Type()
	if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct {
		typ = typ.Elem()
		val = val.Elem()
	}
	for i := 0; i < val.NumField(); i++ {
		v := fn(typ, val, i)
		if v != "" {
			ret = append(ret, v)
		}
	}
	return ret
}

var (
	searchMultiKwRule   *regexp.Regexp = regexp.MustCompile(`[\s]+`)       //多个关键词
	searchMultiIdRule   *regexp.Regexp = regexp.MustCompile(`[\s,]+`)      //多个Id
	searchIdRule        *regexp.Regexp = regexp.MustCompile(`^[\s\d-,]+$`) //多个Id
	searchParagraphRule *regexp.Regexp = regexp.MustCompile(`"[^"]+"`)     //段落
)

/**
 * 搜索某个字段
 * @param field 字段名。支持搜索多个字段，各个字段之间用半角逗号“,”隔开
 * @param keywords 关键词
 * @param idFields 如要搜索id字段需要提供id字段名
 * @return sql
 * @author swh <swh@admpub.com>
 */
func (this *Orm) SearchField(field string, keywords string, idFields ...string) string {
	if keywords == "" || field == "" {
		return ""
	}
	idField := ""
	if len(idFields) > 0 {
		idField = idFields[0]
	}
	keywords = strings.TrimSpace(keywords)
	sql := ""
	if idField != "" && searchIdRule.FindString(keywords) != "" {
		kws := searchMultiIdRule.Split(keywords, -1)
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
					sql += " OR " + this.Quote(idField) + "<='" + v + "'"
					continue
				}
				if v[length-1] == '-' {
					v = strings.Trim(v, "-")
					if v == "" {
						continue
					}
					sql += " OR " + this.Quote(idField) + ">='" + v + "'"
					continue
				}

				v = strings.Trim(v, "-")
				if v == "" {
					continue
				}
				vs := strings.SplitN(v, "-", 2)
				sql += " OR " + this.Quote(idField) + " BETWEEN ('" + vs[0] + "','" + vs[1] + "')"
			} else {
				sql += " OR " + this.Quote(idField) + "='" + v + "'"
			}
		}
		if sql != "" {
			sql = "(" + sql[3:] + ")"
		}
	} else {
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
				if v == "" {
					continue
				}
				cond := ""
				if strings.Contains(v, "||") {
					vals := strings.Split(v, "||")
					for _, val := range vals {
						val = AddSlashes(val)
						val = strings.Replace(val, "_", `\_`, -1)
						val = strings.Replace(val, "%", `\%`, -1)
						cond += " OR 'FIELD' LIKE '%" + val + "%'"
					}
					cond = cond[3:]
				} else {
					v = AddSlashes(v)
					v = strings.Replace(v, "_", `\_`, -1)
					v = strings.Replace(v, "%", `\%`, -1)
				}
				for k, f := range fs {
					if cond != "" {
						ds[k] += " AND (" + strings.Replace(cond, "'FIELD'", this.Quote(f), -1) + ")"
						continue
					}
					ds[k] += " AND " + this.Quote(f) + " LIKE '%" + v + "%'"
				}
			}
			for _, v := range ds {
				if v != "" {
					sql += " OR (" + v[4:] + ")"
				}
			}
			if sql != "" {
				sql = "(" + sql[3:] + ")"
			}
		} else {
			for _, v := range kws {
				v = strings.TrimSpace(v)
				if v == "" {
					continue
				}
				if strings.Contains(v, "||") {
					vals := strings.Split(v, "||")
					cond := ""
					for _, val := range vals {
						val = AddSlashes(val)
						val = strings.Replace(val, "_", `\_`, -1)
						val = strings.Replace(val, "%", `\%`, -1)
						cond += " OR " + this.Quote(field) + " LIKE '%" + val + "%'"
					}
					sql += " AND (" + cond[3:] + ")"
					continue
				}
				v = AddSlashes(v)
				v = strings.Replace(v, "_", `\_`, -1)
				v = strings.Replace(v, "%", `\%`, -1)
				sql += " AND " + this.Quote(field) + " LIKE '%" + v + "%'"
			}
			if sql != "" {
				sql = "(" + sql[4:] + ")"
			}
		}
	}
	return sql
}

func AddSlashes(s string) string {
	s = strings.Replace(s, `\`, `\\`, -1)
	s = strings.Replace(s, "'", `\'`, -1)
	return s
}
