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
func interface_to_sql(value interface{}) string {
	switch value.(type) {
	case string:
		return fmt.Sprintf("'%v'", strings.ReplaceAll(value.(string), "'", "\\'"))
	case float32:
		return fmt.Sprintf("%s", float_to_string(float64(value.(float32))))
	case float64:
		return fmt.Sprintf("%s", float_to_string(value.(float64)))
	case time.Time:
		d := value.(time.Time)
		return interface_to_sql(d.Format("2006-01-02 15:04:05"))
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, bool:
		return fmt.Sprintf("%v", value)
	default:
		v := reflect.ValueOf(value)
		if v.Kind() != reflect.Slice {
			return ""
		}
		ln := v.Len()
		op:= "("
		for i:=0;i<ln ;i++  {
			val := interface_to_sql(v.Index(i).Interface())
			if val != "" {
				op += val
				if i != ln - 1 {
					op += ","
				}
			}
		}
		op+=")"
		return op
	}
}
