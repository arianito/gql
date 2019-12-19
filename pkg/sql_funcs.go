package gql

type SqlOp uint8

const (
	SqlAnd    = SqlOp(0)
	SqlOr     = SqlOp(1)
	SqlAndNot = SqlOp(2)
)

type SqlTyp uint8

type SqlReserved struct {
	content string
}

const (
	SqlTypCreate = SqlTyp(0)
	SqlTypRead   = SqlTyp(1)
	SqlTypUpdate = SqlTyp(2)
	SqlTypDelete = SqlTyp(3)
	SqlTypCustom = SqlTyp(4)
)

func Now() SqlReserved {
	return SqlReserved{content: "NOW()"}
}
func Sql(sql string) SqlReserved {
	return SqlReserved{content: sql}
}

func Count(name string, alias ...string) SqlReserved {
	op := "COUNT(" + name + ")"
	if len(alias) > 0 {
		op += " " + alias[0]
	}
	return SqlReserved{content:op}
}
func CountDistinct(name string, alias ...string) SqlReserved {
	op := "COUNT(DISTINCT " + name + ")"
	if len(alias) > 0 {
		op += " " + alias[0]
	}
	return SqlReserved{content:op}

}

func Sum(expression string, alias ...string) SqlReserved {
	op := "SUM(" + expression + ")"
	if len(alias) > 0 {
		op += " " + alias[0]
	}
	return SqlReserved{content:op}
}

func Query(fn func(builder Builder), alias string) string {
	b := &QueryBuilder{
		typ: SqlTypRead,
	}
	fn(b)
	return "(" + b.Query() + ") " + alias
}

func Custom(query string) Builder {
	q := QueryBuilder{}
	q.typ = SqlTypCustom
	q.customQuery = query
	return &q
}

func Read(table string) Builder {
	q := QueryBuilder{}
	q.typ = SqlTypRead
	q.Table(table)
	return &q
}

func Create(table string) Builder {
	q := QueryBuilder{}
	q.typ = SqlTypCreate
	q.Table(table)
	return &q
}
func Update(table string) Builder {
	q := QueryBuilder{}
	q.typ = SqlTypUpdate
	q.Table(table)
	return &q
}
func Delete(table string) Builder {
	q := QueryBuilder{}
	q.typ = SqlTypDelete
	q.Table(table)
	return &q
}

type OBJ map[string]interface{}
