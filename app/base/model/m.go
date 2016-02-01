package model

import (
	"github.com/webx-top/blog/app/base"
	"github.com/webx-top/blog/app/base/lib/database"
	"github.com/webx-top/webx/lib/i18n"
)

func NewM(lang string) *M {
	return &M{
		DB:   base.DB,
		Lang: lang,
	}
}

type M struct {
	DB   *database.Orm
	Lang string
}

func (m *M) T(key string, args ...interface{}) string {
	return i18n.T(m.Lang, key, args...)
}
