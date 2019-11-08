package gql

import (
	"database/sql"
	"encoding/json"
	"time"
)

type NullInt64 struct {
	sql.NullInt64
}

func NewInt64(value int64) NullInt64 {
	x := NullInt64{}
	x.Int64 = value
	x.Valid = true
	return x
}

func EmptyInt64() NullInt64 {
	x := NullInt64{}
	return x
}

func (m NullInt64) MarshalJSON() ([]byte, error) {
	if !m.Valid {
		return []byte("null"), nil
	}
	val, err := json.Marshal(m.Int64)
	if err != nil {
		return []byte("null"), nil
	}
	return val, nil
}
func (m *NullInt64) UnmarshalJSON(data []byte) error {
	var  val *int64
	err := json.Unmarshal(data, &val)
	if err != nil {
		return err
	}
	if val != nil {
		m.Valid = true
		m.Int64 = *val
	}else {
		m.Valid = false
	}
	return nil
}

type NullInt32 struct {
	sql.NullInt32
}

func NewInt32(value int32) NullInt32 {
	x := NullInt32{}
	x.Int32 = value
	x.Valid = true
	return x
}

func EmptyInt32() NullInt32 {
	x := NullInt32{}
	return x
}
func (m NullInt32) MarshalJSON() ([]byte, error) {
	if !m.Valid {
		return []byte("null"), nil
	}
	val, err := json.Marshal(m.Int32)
	if err != nil {
		return []byte("null"), nil
	}
	return val, nil
}
func (m *NullInt32) UnmarshalJSON(data []byte) error {
	var  val *int32
	err := json.Unmarshal(data, &val)
	if err != nil {
		return err
	}
	if val != nil {
		m.Valid = true
		m.Int32 = *val
	}else {
		m.Valid = false
	}
	return nil
}

type NullBool struct {
	sql.NullBool
}

func EmptyBool() NullBool {
	x := NullBool{}
	return x
}
func NewBool(value bool) NullBool {
	x := NullBool{}
	x.Bool = value
	x.Valid = true
	return x
}
func (m NullBool) MarshalJSON() ([]byte, error) {
	if !m.Valid {
		return []byte("null"), nil
	}
	val, err := json.Marshal(m.Bool)
	if err != nil {
		return []byte("null"), nil
	}
	return val, nil
}
func (m *NullBool) UnmarshalJSON(data []byte) error {
	var  val *bool
	err := json.Unmarshal(data, &val)
	if err != nil {
		return err
	}
	if val != nil {
		m.Valid = true
		m.Bool = *val
	}else {
		m.Valid = false
	}
	return nil
}
type NullFloat64 struct {
	sql.NullFloat64
}

func NewFloat64(value float64) NullFloat64 {
	x := NullFloat64{}
	x.Float64 = value
	x.Valid = true
	return x
}

func EmptyFloat64() NullFloat64 {
	x := NullFloat64{}
	return x
}
func (m NullFloat64) MarshalJSON() ([]byte, error) {
	if !m.Valid {
		return []byte("null"), nil
	}
	val, err := json.Marshal(m.Float64)
	if err != nil {
		return []byte("null"), nil
	}
	return val, nil
}
func (m *NullFloat64) UnmarshalJSON(data []byte) error {
	var  val *float64
	err := json.Unmarshal(data, &val)
	if err != nil {
		return err
	}
	if val != nil {
		m.Valid = true
		m.Float64 = *val
	}else {
		m.Valid = false
	}
	return nil
}
type NullString struct {
	sql.NullString
}

func NewString(value string) NullString {
	x := NullString{}
	x.String = value
	x.Valid = true
	return x
}

func EmptyString() NullString {
	x := NullString{}
	return x
}
func (m NullString) MarshalJSON() ([]byte, error) {
	if !m.Valid {
		return []byte("null"), nil
	}
	val, err := json.Marshal(m.String)
	if err != nil {
		return []byte("null"), nil
	}
	return val, nil
}
func (m *NullString) UnmarshalJSON(data []byte) error {
	var  val *string
	err := json.Unmarshal(data, &val)
	if err != nil {
		return err
	}
	if val != nil {
		m.Valid = true
		m.String = *val
	}else {
		m.Valid = false
	}
	return nil
}

type NullTime struct {
	sql.NullTime
}

func NewTime(value time.Time) NullTime {
	x := NullTime{}
	x.Time = value
	x.Valid = true
	return x
}

func EmptyTime() NullTime {
	x := NullTime{}
	return x
}
func (m NullTime) MarshalJSON() ([]byte, error) {
	if !m.Valid {
		return []byte("null"), nil
	}
	val, err := json.Marshal(m.Time)
	if err != nil {
		return []byte("null"), nil
	}
	return val, nil
}
func (m *NullTime) UnmarshalJSON(data []byte) error {
	var  val *time.Time
	err := json.Unmarshal(data, &val)
	if err != nil {
		return err
	}
	if val != nil {
		m.Valid = true
		m.Time = *val
	}else {
		m.Valid = false
	}
	return nil
}