package gql

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
	Top(top int64) Builder
	Offset(offset int64) Builder
	WhereGroup(fn func(b Builder)) Builder
	Where(clause string, value interface{}) Builder
	Find(value interface{}) Builder
	WhereNull(clause string) Builder
	WhereNotNull(clause string) Builder

	Model(ifc interface{}) Builder
	WhereLike(clause string, value interface{}) Builder


	WhereNotLike(clause string, value interface{}) Builder

	WhereNot(clause string, value interface{}) Builder
	WhereGT(clause string, value interface{}) Builder
	WhereGTE(clause string, value interface{}) Builder
	WhereLT(clause string, value interface{}) Builder
	WhereLTE(clause string, value interface{}) Builder
	WhereBetween(clause string, value1 interface{}, value2 interface{}) Builder
	WhereIn(field string, value []interface{}) Builder
	WhereInQuery(field string, fn func(b Builder)) Builder
	Or() Builder
	And() Builder
	AndNot() Builder
	Fill(values ...*OBJ) Builder
	Set(key string, value interface{}) Builder
	Bind(o interface{}) Builder
	BindExclude(o interface{}, keys ...string) Builder
	BindOnly(o interface{}, keys ...string) Builder
	Field(name string) SqlReserved
	Use(a interface{}) Builder
	Query() string
	Chunk(length int64, callback func(Scan func(o interface{}) Builder)) Builder
	Paginate(page int64, take int64) Builder
	Scan(o interface{}) Builder
	First(o interface{}) Builder
	Count(count *int64) Builder
	LastInsertionId(id *int64) Builder
	RowsAffected(count *int64) Builder
	GetScanLength(length *int64) Builder
	HasValue() bool
	Run() Builder
	GetError() error
}
