package gql

import (
	"fmt"
	"testing"
)

func TestMysqlBuilder(t *testing.T) {
	r := Read("files f").
		Join(Table(func(b Builder) {
			b.Table("hello")
		}, "t"), "t.user_id = f.id").
		GroupBy("f.id").
		Query()
	fmt.Println(r)
	c := Create("files").
		Fill(&map[string]interface{}{
			"hello": 123,
			"name": "aryan",
		}).
		Fill(&map[string]interface{}{
			"hello": 55,
			"name": "arash",
		}).
		Query()
	fmt.Println(c)
	u := Update("files").
		Fill(&map[string]interface{}{
			"hello": 123,
			"name": "aryan",
		}).
		Where("id", 2).
		Query()
	fmt.Println(u)
	d := Delete("files").
		Fill(&map[string]interface{}{
			"hello": 123,
			"name": "aryan",
		}).
		Where("id", 2).
		Query()
	fmt.Println(d)
}
