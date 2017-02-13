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

type Closer interface {
	//关闭缓存引擎
	Close() error
}

type Storer interface {
	//保存到缓存
	Put(key string, value interface{}) error

	//删除缓存
	Del(key string) error

	//获取缓存
	Get(key string) (interface{}, error)
}
