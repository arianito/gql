package gql

import (
	"bytes"
	"database/sql"
	"encoding/hex"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func float_to_string(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
func Convert(value interface{}) (out string) {
	out =  "NULL"
	switch value.(type) {
	case string:
		d := value.(string)
		out =   `'`+strings.ReplaceAll(d, `'`, `\'`)+`'`
		return
	case []byte:
		d := value.([]byte)
		out = "X'"+hex.EncodeToString(d)+"'"
	case sql.RawBytes:
		d := value.(sql.RawBytes)
		out = "X'"+hex.EncodeToString(d)+"'"
	case bytes.Buffer:
		d := value.(bytes.Buffer)
		out = "X'"+hex.EncodeToString(d.Bytes())+"'"
		return
	case *bytes.Buffer:
		d := value.(*bytes.Buffer)
		out = "X'"+hex.EncodeToString(d.Bytes())+"'"
		return
	case NullString:
		d := value.(NullString)
		if d.Valid {
			out =  `'`+strings.ReplaceAll(d.String, `'`, `\'`)+`'`
		}else {
			out = "NULL"
		}
		return
	case float32:
		out =  float_to_string(float64(value.(float32)))
		return
	case float64:
		out =  float_to_string(value.(float64))
		return

	case NullFloat64:
		d := value.(NullFloat64)
		if d.Valid {
			out =  float_to_string(d.Float64)
		}else {
			out = "NULL"
		}
	case SqlReserved:
		out =  (value.(SqlReserved)).content
		return
	case time.Time:
		d := value.(time.Time)
		out =  Convert(d.UTC().Format("2006-01-02 15:04:05"))
		return
	case NullTime:
		d := value.(NullTime)
		if d.Valid {
			out = Convert(d.Time.UTC().Format("2006-01-02 15:04:05"))
		}else {
			out = "NULL"
		}
		return

	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, bool:
		out =  fmt.Sprintf("%v", value)
		return
	case NullInt64:
		d := value.(NullInt64)
		if d.Valid {
			out = strconv.FormatInt(d.Int64, 10)
		}else {
			out = "NULL"
		}
	case NullInt32:
		d := value.(NullInt32)
		if d.Valid {
			out =  strconv.FormatInt(int64(d.Int32), 10)
		}else {
			out = "NULL"
		}
	case NullBool:
		d := value.(NullBool)
		if d.Valid {
			if d.Bool {
				out = "true"
			} else  {
				out = "false"
			}
		}else {
			out = "NULL"
		}
	default:
		v := reflect.ValueOf(value)
		if v.Kind() != reflect.Slice {
			return
		}
		ln := v.Len()
		out = "("
		for i:=0;i<ln ;i++  {
			val := Convert(v.Index(i).Interface())
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
	return
}
