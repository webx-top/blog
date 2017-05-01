package xorm

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/coscms/xorm/core"
)

func cleanupProcessorsClosures(slices *[]func(interface{})) {
	if len(*slices) > 0 {
		*slices = make([]func(interface{}), 0)
	}
}

func (session *Session) scanMapIntoStruct(obj interface{}, objMap map[string][]byte) error {
	dataStruct := rValue(obj)
	if dataStruct.Kind() != reflect.Struct {
		return errors.New("Expected a pointer to a struct")
	}

	var col *core.Column
	session.Statement.setRefValue(dataStruct)
	table := session.Statement.RefTable
	tableName := session.Statement.tableName

	for key, data := range objMap {
		if col = table.GetColumn(key); col == nil {
			session.Engine.logger.Warnf("struct %v's has not field %v. %v",
				table.Type.Name(), key, table.ColumnsSeq())
			continue
		}

		fieldName := col.FieldName
		fieldPath := strings.Split(fieldName, ".")
		var fieldValue reflect.Value
		if len(fieldPath) > 2 {
			session.Engine.logger.Error("Unsupported mutliderive", fieldName)
			continue
		} else if len(fieldPath) == 2 {
			parentField := dataStruct.FieldByName(fieldPath[0])
			if parentField.IsValid() {
				fieldValue = parentField.FieldByName(fieldPath[1])
			}
		} else {
			fieldValue = dataStruct.FieldByName(fieldName)
		}
		if !fieldValue.IsValid() || !fieldValue.CanSet() {
			session.Engine.logger.Warnf("table %v's column %v is not valid or cannot set", tableName, key)
			continue
		}

		err := session.bytes2Value(col, &fieldValue, data)
		if err != nil {
			return err
		}
	}

	return nil
}

// Cell cell is a result of one column field
type Cell *interface{}

func (session *Session) rows2Beans(rows *core.Rows, fields []string, fieldsCount int,
	table *core.Table, newElemFunc func() reflect.Value,
	sliceValueSetFunc func(*reflect.Value)) error {
	for rows.Next() {
		var newValue = newElemFunc()
		bean := newValue.Interface()
		dataStruct := rValue(bean)
		err := session._row2Bean(rows, fields, fieldsCount, bean, &dataStruct, table)
		if err != nil {
			return err
		}
		sliceValueSetFunc(&newValue)
	}
	return nil
}

func (session *Session) row2Bean(rows *core.Rows, fields []string, fieldsCount int, bean interface{}) error {
	dataStruct := rValue(bean)
	if dataStruct.Kind() != reflect.Struct {
		return errors.New("Expected a pointer to a struct")
	}

	session.Statement.setRefValue(dataStruct)

	return session._row2Bean(rows, fields, fieldsCount, bean, &dataStruct, session.Engine.autoMapType(dataStruct)) //[SWH|M]fixbug
}

func (session *Session) _row2Bean(rows *core.Rows, fields []string, fieldsCount int, bean interface{}, dataStruct *reflect.Value, table *core.Table) error {
	//定义保存数据库单行数据结果的变量
	scanResults := make([]interface{}, fieldsCount)
	for i := 0; i < fieldsCount; i++ {
		var cell interface{}
		scanResults[i] = &cell
	}

	//从数据库读取数据到scanResults
	if err := rows.Scan(scanResults...); err != nil {
		return err
	}

	//允许在把数据库结果赋予结构体之前进行一些自定义操作
	if b, hasBeforeSet := bean.(BeforeSetProcessor); hasBeforeSet {
		for ii, key := range fields {
			b.BeforeSet(key, Cell(scanResults[ii].(*interface{})))
		}
	}

	//允许在把数据库结果赋予结构体之后进行一些自定义操作
	defer func() {
		if b, hasAfterSet := bean.(AfterSetProcessor); hasAfterSet {
			for ii, key := range fields {
				b.AfterSet(key, Cell(scanResults[ii].(*interface{})))
			}
		}
	}()

	var tempMap = make(map[string]int)

	//遍历数据库单行中的各个字段，并将其值赋予相应的结构体属性中
	for ii, key := range fields {
		lKey := strings.ToLower(key) //小写表字段名
		idx, ok := tempMap[lKey]
		if !ok {
			idx = 0
		} else {
			idx = idx + 1
		}
		tempMap[lKey] = idx

		if fieldValue := session.getField(dataStruct, key, table, idx); fieldValue != nil {
			rawValue := reflect.Indirect(reflect.ValueOf(scanResults[ii]))

			// if row is null then ignore
			if rawValue.Interface() == nil {
				continue
			}

			if fieldValue.CanAddr() {
				if structConvert, ok := fieldValue.Addr().Interface().(core.Conversion); ok {
					if data, err := value2Bytes(&rawValue); err == nil {
						structConvert.FromDB(data)
					} else {
						session.Engine.logger.Error(err)
					}
					continue
				}
			}

			if _, ok := fieldValue.Interface().(core.Conversion); ok {
				if data, err := value2Bytes(&rawValue); err == nil {
					if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
						fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
					}
					fieldValue.Interface().(core.Conversion).FromDB(data)
				} else {
					session.Engine.logger.Error(err)
				}
				continue
			}

			rawValueType := reflect.TypeOf(rawValue.Interface())
			vv := reflect.ValueOf(rawValue.Interface())

			fieldType := fieldValue.Type()
			hasAssigned := false
			col := table.GetColumnIdx(key, idx)

			if col.SQLType.IsJson() {
				var bs []byte
				if rawValueType.Kind() == reflect.String {
					bs = []byte(vv.String())
				} else if rawValueType.ConvertibleTo(core.BytesType) {
					bs = vv.Bytes()
				} else {
					return fmt.Errorf("unsupported database data type: %s %v", key, rawValueType.Kind())
				}

				hasAssigned = true

				if len(bs) > 0 {
					if fieldValue.CanAddr() {
						err := json.Unmarshal(bs, fieldValue.Addr().Interface())
						if err != nil {
							session.Engine.logger.Error(key, err)
							return err
						}
					} else {
						x := reflect.New(fieldType)
						err := json.Unmarshal(bs, x.Interface())
						if err != nil {
							session.Engine.logger.Error(key, err)
							return err
						}
						fieldValue.Set(x.Elem())
					}
				}

				continue
			}

			switch fieldType.Kind() {
			case reflect.Complex64, reflect.Complex128:
				// TODO: reimplement this
				var bs []byte
				if rawValueType.Kind() == reflect.String {
					bs = []byte(vv.String())
				} else if rawValueType.ConvertibleTo(core.BytesType) {
					bs = vv.Bytes()
				}

				hasAssigned = true
				if len(bs) > 0 {
					if fieldValue.CanAddr() {
						err := json.Unmarshal(bs, fieldValue.Addr().Interface())
						if err != nil {
							session.Engine.logger.Error(err)
							return err
						}
					} else {
						x := reflect.New(fieldType)
						err := json.Unmarshal(bs, x.Interface())
						if err != nil {
							session.Engine.logger.Error(err)
							return err
						}
						fieldValue.Set(x.Elem())
					}
				}
			case reflect.Slice, reflect.Array:
				switch rawValueType.Kind() {
				case reflect.Slice, reflect.Array:
					switch rawValueType.Elem().Kind() {
					case reflect.Uint8:
						if fieldType.Elem().Kind() == reflect.Uint8 {
							hasAssigned = true
							fieldValue.Set(vv)
						}
					}
				}
			case reflect.String:
				if rawValueType.Kind() == reflect.String {
					hasAssigned = true
					fieldValue.SetString(vv.String())
				}
			case reflect.Bool:
				if rawValueType.Kind() == reflect.Bool {
					hasAssigned = true
					fieldValue.SetBool(vv.Bool())
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				switch rawValueType.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					hasAssigned = true
					fieldValue.SetInt(vv.Int())
				}
			case reflect.Float32, reflect.Float64:
				switch rawValueType.Kind() {
				case reflect.Float32, reflect.Float64:
					hasAssigned = true
					fieldValue.SetFloat(vv.Float())
				}
			case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
				switch rawValueType.Kind() {
				case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
					hasAssigned = true
					fieldValue.SetUint(vv.Uint())
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					hasAssigned = true
					fieldValue.SetUint(uint64(vv.Int()))
				}
			case reflect.Struct:
				if fieldType.ConvertibleTo(core.TimeType) {
					if rawValueType == core.TimeType {
						hasAssigned = true

						t := vv.Convert(core.TimeType).Interface().(time.Time)
						z, _ := t.Zone()
						if len(z) == 0 || t.Year() == 0 { // !nashtsai! HACK tmp work around for lib/pq doesn't properly time with location
							dbTZ := session.Engine.DatabaseTZ
							if dbTZ == nil {
								dbTZ = time.Local
							}
							session.Engine.logger.Debugf("empty zone key[%v] : %v | zone: %v | location: %+v\n", key, t, z, *t.Location())
							t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(),
								t.Minute(), t.Second(), t.Nanosecond(), dbTZ)
						}
						// !nashtsai! convert to engine location
						if col.TimeZone == nil {
							t = t.In(session.Engine.TZLocation)
						} else {
							t = t.In(col.TimeZone)
						}
						fieldValue.Set(reflect.ValueOf(t).Convert(fieldType))

						// t = fieldValue.Interface().(time.Time)
						// z, _ = t.Zone()
						// session.Engine.LogDebug("fieldValue key[%v]: %v | zone: %v | location: %+v\n", key, t, z, *t.Location())
					} else if rawValueType == core.IntType || rawValueType == core.Int64Type ||
						rawValueType == core.Int32Type {
						hasAssigned = true
						var tz *time.Location
						if col.TimeZone == nil {
							tz = session.Engine.TZLocation
						} else {
							tz = col.TimeZone
						}
						t := time.Unix(vv.Int(), 0).In(tz)
						//vv = reflect.ValueOf(t)
						fieldValue.Set(reflect.ValueOf(t).Convert(fieldType))
					} else {
						if d, ok := vv.Interface().([]uint8); ok {
							hasAssigned = true
							t, err := session.byte2Time(col, d)
							if err != nil {
								session.Engine.logger.Error("byte2Time error:", err.Error())
								hasAssigned = false
							} else {
								fieldValue.Set(reflect.ValueOf(t).Convert(fieldType))
							}
						} else if d, ok := vv.Interface().(string); ok {
							hasAssigned = true
							t, err := session.str2Time(col, d)
							if err != nil {
								session.Engine.logger.Error("byte2Time error:", err.Error())
								hasAssigned = false
							} else {
								fieldValue.Set(reflect.ValueOf(t).Convert(fieldType))
							}
						} else {
							panic(fmt.Sprintf("rawValueType is %v, value is %v", rawValueType, vv.Interface()))
						}
					}
				} else if nulVal, ok := fieldValue.Addr().Interface().(sql.Scanner); ok {
					// !<winxxp>! 增加支持sql.Scanner接口的结构，如sql.NullString
					hasAssigned = true
					if err := nulVal.Scan(vv.Interface()); err != nil {
						session.Engine.logger.Error("sql.Sanner error:", err.Error())
						hasAssigned = false
					}
				} else if col.SQLType.IsJson() {
					if rawValueType.Kind() == reflect.String {
						hasAssigned = true
						x := reflect.New(fieldType)
						b := []byte(vv.String())
						if len(b) > 0 {
							err := json.Unmarshal(b, x.Interface())
							if err != nil {
								session.Engine.logger.Error(err)
								return err
							}
							fieldValue.Set(x.Elem())
						}
					} else if rawValueType.Kind() == reflect.Slice {
						hasAssigned = true
						x := reflect.New(fieldType)
						b := vv.Bytes()
						if len(b) > 0 {
							err := json.Unmarshal(b, x.Interface())
							if err != nil {
								session.Engine.logger.Error(err)
								return err
							}
							fieldValue.Set(x.Elem())
						}
					}
				} else if session.Statement.UseCascade {
					table := session.Engine.autoMapType(*fieldValue)
					if table != nil {
						hasAssigned = true
						if len(table.PrimaryKeys) != 1 {
							panic("unsupported non or composited primary key cascade")
						}
						var pk = make(core.PK, len(table.PrimaryKeys))

						switch rawValueType.Kind() {
						case reflect.Int64:
							pk[0] = vv.Int()
						case reflect.Int:
							pk[0] = int(vv.Int())
						case reflect.Int32:
							pk[0] = int32(vv.Int())
						case reflect.Int16:
							pk[0] = int16(vv.Int())
						case reflect.Int8:
							pk[0] = int8(vv.Int())
						case reflect.Uint64:
							pk[0] = vv.Uint()
						case reflect.Uint:
							pk[0] = uint(vv.Uint())
						case reflect.Uint32:
							pk[0] = uint32(vv.Uint())
						case reflect.Uint16:
							pk[0] = uint16(vv.Uint())
						case reflect.Uint8:
							pk[0] = uint8(vv.Uint())
						case reflect.String:
							pk[0] = vv.String()
						case reflect.Slice:
							pk[0], _ = strconv.ParseInt(string(rawValue.Interface().([]byte)), 10, 64)
						default:
							panic(fmt.Sprintf("unsupported primary key type: %v, %v", rawValueType, fieldValue))
						}

						if !isPKZero(pk) {
							// !nashtsai! TODO for hasOne relationship, it's preferred to use join query for eager fetch
							// however, also need to consider adding a 'lazy' attribute to xorm tag which allow hasOne
							// property to be fetched lazily
							structInter := reflect.New(fieldValue.Type())
							newsession := session.Engine.NewSession()
							defer newsession.Close()
							has, err := newsession.Id(pk).NoCascade().Get(structInter.Interface())
							if err != nil {
								return err
							}
							if has {
								//v := structInter.Elem().Interface()
								//fieldValue.Set(reflect.ValueOf(v))
								fieldValue.Set(structInter.Elem())
							} else {
								return errors.New("cascade obj is not exist")
							}
						}
					} else {
						session.Engine.logger.Error("unsupported struct type in Scan: ", fieldValue.Type().String())
					}
				}
			case reflect.Ptr:
				// !nashtsai! TODO merge duplicated codes above
				//typeStr := fieldType.String()
				switch fieldType {
				// following types case matching ptr's native type, therefore assign ptr directly
				case core.PtrStringType:
					if rawValueType.Kind() == reflect.String {
						x := vv.String()
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrBoolType:
					if rawValueType.Kind() == reflect.Bool {
						x := vv.Bool()
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrTimeType:
					if rawValueType == core.PtrTimeType {
						hasAssigned = true
						var x = rawValue.Interface().(time.Time)
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrFloat64Type:
					if rawValueType.Kind() == reflect.Float64 {
						x := vv.Float()
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrUint64Type:
					if rawValueType.Kind() == reflect.Int64 {
						var x = uint64(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrInt64Type:
					if rawValueType.Kind() == reflect.Int64 {
						x := vv.Int()
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrFloat32Type:
					if rawValueType.Kind() == reflect.Float64 {
						var x = float32(vv.Float())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrIntType:
					if rawValueType.Kind() == reflect.Int64 {
						var x = int(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrInt32Type:
					if rawValueType.Kind() == reflect.Int64 {
						var x = int32(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrInt8Type:
					if rawValueType.Kind() == reflect.Int64 {
						var x = int8(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrInt16Type:
					if rawValueType.Kind() == reflect.Int64 {
						var x = int16(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrUintType:
					if rawValueType.Kind() == reflect.Int64 {
						var x = uint(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrUint32Type:
					if rawValueType.Kind() == reflect.Int64 {
						var x = uint32(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.Uint8Type:
					if rawValueType.Kind() == reflect.Int64 {
						var x = uint8(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.Uint16Type:
					if rawValueType.Kind() == reflect.Int64 {
						var x = uint16(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.Complex64Type:
					var x complex64
					b := []byte(vv.String())
					if len(b) > 0 {
						err := json.Unmarshal(b, &x)
						if err != nil {
							session.Engine.logger.Error(err)
						} else {
							fieldValue.Set(reflect.ValueOf(&x))
						}
					}
					hasAssigned = true
				case core.Complex128Type:
					var x complex128
					b := []byte(vv.String())
					if len(b) > 0 {
						err := json.Unmarshal(b, &x)
						if err != nil {
							session.Engine.logger.Error(err)
						} else {
							fieldValue.Set(reflect.ValueOf(&x))
						}
					}
					hasAssigned = true
				} // switch fieldType
				// default:
				// 	session.Engine.LogError("unsupported type in Scan: ", reflect.TypeOf(v).String())
			} // switch fieldType.Kind()

			// !nashtsai! for value can't be assigned directly fallback to convert to []byte then back to value
			if !hasAssigned {
				data, err := value2Bytes(&rawValue)
				if err == nil {
					session.bytes2Value(col, fieldValue, data)
				} else {
					session.Engine.logger.Error(err.Error())
				}
			}
		}
	}
	return nil

}
