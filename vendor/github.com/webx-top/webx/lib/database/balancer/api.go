package balancer

import (
	"database/sql"

	"github.com/admpub/log"
	"github.com/coscms/xorm"
)

// Stmt is an aggregate prepared statement.
// It holds a prepared statement for each underlying physical db.
type Stmt interface {
	Close() error
	Exec(...interface{}) (sql.Result, error)
	Query(...interface{}) (*sql.Rows, error)
	QueryRow(...interface{}) *sql.Row
}

// Prepare creates a prepared statement for later queries or executions on each physical database.
// Multiple queries or executions may be run concurrently from the returned statement.
// This is equivalent to running: Prepare() using database/sql
func (b *Balancer) Prepare(query string) (Stmt, error) {
	dbs := b.GetAllDbs()
	stmts := make([]*sql.Stmt, len(dbs))
	for i := range stmts {
		s, err := dbs[i].DB().Prepare(query)
		if err != nil {
			return nil, err
		}
		stmts[i] = s.Stmt
	}
	return &stmt{bl: b, stmts: stmts}, nil
}

func (b *Balancer) TraceOn(prefix string, logger *log.Logger) {
	if len(prefix) > 0 {
		prefix += ` `
	}
	l := xorm.NewAdmpubLoggerWithPrefix(prefix+"<slave>", logger)
	for _, s := range b.replicas {
		s.SetLogger(l)
	}
	b.Engine.SetLogger(xorm.NewAdmpubLoggerWithPrefix(prefix+"<master>", logger.GetLogger(logger.Category)))
}

func (b *Balancer) TraceOff() {
	for _, db := range b.GetAllDbs() {
		db.CloseLog()
	}
}
