package gql

import "database/sql"

type FKType byte

const (
	FKNormal  = FKType(0)
	FKCascade = FKType(1)
	FKSetNull    = FKType(2)
)

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
	Find(value interface{}) Builder
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
	Fill(values ...*OBJ) Builder
	Or() Builder
	And() Builder
	Count() Builder
	AndNot() Builder

	Field(name string, attributes string) Builder
	Unique(keys ...string) Builder
	Index(keys ...string) Builder
	PrimaryKey(key string) Builder
	ForeignKey(localField string, remoteTable string, remoteField string, typ ...FKType) Builder

	UseTx(tx *sql.Tx) Builder
	UseDb(db *sql.DB) Builder
	Query() string
	QueryRows() (*sql.Rows, error)
	QueryRow(args ...interface{}) error
	Exec() (int64, int64, error)
}
