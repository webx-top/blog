package model

import (
	//"errors"
	//"strings"

	D "github.com/webx-top/blog/app/base/dbschema"
	X "github.com/webx-top/webx"
	//"github.com/webx-top/webx/lib/com"
)

func NewPost(ctx *X.Context) *Post {
	return &Post{M: NewM(ctx)}
}

type Post struct {
	*M
}

func (a *Post) List(s *Select) (countFn func() int64, m []*D.Post, err error) {
	m = []*D.Post{}
	err = s.Do().Find(&m)
	if err != nil {
		return
	}
	countFn = func() int64 {
		return s.Count(D.Post{})
	}
	return
}

func (a *Post) Add(m *D.Post) (err error) {
	_, err = a.DB.Insert(m)
	return
}
