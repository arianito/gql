package gql

import "database/sql"

type Builder interface {
	Table(table string) Builder
	Columns(columns ...string) Builder
	BitwiseOr(field string, with int64, value int64) Builder
	BitwiseAnd(field string, with int64, value int64) Builder
	Join(table string, on string, fn ...func(b Builder)) Builder
	LeftJoin(table string, on string, fn ...func(b Builder)) Builder
	RightJoin(table string, on string, fn ...func(b Builder)) Builder
	JoinUsing(table string, using string) Builder
	OrderBy(clause ...string) Builder
	GroupBy(clause ...string) Builder
	Having(fn func(b Builder)) Builder
	WhereGroup(fn func(b Builder)) Builder
	Where(clause string, value interface{}) Builder
	WhereNull(clause string) Builder
	WhereNotNull(clause string) Builder
	WhereNot(clause string, value interface{}) Builder
	WhereGT(clause string, value interface{}) Builder
	WhereGTE(clause string, value interface{}) Builder
	WhereLT(clause string, value interface{}) Builder
	WhereLTE(clause string, value interface{}) Builder
	WhereBetween(clause string, value1 interface{}, value2 interface{}) Builder
	WhereIn(field string, value []interface{}) Builder
	WhereInQuery(field string, fn func(b Builder)) Builder
	Fill(values ...*map[string]interface{}) Builder
	Or() Builder
	And() Builder
	Count() Builder
	AndNot() Builder

	Query() string
	QueryTx(tx *sql.Tx) (*sql.Rows, error)
	QueryDb(tx *sql.DB) (*sql.Rows, error)
	QueryRowTx(tx *sql.Tx, args...interface{}) error
	QueryRowDb(tx *sql.DB, args...interface{}) error
	ExecTx(tx *sql.Tx) (int64, int64, error)
	ExecDb(tx *sql.DB) (int64, int64, error)
}