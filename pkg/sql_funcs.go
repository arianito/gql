package gql

type SqlOp uint8

const (
	SqlAnd    = SqlOp(0)
	SqlOr     = SqlOp(1)
	SqlAndNot = SqlOp(2)
)

type SqlTyp uint8

const (
	SqlTypCreate = SqlTyp(0)
	SqlTypRead    = SqlTyp(1)
	SqlTypUpdate     = SqlTyp(2)
	SqlTypDelete = SqlTyp(3)
)

func Count(name string, alias ...string) string {
	op := "COUNT(" + name + ")"
	if len(alias) > 0 {
		op += " " + alias[0]
	}
	return op
}
func CountDistinct(name string, alias ...string) string {
	op := "COUNT(DISTINCT " + name + ")"
	if len(alias) > 0 {
		op += " " + alias[0]
	}
	return op
}

func Sum(expression string, alias ...string) string {
	op := "SUM(" + expression + ")"
	if len(alias) > 0 {
		op += " " + alias[0]
	}
	return op
}

func Query(fn func(builder Builder), alias string) string {
	b := &QueryBuilder{}
	fn(b)
	return "(" + b.Query() + ") " + alias
}

func Read(table string) Builder {
	q :=  QueryBuilder{}
	q.qtyp = SqlTypRead
	q.Table(table)
	return &q
}

func Create(table string) Builder {
	q :=  QueryBuilder{}
	q.qtyp = SqlTypCreate
	q.Table(table)
	return &q
}
func Update(table string) Builder {
	q :=  QueryBuilder{}
	q.qtyp = SqlTypUpdate
	q.Table(table)
	return &q
}
func Delete(table string) Builder {
	q :=  QueryBuilder{}
	q.qtyp = SqlTypDelete
	q.Table(table)
	return &q
}