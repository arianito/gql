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
	u     interface{}
	//
	obj            interface{}
	fldTag         map[string]string
	lastInsertedId int64
	rowsAffected   int64
	fln            int64
	cursor         *sql.Rows
	err            error
	structFields   map[string]int
	//
}

func (b *QueryBuilder) extractName(name string) string {
	if name[0] == '(' || b.obj == nil || b.fldTag == nil {
		return name
	}
	spl := strings.Split(name, ".")
	if len(spl) == 1 {
		spl[0] = strings.Trim(spl[0], " `")
		if val, ok := b.fldTag[spl[0]]; ok {
			return val
		}
	} else {
		spl[0] = strings.Trim(spl[0], " `")
		spl[1] = strings.Trim(spl[1], " `")
		if val, ok := b.fldTag[spl[1]]; ok {
			return spl[0] + "." + val
		}
	}
	return name
}

func (b *QueryBuilder) Field(name string) SqlReserved {
	return Sql(b.extractName(name))
}
func (b *QueryBuilder) Fill(values ...*OBJ) Builder {
	b.values = values
	return b
}

func (b *QueryBuilder) Columns(columns ...string) Builder {
	for _, column := range columns {
		b.columns = append(b.columns, b.extractName(column))
	}
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
	join := "JOIN " + table + " USING(" + b.extractName(using) + ")"
	b.joins = append(b.joins, join)
	return b
}

func (b *QueryBuilder) BitwiseAnd(field string, with int64, value int64) Builder {
	b.ops = append(b.ops, b.stp)
	b.wheres = append(b.wheres, fmt.Sprintf("%s & %v = %v", b.extractName(field), with, value))
	return b
}

func (b *QueryBuilder) BitwiseOr(field string, with int64, value int64) Builder {
	b.ops = append(b.ops, b.stp)
	b.wheres = append(b.wheres, fmt.Sprintf("%s | %v = %v", b.extractName(field), with, value))
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
	b.wheres = append(b.wheres, fmt.Sprintf("%s = %s", b.extractName(field), Convert(value)))
	return b
}
func (b *QueryBuilder) WhereLike(field string, value interface{}) Builder {
	b.ops = append(b.ops, b.stp)
	b.wheres = append(b.wheres, fmt.Sprintf("%s LIKE %s", b.extractName(field), Convert(value)))
	return b
}
func (b *QueryBuilder) WhereNotLike(field string, value interface{}) Builder {
	b.ops = append(b.ops, b.stp)
	b.wheres = append(b.wheres, fmt.Sprintf("%s NOT LIKE %s", b.extractName(field), Convert(value)))
	return b
}
func (b *QueryBuilder) Find(value interface{}) Builder {
	b.ops = append(b.ops, b.stp)
	b.wheres = append(b.wheres, fmt.Sprintf("%s = %s", b.extractName("id"), Convert(value)))
	return b
}
func (b *QueryBuilder) WhereNot(field string, value interface{}) Builder {
	b.ops = append(b.ops, b.stp)
	b.wheres = append(b.wheres, fmt.Sprintf("%s != %s", b.extractName(field), Convert(value)))
	return b
}

func (b *QueryBuilder) WhereNull(field string) Builder {
	b.ops = append(b.ops, b.stp)
	b.wheres = append(b.wheres, fmt.Sprintf("%s IS NULL", b.extractName(field)))
	return b
}

func (b *QueryBuilder) WhereNotNull(field string) Builder {
	b.ops = append(b.ops, b.stp)
	b.wheres = append(b.wheres, fmt.Sprintf("%s IS NOT NULL", b.extractName(field)))
	return b
}

func (b *QueryBuilder) WhereBetween(field string, value1 interface{}, value2 interface{}) Builder {
	b.ops = append(b.ops, b.stp)
	b.wheres = append(b.wheres, fmt.Sprintf("%s BETWEEN %s AND %s", b.extractName(field), Convert(value1), Convert(value2)))
	return b
}

func (b *QueryBuilder) WhereGT(field string, value interface{}) Builder {
	b.ops = append(b.ops, b.stp)
	b.wheres = append(b.wheres, fmt.Sprintf("%s > %s", b.extractName(field), Convert(value)))
	return b
}
func (b *QueryBuilder) WhereGTE(field string, value interface{}) Builder {
	b.ops = append(b.ops, b.stp)
	b.wheres = append(b.wheres, fmt.Sprintf("%s >= %s", b.extractName(field), Convert(value)))
	return b
}

func (b *QueryBuilder) WhereLT(field string, value interface{}) Builder {
	b.ops = append(b.ops, b.stp)
	b.wheres = append(b.wheres, fmt.Sprintf("%s < %s", b.extractName(field), Convert(value)))
	return b
}
func (b *QueryBuilder) WhereLTE(field string, value interface{}) Builder {
	b.ops = append(b.ops, b.stp)
	b.wheres = append(b.wheres, fmt.Sprintf("%s <= %s", b.extractName(field), Convert(value)))
	return b
}

func (b *QueryBuilder) WhereIn(field string, value []interface{}) Builder {
	b.ops = append(b.ops, b.stp)
	b.wheres = append(b.wheres, fmt.Sprintf("%s in %s", b.extractName(field), Convert(value)))
	return b
}

func (b *QueryBuilder) WhereInQuery(field string, fn func(b Builder)) Builder {
	b.ops = append(b.ops, b.stp)
	bld := &QueryBuilder{
		typ: SqlTypRead,
	}
	fn(bld)
	b.wheres = append(b.wheres, fmt.Sprintf("%s IN (%s)", b.extractName(field), bld.Query()))
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
				itm := (*item)[key]
				values += Convert(itm)
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
			values += key + "=" + Convert(value)
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
	}
	return
}

func (b *QueryBuilder) Use(a interface{}) Builder {
	b.u = a
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
	b.Scan(o)
	if b.fln < 1 {
		b.err = fmt.Errorf("no row found")
	}
	return b
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
	type LenObj struct {
		Len int64 `json:"len"`
	}
	var obj LenObj
	a := Read("("+b.Query()+") a").Columns("COUNT(*) len").Use(b.u).Scan(&obj)
	b.err = a.GetError()
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
func (b *QueryBuilder) Set(key string, val interface{}) (out Builder) {
	out = b
	if len(b.values) == 0 {
		b.values = []*OBJ{{}}
	}
	for _, value := range b.values {
		(*value)[b.extractName(key)] = val
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

func checkEqual(tag, json, fieldName, name string) bool {
	if tag != name {
		jsn := json
		if jsn != name {
			if strings.ToLower(name) != fieldName {
				return false
			}
		}
	}
	return true
}

func (b *QueryBuilder) getStructFields(elemType reflect.Type, mode int, keys ...string) (out map[string]int) {
	out = make(map[string]int)
	b.fldTag = make(map[string]string)

	fln := elemType.NumField()
	for j := 0; j < fln; j++ {

		field := elemType.Field(j)
		tag := field.Tag.Get("gql")
		jsn := field.Tag.Get("json")
		fld := strings.ToLower(field.Name)
		if jsn != "" {
			b.fldTag[jsn] = tag
		}
		b.fldTag[field.Name] = tag
		b.fldTag[fld] = tag
		allow := true

		if mode == 1 { // only
			if len(keys) > 0 {
				allow = some(keys, func(key string) bool {
					return checkEqual(tag, jsn, fld, key)
				})
			} else {
				allow = false
			}
		} else if mode == 2 { // exclude
			if len(keys) > 0 {
				allow = !some(keys, func(key string) bool {
					return checkEqual(tag, jsn, fld, key)
				})
			}
		}
		if tag != "-" && tag != "" && allow && tag != "id" {
			out[tag] = j
		} else {
			out[tag] = -1
		}
	}
	return
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
	tf := reflect.TypeOf(o).Elem()

	if tf.Kind() == reflect.Slice {
		elemType := tf.Elem()
		if elemType.Kind() != reflect.Struct {
			elemType = elemType.Elem()
		}
		structFields := b.getStructFields(elemType, mode, keys...)

		b.values = make([]*OBJ, 0)
		ln := vf.Len()
		for i := 0; i < ln; i++ {
			val := vf.Index(i)
			if val.Kind() != reflect.Struct {
				val = val.Elem()
			}
			data := make(OBJ)
			for key, j := range structFields {
				if j > -1 {
					value := val.Field(j)
					if value.CanInterface() {
						data[key] = value.Interface()
					}
				}
			}
			b.values = append(b.values, &data)
		}
		return
	}
	if tf.Kind() != reflect.Struct {
		tf = tf.Elem()
		vf = vf.Elem()
	}

	structFields := b.getStructFields(tf, mode, keys...)

	data := make(OBJ)
	for key, j := range structFields {
		if j > -1 {
			value := vf.Field(j)
			if value.CanInterface() {
				data[key] = value.Interface()
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
					ifc[i] = obj
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
	if err != nil {
		return
	}
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
				if tags.Get("gql") == "id" {
					val.Set(reflect.ValueOf(b.lastInsertedId))
					break
				}
			}
		}
		recover()
	}
	b.rowsAffected, err = a.RowsAffected()
	return
}
