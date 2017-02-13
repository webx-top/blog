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
package cache

import (
	"errors"
)

var ErrNoData = errors.New(`no data`)

var caches = make(map[string]func() Storer)
var defaults = &defaultCache{}

func Reg(name string, c func() Storer) {
	caches[name] = c
}

func Get(name string) Storer {
	fn, _ := caches[name]
	if fn == nil {
		fn = func() Storer {
			return defaults
		}
	}
	return fn()
}

func Has(name string) bool {
	_, ok := caches[name]
	return ok
}

func Del(name string) {
	if _, ok := caches[name]; ok {
		delete(caches, name)
	}
}

type defaultCache struct {
}

func (s *defaultCache) Put(key string, value interface{}) error {
	return nil
}

func (s *defaultCache) Get(key string) (interface{}, error) {
	return nil, ErrNoData
}

func (s *defaultCache) Del(key string) error {
	return nil
}
