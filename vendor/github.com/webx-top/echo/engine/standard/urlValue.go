package standard

import (
	"net/url"
)

type UrlValue struct {
	values *url.Values
	initFn func() *url.Values
}

func (u *UrlValue) Add(key string, value string) {
	u.init()
	u.values.Add(key, value)
}

func (u *UrlValue) Del(key string) {
	u.init()
	u.values.Del(key)
}

func (u *UrlValue) Get(key string) string {
	u.init()
	return u.values.Get(key)
}

func (u *UrlValue) Gets(key string) []string {
	u.init()
	if v, ok := (*u.values)[key]; ok {
		return v
	}
	return []string{}
}

func (u *UrlValue) Set(key string, value string) {
	u.init()
	u.values.Set(key, value)
}

func (u *UrlValue) Encode() string {
	u.init()
	return u.values.Encode()
}

func (u *UrlValue) All() map[string][]string {
	u.init()
	return *u.values
}

func (u *UrlValue) Reset(data url.Values) {
	u.values = &data
}

func (u *UrlValue) init() {
	if u.values != nil {
		return
	}
	u.values = u.initFn()
}

func NewValue(r *Request) *Value {
	v := &Value{
		queryArgs: &UrlValue{initFn: func() *url.Values {
			q := r.url.Query()
			return &q
		}},
		request: r,
	}
	v.postArgs = &UrlValue{initFn: func() *url.Values {
		r.MultipartForm()
		return &r.request.PostForm
	}}
	return v
}

type Value struct {
	request   *Request
	queryArgs *UrlValue
	postArgs  *UrlValue
	form      *url.Values
}

func (v *Value) Add(key string, value string) {
	v.init()
	v.form.Add(key, value)
}

func (v *Value) Del(key string) {
	v.init()
	v.form.Del(key)
}

func (v *Value) Get(key string) string {
	v.init()
	return v.form.Get(key)
}

func (v *Value) Gets(key string) []string {
	v.init()
	form := *v.form
	if v, ok := form[key]; ok {
		return v
	}
	return []string{}
}

func (v *Value) Set(key string, value string) {
	v.init()
	v.form.Set(key, value)
}

func (v *Value) Encode() string {
	v.init()
	return v.form.Encode()
}

func (v *Value) init() {
	if v.form != nil {
		return
	}
	//v.request.request.ParseForm()
	v.request.MultipartForm()
	v.form = &v.request.request.Form
}

func (v *Value) All() map[string][]string {
	v.init()
	return *v.form
}

func (v *Value) Reset(data url.Values) {
	v.form = &data
}
