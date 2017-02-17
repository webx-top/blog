package orm

import "database/sql"

type ORM interface {
	DB() *sql.DB
	Prepare(query string) (*sql.Stmt, error)
	TraceOn(prefix string, logger interface{})
	TraceOff()
}

type Logger interface {
	Printf(format string, v ...interface{})
}

// Stmt is an aggregate prepared statement.
// It holds a prepared statement for each underlying physical db.
type Stmt interface {
	Close() error
	Exec(...interface{}) (sql.Result, error)
	Query(...interface{}) (*sql.Rows, error)
	QueryRow(...interface{}) *sql.Row
}

type Connector func(driverName string, dsn string) (ORM, error)
