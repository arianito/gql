package gql

import (
	"fmt"
	"testing"
)

func TestMysqlBuilder(t *testing.T) {
	r := Read("files f").
		Join(Query(func(b Builder) {
			b.Table("hello")
		}, "t"), "t.user_id = f.id").
		GroupBy("f.id").
		Query()
	fmt.Println(r)
	c := Create("files").
		Fill(
			&OBJ{
				"hello": 1,
				"name":  "aryan",
			},
			&OBJ{
				"name":  "arash",
				"hello": 55,
			},
		).
		Query()
	fmt.Println(c)
	u := Update("files").
		Fill(&OBJ{
			"hello": 123,
			"name":  "aryan",
		}).
		Where("id", 2).
		Query()
	fmt.Println(u)
	d := Delete("files").
		Fill(&OBJ{
			"hello": 123,
			"name":  "aryan",
		}).
		Where("id", 2).
		Query()
	fmt.Println(d)
}
