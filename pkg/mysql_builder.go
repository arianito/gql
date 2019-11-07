package gql

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"
)

type QueryBuilder struct {
	values  []*OBJ
	tables  []string
	columns []string
	wheres  []string
	orders  []string
	groups  []string
	joins   []string
	having  string
	ops     []SqlOp
	stp     SqlOp
	typ     SqlTyp

	limit  int64
	offset int64
	tx     *sql.Tx
	db     *sql.DB
	//
	obj            interface{}
	lastInsertedId int64
	rowsAffected   int64
	fln            int64
	cursor         *sql.Rows
	err            error
	//
}

func (b *QueryBuilder) Field(name string, attributes string) Builder {
	b.columns = append(b.columns, name+" "+attributes)
	return b
}
func (b *QueryBuilder) Unique(keys ...string) Builder {
	b.orders = append(b.orders, "UNIQUE("+strings.Join(keys, ", ")+")")
	return b
}
func (b *QueryBuilder) Index(keys ...string) Builder {
	b.orders = append(b.orders, "INDEX("+strings.Join(keys, ", ")+")")
	return b
}
func (b *QueryBuilder) PrimaryKey(key string) Builder {
	b.orders = append(b.orders, "PRIMARY KEY ("+key+")")

	return b
}
func (b *QueryBuilder) ForeignKey(localField string, remoteTable string, remoteField string, onDelete ...FKType) Builder {
	key := "FOREIGN KEY (" + localField + ") REFERENCES " + remoteTable + " (" + remoteField + ")"
	if len(onDelete) > 0 {
		if onDelete[0] == FKCascade {
			key += " ON DELETE SET NULL"
		} else if onDelete[0] == FKSetNull {
			key += " ON DELETE CASCADE"
		}
	}
	b.orders = append(b.orders, key)
	return b
}

func (b *QueryBuilder) Fill(values ...*OBJ) Builder {
	b.values = values
	return b
}

func (b *QueryBuilder) Columns(columns ...string) Builder {
	b.columns = append(b.columns, columns...)
	return b
}
func (b *QueryBuilder) Table(table string) Builder {
	b.tables = append(b.tables, table)
	return b
}

func (b *QueryBuilder) Join(table string, condition string, fn ...func(b Builder)) Builder {
	bld := &QueryBuilder{
		typ: SqlTypRead,
	}
	if len(fn) > 0 {
		fn[0](bld)
	}
	join := "JOIN " + table + " ON " + condition + " " + bld.getWhereClauses(true)
	b.joins = append(b.joins, join)
	return b
}

func (b *QueryBuilder) LeftJoin(table string, condition string, fn ...func(b Builder)) Builder {
	bld := &QueryBuilder{
		typ: SqlTypRead,
	}
	if len(fn) > 0 {
		fn[0](bld)
	}
	join := "LEFT JOIN " + table + " ON " + condition + " " + bld.getWhereClauses(true)
	b.joins = append(b.joins, join)
	return b
}

func (b *QueryBuilder) RightJoin(table string, condition string, fn ...func(b Builder)) Builder {
	bld := &QueryBuilder{
		typ: SqlTypRead,
	}
	if len(fn) > 0 {
		fn[0](bld)
	}
	join := "RIGHT JOIN " + table + " ON " + condition + " " + bld.getWhereClauses(true)
	b.joins = append(b.joins, join)
	return b
}

func (b *QueryBuilder) JoinUsing(table string, using string) Builder {
	join := "JOIN " + table + " USING(" + using + ")"
	b.joins = append(b.joins, join)
	return b
}

func (b *QueryBuilder) BitwiseAnd(field string, with int64, value int64) Builder {
	b.ops = append(b.ops, b.stp)
	b.wheres = append(b.wheres, fmt.Sprintf("%s & %v = %v", field, with, value))
	return b
}

func (b *QueryBuilder) BitwiseOr(field string, with int64, value int64) Builder {
	b.ops = append(b.ops, b.stp)
	b.wheres = append(b.wheres, fmt.Sprintf("%s | %v = %v", field, with, value))
	return b
}

func (b *QueryBuilder) OrderBy(clause ...string) Builder {
	b.orders = append(b.orders, clause...)
	return b
}

func (b *QueryBuilder) GroupBy(clause ...string) Builder {
	b.groups = append(b.groups, clause...)
	return b
}
func (b *QueryBuilder) Having(fn func(b Builder)) Builder {
	bld := &QueryBuilder{
		typ: SqlTypRead,
	}
	fn(bld)
	b.having = bld.getWhereClauses(false)
	return b
}

func (b *QueryBuilder) Where(field string, value interface{}) Builder {
	b.ops = append(b.ops, b.stp)
	b.wheres = append(b.wheres, fmt.Sprintf("%s = %s", field, interface_to_sql(value)))
	return b
}
func (b *QueryBuilder) Find(value interface{}) Builder {
	b.ops = append(b.ops, b.stp)
	b.wheres = append(b.wheres, fmt.Sprintf("id = %s", interface_to_sql(value)))
	return b
}
func (b *QueryBuilder) WhereNot(field string, value interface{}) Builder {
	b.ops = append(b.ops, b.stp)
	b.wheres = append(b.wheres, fmt.Sprintf("%s != %s", field, interface_to_sql(value)))
	return b
}

func (b *QueryBuilder) WhereNull(field string) Builder {
	b.ops = append(b.ops, b.stp)
	b.wheres = append(b.wheres, fmt.Sprintf("%s IS NULL", field))
	return b
}

func (b *QueryBuilder) WhereNotNull(field string) Builder {
	b.ops = append(b.ops, b.stp)
	b.wheres = append(b.wheres, fmt.Sprintf("%s IS NOT NULL", field))
	return b
}

func (b *QueryBuilder) WhereBetween(field string, value1 interface{}, value2 interface{}) Builder {
	b.ops = append(b.ops, b.stp)
	b.wheres = append(b.wheres, fmt.Sprintf("%s BETWEEN %s AND %s", field, interface_to_sql(value1), interface_to_sql(value2)))
	return b
}

func (b *QueryBuilder) WhereGT(field string, value interface{}) Builder {
	b.ops = append(b.ops, b.stp)
	b.wheres = append(b.wheres, fmt.Sprintf("%s > %s", field, interface_to_sql(value)))
	return b
}
func (b *QueryBuilder) WhereGTE(field string, value interface{}) Builder {
	b.ops = append(b.ops, b.stp)
	b.wheres = append(b.wheres, fmt.Sprintf("%s >= %s", field, interface_to_sql(value)))
	return b
}

func (b *QueryBuilder) WhereLT(field string, value interface{}) Builder {
	b.ops = append(b.ops, b.stp)
	b.wheres = append(b.wheres, fmt.Sprintf("%s < %s", field, interface_to_sql(value)))
	return b
}
func (b *QueryBuilder) WhereLTE(field string, value interface{}) Builder {
	b.ops = append(b.ops, b.stp)
	b.wheres = append(b.wheres, fmt.Sprintf("%s <= %s", field, interface_to_sql(value)))
	return b
}

func (b *QueryBuilder) WhereIn(field string, value []interface{}) Builder {
	b.ops = append(b.ops, b.stp)
	b.wheres = append(b.wheres, fmt.Sprintf("%s in %s", field, interface_to_sql(value)))
	return b
}

func (b *QueryBuilder) WhereInQuery(field string, fn func(b Builder)) Builder {
	b.ops = append(b.ops, b.stp)
	bld := &QueryBuilder{
		typ: SqlTypRead,
	}
	fn(bld)
	b.wheres = append(b.wheres, fmt.Sprintf("%s IN (%s)", field, bld.Query()))
	return b
}

func (b *QueryBuilder) WhereGroup(fn func(b Builder)) Builder {
	builder := &QueryBuilder{
		typ: SqlTypRead,
	}
	fn(builder)
	b.wheres = append(b.wheres, fmt.Sprintf("(%s)", builder.getWhereClauses(false)))
	b.ops = append(b.ops, b.stp)
	return b
}

func (b *QueryBuilder) Or() Builder {
	b.stp = SqlOr
	return b
}
func (b *QueryBuilder) And() Builder {
	b.stp = SqlAnd
	return b
}
func (b *QueryBuilder) AndNot() Builder {
	b.stp = SqlAndNot
	return b
}

func (b *QueryBuilder) getWhereClauses(flag bool) string {
	where := ""
	ln := len(b.wheres)
	for i := 0; i < ln; i++ {
		cls := b.wheres[i]
		if flag || i != 0 {
			op := b.ops[i]
			switch op {
			case SqlAnd:
				where += "AND "
				break
			case SqlOr:
				where += "OR "
				break
			case SqlAndNot:
				where += "AND NOT "
				break
			}
		}
		where += cls
		if i != ln-1 {
			where += " "
		}
	}
	return where
}

func (b *QueryBuilder) Query() (out string) {
	if enableLog {
		defer func() {
			log.Println(out)
		}()
	}
	if b.typ == SqlTypRead {
		columns := "*"
		if len(b.columns) > 0 {
			columns = strings.Join(b.columns, ", ")
		}
		orderBy := ""
		if len(b.orders) > 0 {
			orderBy = " ORDER BY " + strings.Join(b.orders, ", ")
		}
		where := ""
		if len(b.wheres) > 0 {
			where = " WHERE 1 " + b.getWhereClauses(true)
		}
		joins := ""
		if len(b.joins) > 0 {
			joins = " " + strings.Join(b.joins, " ")
		}
		groupBy := ""
		if len(b.groups) > 0 {
			groupBy = " GROUP BY " + strings.Join(b.groups, ", ")
			if b.having != "" {
				groupBy += " HAVING " + b.having
			}
		}
		limit := ""
		if b.limit > 0 {
			limit = fmt.Sprintf(" LIMIT %v", b.limit)
		}
		offset := ""
		if b.offset > 0 {
			offset = fmt.Sprintf(" OFFSET %v", b.offset)
		}
		query := "SELECT " + columns + " FROM " + strings.Join(b.tables, ", ") + joins + where + groupBy + orderBy + limit + offset
		out = strings.Trim(query, " ")
	} else if b.typ == SqlTypCreate {
		keys := make([]string, 0)
		for key := range *b.values[0] {
			keys = append(keys, key)
		}
		ln := len(keys)

		stm := make([]string, 0)

		for _, item := range b.values {
			values := ""
			i := 0
			for _, key := range keys {
				values += interface_to_sql((*item)[key])
				if i != ln-1 {
					values += ", "
				}
				i++
			}
			stm = append(stm, "("+values+")")
		}
		out = "INSERT INTO " + b.tables[0] + "(" + strings.Join(keys, ", ") + ") VALUES" + strings.Join(stm, ", ")
	} else if b.typ == SqlTypUpdate {
		item := b.values[0]
		values := ""
		i := 0
		ln := len(*item)
		for key, value := range *item {
			values += key + "=" + interface_to_sql(value)
			if i != ln-1 {
				values += ", "
			}
			i++
		}

		where := ""
		if len(b.wheres) > 0 {
			where = " WHERE 1 " + b.getWhereClauses(true)
		}

		set := ""
		if len(b.values) > 0 {
			set = " SET " + values
		}

		out = "UPDATE " + strings.Join(b.tables, ", ") + set + where
	} else if b.typ == SqlTypDelete {
		out = "DELETE FROM " + strings.Join(b.tables, ", ") + " WHERE " + b.getWhereClauses(false)
	} else if b.typ == SqlTypTable {
		table := "CREATE TABLE " + strings.Join(b.tables, ", ") + "(" + strings.Join(b.columns, ", ")
		if len(b.columns) > 0 && len(b.orders) > 0 {
			table += ", " + strings.Join(b.orders, ", ")
		}
		out = table + ")"
	}
	return
}

func (b *QueryBuilder) Use(a interface{}) Builder {
	tx, ok := a.(*sql.Tx)
	if ok {
		b.tx = tx
		b.db = nil
	} else {
		db, ok := a.(*sql.DB)
		if ok {
			b.tx = nil
			b.db = db
		} else {
			panic("wrong db provider used")
		}
	}
	return b
}

func (b *QueryBuilder) Top(top int64) Builder {
	b.limit = top
	return b
}
func (b *QueryBuilder) Offset(offset int64) Builder {
	b.offset = offset
	return b
}

func (b *QueryBuilder) First(o interface{}) Builder {
	b.limit = 1
	return b.Scan(o)
}

func (b *QueryBuilder) query() (*sql.Rows, error) {
	if b.db == nil && b.tx == nil {
		panic("db driver not defined")
	}
	if b.db != nil {
		return b.db.Query(b.Query())
	} else {
		return b.tx.Query(b.Query())
	}
}

func (b *QueryBuilder) Count(count *int64) Builder {
	var cpy []string
	cpy = b.columns
	b.columns = []string{Count("*", "len")}
	type LenObj struct {
		Len int64 `json:"len"`
	}
	var obj LenObj
	b.Scan(&obj)
	*count = obj.Len
	b.columns = cpy
	return b
}
func (b *QueryBuilder) LastInsertionId(id *int64) Builder {
	*id = b.lastInsertedId
	return b
}
func (b *QueryBuilder) RowsAffected(count *int64) Builder {
	*count = b.rowsAffected
	return b
}
func (b *QueryBuilder) GetScanLength(length *int64) Builder {
	*length = b.fln
	return b
}

func (b *QueryBuilder) Paginate(page int64, take int64) (out Builder) {
	out = b
	b.Top(take)
	b.Offset(page * take)
	return
}
func (b *QueryBuilder) Set(key string, value interface{}) (out Builder) {
	out = b
	for _, value := range b.values {
		(*value)[key] = value
	}
	return
}
func some(stack []string, check func(key string) bool) bool {
	for _, key := range stack {
		if check(key) {
			return true
		}
	}
	return false
}

func (b *QueryBuilder) Bind(o interface{}) (out Builder) {
	out = b
	b.bind(0, o)
	return
}

func (b *QueryBuilder) BindOnly(o interface{}, keys ...string) (out Builder) {
	out = b
	b.bind(1, o, keys...)
	return
}

func (b *QueryBuilder) BindExclude(o interface{}, keys ...string) (out Builder) {
	out = b
	b.bind(2, o, keys...)
	return
}


func checkEqual(field reflect.StructField, name string) bool {
	tag := field.Tag.Get("gql")
	if tag != name {
		jsn := field.Tag.Get("json")
		if jsn != name {
			if strings.ToLower(name) != strings.ToLower(field.Name) {
				return false
			}
		}
	}
	return true
}

func (b *QueryBuilder) bind(mode int, o interface{}, keys ...string) (out Builder) {
	out = b
	b.obj = o

	var err error
	defer func() {
		if p := recover(); p != nil {
			b.err = fmt.Errorf("%v", p)
			return
		} else if err != nil {
			b.err = err
		}
	}()

	vf := reflect.ValueOf(o).Elem()

	if vf.Kind() == reflect.Slice {
		b.values = make([]*OBJ, 0)
		ln := vf.Len()
		for i := 0; i < ln; i++ {
			val := vf.Index(i)
			if val.Kind() != reflect.Struct {
				val = val.Elem()
			}
			fln := val.NumField()
			data := make(OBJ)
			for j := 0; j < fln; j++ {
				value := val.Field(j)
				field := val.Type().Field(j)
				tag := field.Tag.Get("gql")


				allow := true

				if mode == 1 { // only
					allow = some(keys, func(key string) bool {
						return checkEqual(field, key)
					})
				}else if mode == 2 { // exclude
					allow = !some(keys, func(key string) bool {
						return checkEqual(field, key)
					})
				}

				if tag != "-" && tag != "" && allow {
					if value.CanInterface() {
						data[tag] = value.Interface()

					}
				}
			}
			b.values = append(b.values, &data)
		}
		return
	}
	if vf.Kind() != reflect.Struct {
		vf = vf.Elem()
	}
	data := make(OBJ)
	ln := vf.NumField()
	for i := 0; i < ln; i++ {
		field := vf.Type().Field(i)
		value := vf.Field(i)
		tag := field.Tag.Get("gql")

		allow := true

		if mode == 1 { // only
			allow = some(keys, func(key string) bool {
				return checkEqual(field, key)
			})
		}else if mode == 2 { // exclude
			allow = !some(keys, func(key string) bool {
				return checkEqual(field, key)
			})
		}

		if tag != "-" && tag != "" && allow {
			if value.CanInterface() {
				data[tag] = value.Interface()
			}
		}
	}
	b.values = []*OBJ{&data}
	return
}
func (b *QueryBuilder) Scan(o interface{}) (out Builder) {
	out = b
	var err error
	defer func() {
		if p := recover(); p != nil {
			b.err = fmt.Errorf("%v", p)
			return
		} else if err != nil {
			b.err = err
		}
	}()

	tf := reflect.TypeOf(o).Elem()
	vf := reflect.ValueOf(o).Elem()

	if tf.Kind() == reflect.Slice {
		stc := true
		elem := tf.Elem()
		if elem.Kind() != reflect.Struct {
			elem = elem.Elem()
			stc = false
		}
		pairs := make(map[string]string)
		ln := elem.NumField()
		for i := 0; i < ln; i++ {
			field := elem.Field(i)
			tag := field.Tag.Get("gql")
			if tag != "-" {
				if tag == "" {
					tag = field.Tag.Get("json")
				}

				if tag != "" {
					pairs[tag] = field.Name
				}

			}
		}
		var rows *sql.Rows
		rows, err = b.query()
		if err != nil {
			return
		}

		vf.Set(reflect.MakeSlice(tf, 0, 0))
		b.fln = 0
		for rows.Next() {
			var data []string
			data, err = rows.Columns()
			if err != nil {
				rows.Close()
				return
			}

			val := reflect.New(elem)
			el := val.Elem()
			ifc := make([]interface{}, len(data))
			for i, str := range data {
				op := pairs[str]
				if op != "" {
					obj := el.FieldByName(op).Addr().Interface()
					ifc[i] =
						obj
				} else {
					var obj interface{}
					ifc[i] = &obj
				}
			}
			err = rows.Scan(ifc...)
			if err != nil {
				return
			}
			if stc {
				vf.Set(reflect.Append(vf, val.Elem()))
			} else {
				vf.Set(reflect.Append(vf, val))
			}
			b.fln++
		}

	} else {
		elem := tf
		val := vf

		if elem.Kind() != reflect.Struct {
			elem = elem.Elem()
			val = val.Elem()
		}

		pairs := make(map[string]string)
		ln := elem.NumField()
		for i := 0; i < ln; i++ {
			field := elem.Field(i)
			tag := field.Tag.Get("gql")
			if tag != "-" {
				if tag == "" {
					tag = field.Tag.Get("json")
				}

				if tag != "" {
					pairs[tag] = field.Name
				}

			}
		}
		var rows *sql.Rows
		rows, err = b.query()
		if err != nil {
			return
		}
		b.fln = 0
		for rows.Next() {
			var data []string
			data, err = rows.Columns()
			if err != nil {
				rows.Close()
				return
			}
			ifc := make([]interface{}, len(data))
			for i, str := range data {
				op := pairs[str]
				if op != "" {
					obj := val.FieldByName(op).Addr().Interface()
					ifc[i] =
						obj
				} else {
					var obj interface{}
					ifc[i] = &obj
				}
			}
			err = rows.Scan(ifc...)
			if err != nil {
				return
			}
			b.fln++
			rows.Close()
			return
		}
	}
	return
}

func (b *QueryBuilder) GetError() error {
	return b.err
}

func (b *QueryBuilder) Chunk(length int64, callback func(Scan func(o interface{}) Builder)) (out Builder) {
	out = b
	var offset int64 = 0
	for {
		b.Top(length)
		b.Offset(offset)
		callback(b.Scan)
		var ln int64
		b.GetScanLength(&ln)
		if ln < 1 || b.GetError() != nil {
			return
		}
		offset += length
	}
}
func (b *QueryBuilder) scan(o interface{}) {
	b.Scan(o)
}
func (b *QueryBuilder) HasValue() bool {
	return b.GetError() == nil && b.fln > 0
}
func (b *QueryBuilder) Run() (out Builder) {
	out = b
	var err error
	defer func() {
		if p := recover(); p != nil {
			b.err = fmt.Errorf("%v", p)
			return
		} else if err != nil {
			b.err = err
		}
	}()

	if b.db == nil && b.tx == nil {
		panic("db driver not defined")
	}

	var a sql.Result

	if b.db != nil {
		a, err = b.db.Exec(b.Query())
	} else {
		a, err = b.tx.Exec(b.Query())
	}
	if err != nil {
		return
	}
	b.lastInsertedId, err = a.LastInsertId()

	if b.obj != nil {
		vf := reflect.ValueOf(b.obj).Elem()
		if vf.Kind() != reflect.Struct && vf.Kind() != reflect.Slice {
			vf = vf.Elem()
		}
		if vf.Kind() == reflect.Struct {
			ln := vf.NumField()
			for i := 0; i < ln; i++ {
				val := vf.Field(i)
				tags := vf.Type().Field(i).Tag
				if tags.Get("gql") == "id" || tags.Get("pk") == "true" {
					val.Set(reflect.ValueOf(b.lastInsertedId))
					break
				}
			}
		}
	}
	if err != nil {
		return
	}
	b.rowsAffected, err = a.RowsAffected()
	return
}
