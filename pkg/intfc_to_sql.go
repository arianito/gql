package gql

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func float_to_string(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
func interface_to_sql(value interface{}) (out string) {
	defer func() {
		if p:= recover(); p != nil {
			out = "NULL"
		}
	}()
	if value == nil {
		out =  "NULL"
		return
	}
	switch value.(type) {
	case string:
		out =  fmt.Sprintf("'%v'", strings.ReplaceAll(value.(string), "'", "\\'"))
		return
	case float32:
		out =  fmt.Sprintf("%s", float_to_string(float64(value.(float32))))
		return
	case float64:
		out =  fmt.Sprintf("%s", float_to_string(value.(float64)))
		return
	case SqlReserved:
		out =  (value.(SqlReserved)).content
		return
	case *SqlReserved:
		out =  (value.(*SqlReserved)).content
		return
	case time.Time:
		d := value.(time.Time)
		out =  interface_to_sql(d.Format("2006-01-02 15:04:05"))
		return
	case *time.Time:
		d := value.(*time.Time)
		out = interface_to_sql(d.Format("2006-01-02 15:04:05"))
		return
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, bool:
		out =  fmt.Sprintf("%v", value)
		return
	default:
		v := reflect.ValueOf(value)
		if v.Kind() != reflect.Slice {
			out = ""
			return
		}
		ln := v.Len()
		out = "("
		for i:=0;i<ln ;i++  {
			val := interface_to_sql(v.Index(i).Interface())
			if val != "" {
				out += val
				if i != ln - 1 {
					out += ","
				}
			}
		}
		out+=")"
		return
	}
}
