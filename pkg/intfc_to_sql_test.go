package gql

import (
	"fmt"
	"testing"
)

func TestIntfcToSql(t *testing.T) {
	a := interface_to_sql([]interface{}{1,2,3,"hello aryan's aunt"})
	fmt.Println(a)
}
