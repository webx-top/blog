package xorm

import (
	"fmt"
	"reflect"

	"github.com/coscms/xorm/core"
)

// == ORDER BY ==
type orderBy []*orderByParam

type orderByParam struct {
	Field string //Column name
	Sort  string // ASC/DESC
}

// == Fields ==
type fields []string //for Omit

func newJoinTables(stmt *Statement) *joinTables {
	return &joinTables{
		params:    []*joinParam{},
		statement: stmt,
	}
}

// == JOIN ==
type joinTables struct {
	params    []*joinParam
	statement *Statement
}

func (j *joinTables) New(stmt *Statement) *joinParam {
	join := NewJoinParam(stmt)
	j.params = append(j.params, join)
	return join
}

func (j *joinTables) Add(join *joinParam) {
	j.params = append(j.params, join)
}

func (j *joinTables) Size() int {
	return len(j.params)
}

func (j *joinTables) HasJoin() bool {
	return j.Size() > 0
}

func (j *joinTables) String() string {
	var joinStr string
	if j.Size() > 0 {
		var t string
		for _, join := range j.params {
			joinStr += t + join.String()
			t = ` `
		}
	} else {
		joinStr = j.fromRelation()
	}
	return joinStr
}

func (j *joinTables) fromRelation() string {
	r := j.statement.relation
	if r == nil || r.IsTable {
		return ``
	}
	var (
		joinStr string
		t       string
	)
	for i, table := range r.Extends {
		rt := r.RelTables[i]
		if rt == nil || !rt.IsValid() {
			continue
		}
		s := rt.JoinType + ` JOIN ` + j.statement.Engine.Quote(rt.TableName)
		alias, _ := r.ExAlias[table.Name]
		if len(alias) > 0 {
			s += ` AS ` + j.statement.Engine.Quote(alias)
		}
		s += ` ON ` + rt.Where
		joinStr += t + s
		t = ` `
	}
	return joinStr
}

func (j *joinTables) Args() (args []interface{}) {
	for _, join := range j.params {
		args = append(args, join.Args...)
	}
	return
}

func NewJoinParam(stmt *Statement) *joinParam {
	return &joinParam{
		Args:      make([]interface{}, 0),
		statement: stmt,
	}
}

type joinParam struct {
	Operator string //LEFT/RIGHT/INNER...
	Table    string
	Alias    string
	ONStr    string
	Args     []interface{}
	SQLStr   string

	statement *Statement
}

func (j *joinParam) String() string {
	if len(j.SQLStr) > 0 {
		return j.SQLStr
	}
	joinStr := j.Operator + ` JOIN ` + j.statement.Engine.Quote(j.Table)
	if len(j.Alias) == 0 && j.statement.relation != nil {
		j.Alias, _ = j.statement.relation.ExAlias[j.Table]
	}
	if len(j.Alias) > 0 {
		joinStr += ` AS ` + j.statement.Engine.Quote(j.Alias)
	}
	if len(j.ONStr) > 0 {
		joinStr += ` ON ` + j.ONStr
	}
	j.SQLStr = joinStr
	return joinStr
}

// == Extends Statement ==

// Join The joinOP should be one of INNER, LEFT OUTER, CROSS etc - this will be prepended to JOIN
func (statement *Statement) join(joinOP string, tablename interface{}, condition string, args ...interface{}) *Statement {
	join := statement.joinTables.New(statement)
	join.Operator = joinOP
	join.Table = ``
	join.Alias = ``
	join.ONStr = condition
	join.Args = args

	switch t := tablename.(type) {
	case []string:
		if len(t) > 1 {
			join.Table = statement.withPrefix(t[0])
			join.Alias = t[1]
		} else if len(t) == 1 {
			join.Table = statement.withPrefix(t[0])
		}
	case []interface{}:
		l := len(t)
		var table string
		if l > 0 {
			f := t[0]
			v := rValue(f)
			t := v.Type()
			if t.Kind() == reflect.String {
				table = f.(string)
			} else if t.Kind() == reflect.Struct {
				r := statement.Engine.autoMapType(v)
				table = r.Name
			} else {
				table = fmt.Sprintf("%v", f)
			}
			join.Table = statement.withPrefix(table)
		}
		if l > 1 {
			join.Alias = fmt.Sprintf("%v", t[1])
		}
	case core.SQL:
		join.SQLStr = joinOP + ` JOIN ` + string(t)
	case string:
		join.Table = statement.withPrefix(t)
	default:
		v := rValue(tablename)
		typ := v.Type()
		if typ.Kind() == reflect.Struct {
			r := statement.Engine.autoMapType(v)
			join.Table = r.Name
		} else {
			join.Table = statement.withPrefix(fmt.Sprintf("%v", tablename))
		}
	}
	return statement
}

func (statement *Statement) withPrefix(tableName string) string {
	if len(tableName) > 0 && tableName[0] == '~' {
		return statement.Engine.TablePrefix + tableName[1:]
	}
	return tableName
}

func (statement *Statement) JoinStr() string {
	if statement.joinGenerated {
		return statement.joinStr
	}
	statement.joinStr = statement.joinTables.String()
	statement.joinArgs = statement.joinTables.Args()
	statement.joinGenerated = true
	return statement.joinStr
}

func (statement *Statement) SetRelation(r *core.Relation) {
	statement.relation = r
	if r == nil || len(r.Extends) < 1 {
		return
	}
	if !r.IsTable {
		statement.RefTable = r.Extends[0]
		if len(statement.TableAlias) == 0 {
			name := statement.RefTable.Name
			statement.TableAlias, _ = r.ExAlias[name]
		}
		statement.tableName = statement.RefTable.Name
	} else {
		if len(statement.TableAlias) == 0 {
			statement.TableAlias = r.Table.Type.Name()
		}
	}
}
