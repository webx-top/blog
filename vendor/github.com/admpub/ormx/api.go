package ormx

import (
	"database/sql"

	"github.com/admpub/ormx/orm"
)

// Prepare creates a prepared statement for later queries or executions on each physical database.
// Multiple queries or executions may be run concurrently from the returned statement.
// This is equivalent to running: Prepare() using database/sql
func (b *Balancer) Prepare(query string) (orm.Stmt, error) {
	dbs := b.GetAllDbs()
	stmts := make([]*sql.Stmt, len(dbs))
	for i := range stmts {
		s, err := dbs[i].Prepare(query)
		if err != nil {
			return nil, err
		}
		stmts[i] = s
	}
	return &stmt{bl: b, stmts: stmts}, nil
}

func (b *Balancer) TraceOn(prefix string, logger interface{}) {
	if len(prefix) > 0 {
		prefix += ` `
	}
	for _, s := range b.replicas {
		s.TraceOn(prefix+"<slave>", logger)
	}
	b.ORM.TraceOn(prefix+"<master>", logger)
}

func (b *Balancer) TraceOff() {
	for _, db := range b.GetAllDbs() {
		db.TraceOff()
	}
}
