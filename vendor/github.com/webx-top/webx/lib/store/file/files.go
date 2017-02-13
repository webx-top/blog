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
package file

import (
	"errors"
	"io"
)

var ErrNoData = errors.New(`no data`)

var files = make(map[string]func() Storer)
var defaults = &defaultStore{}

func Reg(name string, c func() Storer) {
	files[name] = c
}

func Get(name string) Storer {
	fn, _ := files[name]
	if fn == nil {
		fn = func() Storer {
			return defaults
		}
	}
	return fn()
}

func Has(name string) bool {
	_, ok := files[name]
	return ok
}

func Del(name string) {
	if _, ok := files[name]; ok {
		delete(files, name)
	}
}

type defaultStore struct {
}

func (s *defaultStore) Put(body io.ReadCloser, fileName string) (*Result, error) {
	return nil, nil
}

func (s *defaultStore) Get(key string) (io.ReadCloser, error) {
	return nil, ErrNoData
}

func (s *defaultStore) Del(key string) error {
	return nil
}

func (s *defaultStore) Open() error {
	return nil
}

func (s *defaultStore) Close() error {
	return nil
}

func (s *defaultStore) EName() string {
	return ``
}
