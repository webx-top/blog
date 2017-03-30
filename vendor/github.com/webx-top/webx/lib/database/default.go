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
	. "github.com/coscms/xorm"
	"github.com/coscms/xorm/core"
	_ "github.com/go-sql-driver/mysql"
	//_ "github.com/ziutek/mymysql/godrv"
	xlog "github.com/admpub/log"
	"github.com/webx-top/webx/lib/config"
)

var Default = &Orm{}

func Connect(engine string, cfgMaster *config.DB, cfgSlaves ...*config.DB) error {
	return Default.Connect(engine, cfgMaster, cfgSlaves...)
}

func Connected() bool {
	return Default.Connected()
}

func SetToAllDbs(setter func(*Engine)) *Orm {
	return Default.SetToAllDbs(setter)
}

func OpenLog(tags ...string) *Orm {
	return Default.OpenLog(tags...)
}

func CloseLog(tags ...string) *Orm {
	return Default.CloseLog(tags...)
}

func SetTimezone(timezone string) *Orm {
	return Default.SetTimezone(timezone)
}

func SetPrefix(prefix string) *Orm {
	return Default.SetPrefix(prefix)
}

// TblName 取得完整的表名
func TblName(noPrefixTableName string) string {
	return Default.TblName(noPrefixTableName)
}

func JoinAlias(table string, alias string) []string {
	return Default.JoinAlias(table, alias)
}

func JoinAliasObj(table interface{}, alias string) []interface{} {
	return Default.JoinAliasObj(table, alias)
}

func SetLogger(logger *xlog.Logger) *Orm {
	return Default.SetLogger(logger)
}

func SetCacher(cs core.CacheStore) *Orm {
	return Default.SetCacher(cs)
}

func Close() {
	Default.Close()
}

func CompareField(idField string, keywords string) string {
	return Default.CompareField(idField, keywords)
}

func SearchFields(fields []string, keywords string, idFields ...string) string {
	return Default.SearchFields(fields, keywords, idFields...)
}

/**
 * 搜索某个字段
 * @param field 字段名。支持搜索多个字段，各个字段之间用半角逗号“,”隔开
 * @param keywords 关键词
 * @param idFields 如要搜索id字段需要提供id字段名
 * @return sql
 * @author swh <swh@admpub.com>
 */
func SearchField(field string, keywords string, idFields ...string) string {
	return Default.SearchField(field, keywords, idFields...)
}

func RangeField(idField string, keywords string) string {
	return Default.RangeField(idField, keywords)
}

func EqField(field string, keywords string) string {
	return Default.EqField(field, keywords)
}

/**
 * GenDateRangeSql :
 * 生成日期范围SQL语句
 * @param cond 已有条件sql
 * @param field 字段名。支持搜索多个字段，各个字段之间用半角逗号“,”隔开
 * @param keywords 关键词
 * @return sql
 */
func GenDateRangeSql(cond *string, field string, keywords string) {
	Default.GenDateRangeSql(cond, field, keywords)
}
