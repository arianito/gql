package gql

import (
	"fmt"
	"testing"
	"time"
)



func TestMysqlBuilder(t *testing.T) {
	type File struct {
		Name     string     `gql:"fnam"`
		JobTitle string     `gql:"jtitl"`
		Hint     string     `gql:"hint"`
		Hello     int64
		Time     *time.Time `gql:"tm"`
	}
	file := []*File{
		{
			Name: "aryan",
			JobTitle: "developer",
			Hint: "hello",
		},
		{
			Name:     "jacob",
			JobTitle: "manager",
			Hint: "bye",
		},
	}
	a := Create("files").BindExclude(&file, "hello").Run()
	fmt.Println(a.GetError())
	fmt.Println(a.Query())
}
