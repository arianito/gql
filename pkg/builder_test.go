package gql

import (
	"fmt"
	"testing"
	"time"
)

func TestMysqlBuilder(t *testing.T) {
	type File struct {
		ID        int64  `gql:"id" pk:"true"`
		Name      string `gql:"fnam"`
		JobTitle  string `gql:"jtitl"`
		Hint      string `gql:"hint"`
		Hello     int64
		Time      *time.Time  `gql:"tm"`
		DummyTime interface{} `gql:"tm"`
	}
	file := File{
		Name:      "aryan",
		JobTitle:  "developer",
		Hint:      "hello",
		DummyTime: Now(),
	}

	a := Create("files").BindOnly(&file, "name", "jobtitle", "tm").Run()
	fmt.Println(a.GetError())
	fmt.Println(a.Query())
}
