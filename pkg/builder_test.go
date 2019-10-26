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
				"hello": 55,
				"name":  "arash",
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
	h := Table("files").
		Field("id", "int not null primary key auto_increment").
		Field("name", "varchar(100)").
		ForeignKey("name", "users", "id").
		Query()
	fmt.Println(h)
}
