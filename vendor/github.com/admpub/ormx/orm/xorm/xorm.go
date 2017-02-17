package xorm

import (
	"database/sql"
	"io"

	"github.com/admpub/log"
	"github.com/admpub/ormx"
	"github.com/admpub/ormx/orm"
	"github.com/coscms/xorm"
	"github.com/coscms/xorm/core"
)

func New(driverName string, sources string) (*Balancer, error) {
	b, e := ormx.NewBalancer(Connect(), driverName, sources)
	if e != nil {
		return nil, e
	}
	return &Balancer{
		Xorm:     b.ORM.(*Xorm),
		Balancer: b,
	}, nil
}

func Connect() orm.Connector {
	return func(driverName string, dsn string) (orm.ORM, error) {
		engine, err := xorm.NewEngine(driverName, dsn)
		if err != nil {
			return nil, err
		}
		a := &Xorm{
			Engine: engine,
		}
		return a, nil
	}
}

type Xorm struct {
	*xorm.Engine
}

func (x *Xorm) TraceOn(prefix string, logger interface{}) {
	if v, y := logger.(io.Writer); y {
		x.SetLogger(xorm.NewSimpleLogger2(v, prefix, xorm.DEFAULT_LOG_FLAG))
		return
	}
	x.SetLogger(xorm.NewAdmpubLoggerWithPrefix(prefix, logger.(*log.Logger)))
}

func (x *Xorm) DB() *sql.DB {
	return x.Engine.DB().DB
}

func (x *Xorm) Prepare(query string) (*sql.Stmt, error) {
	return x.DB().Prepare(query)
}

func (x *Xorm) TraceOff() {
	x.CloseLog()
}

type Balancer struct {
	*Xorm
	*ormx.Balancer
}

func (b *Balancer) Master() *Xorm {
	return b.Balancer.Master().(*Xorm)
}

func (b *Balancer) Replica() *Xorm {
	return b.Balancer.Replica().(*Xorm)
}

func (b *Balancer) QueryStr(sql string, args ...interface{}) []map[string]string {
	return b.Replica().QueryStr(sql, args...)
}

func (b *Balancer) QueryRaw(sql string, args ...interface{}) []map[string]interface{} {
	return b.Replica().QueryRaw(sql, args...)
}

// =======================
// 原生SQL查询
// =======================
func (b *Balancer) RawQuery(sql string, args ...interface{}) ([]*xorm.ResultSet, error) {
	return b.Replica().RawQuery(sql, args...)
}

func (b *Balancer) RawQueryCallback(callback func(*core.Rows, []string), sql string, args ...interface{}) error {
	return b.Replica().RawQueryCallback(callback, sql, args...)
}

// RawQueryKv 查询键值对
func (b *Balancer) RawQueryKv(key string, val string, sql string, args ...interface{}) map[string]string {
	return b.Replica().RawQueryKv(key, val, sql, args...)
}

func (b *Balancer) RawQueryAllKvs(key string, sql string, args ...interface{}) map[string][]map[string]string {
	return b.Replica().RawQueryAllKvs(key, sql, args...)
}

// -----------------------
// ResultSet结果
// -----------------------
func (b *Balancer) GetRows(sql string, params ...interface{}) []*xorm.ResultSet {
	return b.Replica().GetRows(sql, params...)
}

func (b *Balancer) GetRow(sql string, params ...interface{}) *xorm.ResultSet {
	return b.Replica().GetRow(sql, params...)
}

func (b *Balancer) GetOne(sql string, params ...interface{}) string {
	return b.Replica().GetOne(sql, params...)
}

// RawSelect .
// RawSelect("*","member","id=?",1)
// RawSelect("*","member","status=? AND sex=?",1,1)
// RawSelect("*","`~member` a,`~order` b","a.status=? AND b.status=?",1,1)
func (b *Balancer) RawSelect(fields string, table string, where string, params ...interface{}) xorm.SelectRows {
	return b.Replica().RawSelect(fields, table, where, params...)
}

// -----------------------
// map结果
// -----------------------
func (b *Balancer) RawFetchAll(fields string, table string, where string, params ...interface{}) []map[string]string {
	return b.Replica().RawFetchAll(fields, table, where, params...)
}

func (b *Balancer) RawFetch(fields string, table string, where string, params ...interface{}) map[string]string {
	return b.Replica().RawFetch(fields, table, where, params...)
}

// RawQueryKvs 查询基于指定字段值为键名的map
func (b *Balancer) RawQueryKvs(key string, sql string, args ...interface{}) map[string]map[string]string {
	return b.Replica().RawQueryKvs(key, sql, args...)
}

// RawQueryStr 查询[]map[string]string
func (b *Balancer) RawQueryStr(sql string, args ...interface{}) []map[string]string {
	return b.Replica().RawQueryStr(sql, args...)
}
