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
				Or(func(builder Builder) {
					builder.Where("type", "3")
				})
		}).
		WhereGroup(func(builder Builder) {
			builder.
				Where("name", "hello").
				Or(func(builder Builder) {
					builder.WhereNot("type", "3")
				})
		}).
		Query()
	fmt.Println(q)
}
