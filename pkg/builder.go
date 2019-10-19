package gql

import (
	"fmt"
	"strings"
)

type SqlOp uint8

const (
	SqlAnd = SqlOp(0)
	SqlOr = SqlOp(1)
	SqlAndNot = SqlOp(2)
)

type Builder interface {
	Select(columns ...string) Builder
	Table(table string) Builder
	Or() Builder
	And() Builder
	AndNot() Builder
	WhereGroup(fn func(builder Builder)) Builder
	Where(clause string, value interface{}) Builder
	WhereNot(clause string, value interface{}) Builder
	WhereGT(clause string, value interface{}) Builder
	WhereGTE(clause string, value interface{}) Builder
	WhereLT(clause string, value interface{}) Builder
	WhereLTE(clause string, value interface{}) Builder
	WhereBetween(clause string, value1 interface{}, value2 interface{}) Builder
	WhereIn(field string, value []interface{}) Builder
	WhereInQuery(field string, fn func(b Builder)) Builder
	Query() string
}

type MYSQLBuilder struct {
	tables       []string
	columns      []string
	whereClauses []string
	ops []SqlOp
	stp SqlOp
}

func NewMYSQLBuilder() Builder  {
	return &MYSQLBuilder{}
}

func (b *MYSQLBuilder) Select(columns ...string) Builder {
	b.columns = append(b.columns, columns...)
	return b
}
func (b *MYSQLBuilder) Table(table string) Builder {
	b.tables = append(b.tables, table)
	return b
}


func (b *MYSQLBuilder) BitwiseAnd(field string, with int64, value int64) Builder {
	b.ops = append(b.ops, b.stp)
	b.whereClauses = append(b.whereClauses, fmt.Sprintf("%s & %v = %v", field, with, value))
	return b
}

func (b *MYSQLBuilder) BitwiseOr(field string, with int64, value int64) Builder {
	b.ops = append(b.ops, b.stp)
	b.whereClauses = append(b.whereClauses, fmt.Sprintf("%s | %v = %v", field, with, value))
	return b
}

func (b *MYSQLBuilder) Where(field string, value interface{}) Builder {
	b.ops = append(b.ops, b.stp)
	b.whereClauses = append(b.whereClauses, fmt.Sprintf("%s = %s", field, interface_to_sql(value)))
	return b
}
func (b *MYSQLBuilder) WhereNot(field string, value interface{}) Builder {
	b.ops = append(b.ops, b.stp)
	b.whereClauses = append(b.whereClauses, fmt.Sprintf("%s != %s", field, interface_to_sql(value)))
	return b
}

func (b *MYSQLBuilder) WhereBetween(field string, value1 interface{}, value2 interface{}) Builder {
	b.ops = append(b.ops, b.stp)
	b.whereClauses = append(b.whereClauses, fmt.Sprintf("%s between %s and %s", field, interface_to_sql(value1), interface_to_sql(value2)))
	return b
}

func (b *MYSQLBuilder) WhereGT(field string, value interface{}) Builder {
	b.ops = append(b.ops, b.stp)
	b.whereClauses = append(b.whereClauses, fmt.Sprintf("%s > %s", field, interface_to_sql(value)))
	return b
}
func (b *MYSQLBuilder) WhereGTE(field string, value interface{}) Builder {
	b.ops = append(b.ops, b.stp)
	b.whereClauses = append(b.whereClauses, fmt.Sprintf("%s >= %s", field, interface_to_sql(value)))
	return b
}

func (b *MYSQLBuilder) WhereLT(field string, value interface{}) Builder {
	b.ops = append(b.ops, b.stp)
	b.whereClauses = append(b.whereClauses, fmt.Sprintf("%s < %s", field, interface_to_sql(value)))
	return b
}
func (b *MYSQLBuilder) WhereLTE(field string, value interface{}) Builder {
	b.ops = append(b.ops, b.stp)
	b.whereClauses = append(b.whereClauses, fmt.Sprintf("%s <= %s", field, interface_to_sql(value)))
	return b
}

func (b *MYSQLBuilder) WhereIn(field string, value []interface{}) Builder {
	b.ops = append(b.ops, b.stp)
	b.whereClauses = append(b.whereClauses, fmt.Sprintf("%s in %s", field, interface_to_sql(value)))
	return b
}


func (b *MYSQLBuilder) WhereInQuery(field string, fn func(b Builder)) Builder {
	b.ops = append(b.ops, b.stp)
	bld := &MYSQLBuilder{}
	fn(bld)
	b.whereClauses = append(b.whereClauses, fmt.Sprintf("%s in (%s)", field, bld.Query()))
	return b
}




func (b *MYSQLBuilder) WhereGroup(fn func(b Builder)) Builder {
	builder := &MYSQLBuilder{}
	fn(builder)
	b.whereClauses = append(b.whereClauses, fmt.Sprintf("(%s)", builder.computeWhereClauses(false)))
	b.ops = append(b.ops, b.stp)
	return b
}



func (b *MYSQLBuilder) Or() Builder {
	b.stp = SqlOr
	return b
}
func (b *MYSQLBuilder) And() Builder {
	b.stp = SqlAnd
	return b
}
func (b *MYSQLBuilder) AndNot() Builder {
	b.stp = SqlAndNot
	return b
}


func (b *MYSQLBuilder) computeWhereClauses(flag bool) string {
	where := ""
	ln := len(b.whereClauses)
	for i := 0; i < ln; i++ {
		cls := b.whereClauses[i]
		if flag || i != 0 {
			op := b.ops[i]
			switch op {
			case SqlAnd:
				where += "and "
				break
			case SqlOr:
				where += "or "
				break
			case SqlAndNot:
				where += "and not "
				break
			}
		}
		where += cls
		if i != ln - 1 {
			where += " "
		}
	}
	return where
}

func (b *MYSQLBuilder) Query() string {
	columns := "*"
	if len(b.columns) > 0 {
		columns = strings.Join(b.columns, ",")
	}
	query := "select " + columns + " from " + strings.Join(b.tables, ",") + " "
	where := "where 1 "
	where += b.computeWhereClauses(true)
	query += where
	return query + ";"
}
