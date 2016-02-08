package model

import (
	"fmt"

	"github.com/coscms/xorm"
	"github.com/webx-top/blog/app/base"
	"github.com/webx-top/blog/app/base/lib/database"
	"github.com/webx-top/blog/app/base/lib/datatable"
	X "github.com/webx-top/webx"
	"github.com/webx-top/webx/lib/i18n"
)

func NewM(ctx *X.Context) *M {
	return &M{
		DB:      base.DB,
		Context: ctx,
	}
}

type M struct {
	DB      *database.Orm
	Context *X.Context
}

func (this *M) T(key string, args ...interface{}) string {
	return i18n.T(this.Context.Language, key, args...)
}

// =====================================
// TransManager
// =====================================
func (this *M) Begin() *xorm.Session {
	ss, ok := this.transSession()
	if ok {
		ss.Close()
	}
	ss = this.DB.NewSession()
	err := ss.Begin()
	if err != nil {
		this.Context.X().Echo().Logger().Error(err)
	}
	this.Context.Set(`webx:transSession`, ss)
	return ss
}

//事务是否已经开始
func (this *M) HasBegun() bool {
	ss, ok := this.transSession()
	return ok && ss != nil
}

func (this *M) transSession() (ss *xorm.Session, ok bool) {
	ss, ok = this.Context.Get(`webx:transSession`).(*xorm.Session)
	return
}

func (this *M) TSess() *xorm.Session { // TransSession
	ss, ok := this.transSession()
	if ok == false || ss == nil {
		return this.Begin()
	}
	return ss
}

func (this *M) Trans(fn func(*xorm.Session) error) *database.Orm {
	ss, ok := this.transSession()
	begun := ok && ss != nil
	if !begun {
		ss = this.Begin()
	}
	result := fn(ss)
	if !begun {
		this.End(result == nil, ss)
	}
	return this.DB
}

func (this *M) Sess() *xorm.Session { // TransSession or Session
	ss, ok := this.transSession()
	if ok == false {
		var ss *xorm.Session = this.DB.NewSession()
		ss.IsAutoClose = true
		return ss
	}
	return ss
}

func (this *M) End(result bool, args ...*xorm.Session) (err error) {
	var ss *xorm.Session
	if len(args) > 0 && args[0] != nil {
		ss = args[0]
	} else {
		ss, _ = this.transSession()
	}
	if result {
		err = ss.Commit()
	} else {
		err = ss.Rollback()
	}
	if err != nil {
		this.Context.X().Echo().Logger().Error(err)
	}
	ss.Close()
	this.Context.Set(`webx:transSession`, nil)
	return
}

func (this *M) NewSelect() *Select {
	s := NewSelect(this.DB)
	return s
}

func NewSelect(orm *database.Orm) *Select {
	s := &Select{}
	s.Orm = orm
	s.Params = make([]interface{}, 0)
	return s
}

type Select struct {
	Offset    int64
	Limit     int64
	OrderBy   string
	Condition string
	Params    []interface{}
	GroupBy   string
	Having    string
	Table     interface{}
	Alias     string
	*database.Orm
}

func (a *Select) Do(args ...interface{}) *xorm.Session {
	return a.GenSS(args...).OrderBy(a.OrderBy).Limit(int(a.Limit), int(a.Offset))
}

func (a *Select) AddP(args ...interface{}) *Select {
	a.Params = append(a.Params, args...)
	return a
}

func (a *Select) FromDT(dt *datatable.DataTable) *Select {
	a.OrderBy = dt.OrderBy
	a.Offset = dt.Offset
	a.Limit = dt.PageSize
	return a
}

func (a *Select) GenSS(args ...interface{}) *xorm.Session {
	var s *xorm.Session = a.Orm.NewSession()
	s.IsAutoClose = true
	switch len(args) {
	case 2:
		alias, _ := args[1].(string)
		if args[0] == nil {
			s = s.Alias(alias)
		} else {
			s = s.Table(args[0]).Alias(alias)
			a.Table = args[0]
		}
		a.Alias = alias
	case 1:
		s = s.Table(args[0])
		a.Table = args[0]
	default:
		if a.Table != nil {
			s = s.Table(a.Table)
		}
		if a.Alias != `` {
			s = s.Alias(a.Alias)
		}
	}
	return s.Where(a.Condition, a.Params...).GroupBy(a.GroupBy).Having(a.Having)
}

func (a *Select) Count(m interface{}) int64 {
	count, err := a.GenSS().Count(m)
	if err != nil {
		fmt.Println(err)
	}
	return count
}
