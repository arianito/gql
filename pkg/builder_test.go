package gql

import (
	"fmt"
	"testing"
)

func TestMysqlBuilder(t *testing.T) {
	q := NewMYSQLBuilder().
		Table("users").
		WhereGroup(func(builder Builder) {
			builder.
				Where("name", "hello").
				Or().Where("type", "3")
		}).
		WhereGroup(func(builder Builder) {
			builder.
				Where("name", "hello").
				Or().WhereNot("type", "3").And().Where("name", "test")
		}).
		Query()
	fmt.Println(q)
}
