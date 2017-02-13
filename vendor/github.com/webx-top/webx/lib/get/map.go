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
package get

import (
	"strconv"
)

type IMap map[string]interface{}

func (m *IMap) IMap(k string) map[string]interface{} {
	if v, y := (*m)[k]; y {
		if v, y := v.(map[string]interface{}); y {
			return v
		}
	}
	return map[string]interface{}{}
}

func (m *IMap) SMap(k string) map[string]string {
	if v, y := (*m)[k]; y {
		if v, y := v.(map[string]string); y {
			return v
		}
	}
	return map[string]string{}
}

func (m *IMap) SSlice(k string) []string {
	if v, y := (*m)[k]; y {
		if v, y := v.([]string); y {
			return v
		}
	}
	return []string{}
}

func (m *IMap) ISlice(k string) []interface{} {
	if v, y := (*m)[k]; y {
		if v, y := v.([]interface{}); y {
			return v
		}
	}
	return make([]interface{}, 0)
}

func (m *IMap) BSlice(k string) []byte {
	if v, y := (*m)[k]; y {
		if v, y := v.([]byte); y {
			return v
		}
	}
	return make([]byte, 0)
}

func (m *IMap) Bytes(k string) []byte {
	return m.BSlice(k)
}

func (m *IMap) String(k string) string {
	if v, y := (*m)[k]; y {
		if v, y := v.(string); y {
			return v
		}
	}
	return ""
}

func (m *IMap) Interface(k string) interface{} {
	if v, y := (*m)[k]; y {
		return v
	}
	return nil
}

func (m *IMap) Float32(k string) float32 {
	if v, y := (*m)[k]; y {
		if v, y := v.(float32); y {
			return v
		}
	}
	return 0
}

func (m *IMap) Float64(k string) float64 {
	if v, y := (*m)[k]; y {
		if v, y := v.(float64); y {
			return v
		}
	}
	return 0
}

func (m *IMap) Int(k string) int {
	if v, y := (*m)[k]; y {
		if v, y := v.(int); y {
			return v
		}
	}
	return 0
}

func (m *IMap) Uint(k string) uint {
	if v, y := (*m)[k]; y {
		if v, y := v.(uint); y {
			return v
		}
	}
	return 0
}

func (m *IMap) Int8(k string) int8 {
	if v, y := (*m)[k]; y {
		if v, y := v.(int8); y {
			return v
		}
	}
	return 0
}

func (m *IMap) Uint8(k string) uint8 {
	if v, y := (*m)[k]; y {
		if v, y := v.(uint8); y {
			return v
		}
	}
	return 0
}

func (m *IMap) Int16(k string) int16 {
	if v, y := (*m)[k]; y {
		if v, y := v.(int16); y {
			return v
		}
	}
	return 0
}

func (m *IMap) Uint16(k string) uint16 {
	if v, y := (*m)[k]; y {
		if v, y := v.(uint16); y {
			return v
		}
	}
	return 0
}

func (m *IMap) Int32(k string) int32 {
	if v, y := (*m)[k]; y {
		if v, y := v.(int32); y {
			return v
		}
	}
	return 0
}

func (m *IMap) Uint32(k string) uint32 {
	if v, y := (*m)[k]; y {
		if v, y := v.(uint32); y {
			return v
		}
	}
	return 0
}

func (m *IMap) Int64(k string) int64 {
	if v, y := (*m)[k]; y {
		if v, y := v.(int64); y {
			return v
		}
	}
	return 0
}

func (m *IMap) Uint64(k string) uint64 {
	if v, y := (*m)[k]; y {
		if v, y := v.(uint64); y {
			return v
		}
	}
	return 0
}

func (m *IMap) Bool(k string) bool {
	if v, y := (*m)[k]; y {
		if v, y := v.(bool); y {
			return v
		}
	}
	return false
}

type ToIMap struct {
	Data map[string]string
	Type map[string]string
}

func (t *ToIMap) Interface() map[string]interface{} {
	return SMapToIMap(t.Data, t.Type)
}

type SMap map[string]string

func (t *SMap) ToIMap(typ map[string]string) map[string]interface{} {
	return SMapToIMap(*t, typ)
}

//map[string]string => map[string]interface{}
func SMapToIMap(data map[string]string, typs map[string]string) map[string]interface{} {
	r := make(map[string]interface{})
	for k, v := range data {
		typ, _ := typs[k]
		r[k] = ParseString(v, typ)
	}
	return r
}

func ParseString(val string, typ string) interface{} {
	switch typ {
	case "int8":
		if v, err := strconv.ParseInt(val, 10, 8); err == nil {
			return int8(v)
		}
		return 0

	case "uint8":
		if v, err := strconv.ParseInt(val, 10, 8); err == nil {
			return uint8(v)
		}
		return 0

	case "int16":
		if v, err := strconv.ParseInt(val, 10, 16); err == nil {
			return int16(v)
		}
		return 0

	case "uint16":
		if v, err := strconv.ParseInt(val, 10, 16); err == nil {
			return uint16(v)
		}
		return 0

	case "int":
		if v, err := strconv.Atoi(val); err == nil {
			return v
		}
		return 0

	case "uint":
		if v, err := strconv.ParseInt(val, 10, 32); err == nil {
			return uint(v)
		}
		return 0

	case "int32":
		if v, err := strconv.ParseInt(val, 10, 32); err == nil {
			return int32(v)
		}
		return 0

	case "uint32":
		if v, err := strconv.ParseInt(val, 10, 32); err == nil {
			return uint32(v)
		}
		return 0

	case "int64":
		if v, err := strconv.ParseInt(val, 10, 64); err == nil {
			return int64(v)
		}
		return 0

	case "uint64":
		if v, err := strconv.ParseInt(val, 10, 64); err == nil {
			return uint64(v)
		}
		return 0

	case "float32":
		if v, err := strconv.ParseFloat(val, 32); err == nil {
			return float32(v)
		}
		return 0

	case "float64":
		if v, err := strconv.ParseFloat(val, 64); err == nil {
			return float64(v)
		}
		return 0

	case "bool":
		if v, err := strconv.ParseBool(val); err == nil {
			return v
		}
		return false

	default:
		return val
	}
}
