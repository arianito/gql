package gql

import (
	"database/sql"
	"fmt"
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
	qtyp    SqlTyp

	limit  int
	offset int
	tx     *sql.Tx
	db     *sql.DB
	//
	obj            interface{}
	lastInsertedId int64
	rowsAffected   int64
	cursor         *sql.Rows
	err            error
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
func (b *QueryBuilder) ForeignKey(localField string, remoteTable string, remoteField string, ondelete ...FKType) Builder {
	key := "FOREIGN KEY (" + localField + ") REFERENCES " + remoteTable + " (" + remoteField + ")"
	if len(ondelete) > 0 {
		if ondelete[0] == FKCascade {
			key += " ON DELETE SET NULL"
		} else if ondelete[0] == FKSetNull {
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
func (b *QueryBuilder) Count() Builder {
	b.columns = []string{Count("*", "count")}
	return b
}
func (b *QueryBuilder) Table(table string) Builder {
	b.tables = append(b.tables, table)
	return b
}

func (b *QueryBuilder) Join(table string, condition string, fn ...func(b Builder)) Builder {
	bld := &QueryBuilder{
		qtyp: SqlTypRead,
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
		qtyp: SqlTypRead,
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
		qtyp: SqlTypRead,
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
		qtyp: SqlTypRead,
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
		qtyp: SqlTypRead,
	}
	fn(bld)
	b.wheres = append(b.wheres, fmt.Sprintf("%s IN (%s)", field, bld.Query()))
	return b
}

func (b *QueryBuilder) WhereGroup(fn func(b Builder)) Builder {
	builder := &QueryBuilder{
		qtyp: SqlTypRead,
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

func (b *QueryBuilder) Query() string {
	if b.qtyp == SqlTypRead {
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
		return strings.Trim(query, " ")
	}
	if b.qtyp == SqlTypCreate {
		keys := make([]string, 0)
		for key, _ := range *b.values[0] {
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
		return "INSERT INTO " + b.tables[0] + "(" + strings.Join(keys, ", ") + ") VALUES" + strings.Join(stm, ", ")
	}

	if b.qtyp == SqlTypUpdate {
		ok := make([]string, 0)
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
		ok = append(ok, "("+values+")")
		return "UPDATE " + strings.Join(b.tables, ", ") + " SET " + values + " WHERE " + b.getWhereClauses(false)
	}
	if b.qtyp == SqlTypDelete {
		return "DELETE FROM " + strings.Join(b.tables, ", ") + " WHERE " + b.getWhereClauses(false)
	}
	if b.qtyp == SqlTypTable {
		table := "CREATE TABLE " + strings.Join(b.tables, ", ") + "(" + strings.Join(b.columns, ", ")
		if len(b.columns) > 0 && len(b.orders) > 0 {
			table += ", " + strings.Join(b.orders, ", ")
		}
		return table + ")"
	}
	return ""
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

func (b *QueryBuilder) Top(top int) Builder {
	b.limit = top
	return b
}
func (b *QueryBuilder) Offset(offset int) Builder {
	b.offset = offset
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

func (b *QueryBuilder) First(o interface{}) Builder {
	b.limit = 1
	return b.Scan(o)
}

func (b *QueryBuilder) Scan(o interface{}) Builder {
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
		rows, err := b.query()
		if err != nil {
			b.err = err
			return b
		}

		for rows.Next() {
			data, err := rows.Columns()
			if err != nil {
				rows.Close()
				b.err = err
				return b
			}

			val := reflect.New(elem)
			el:=val.Elem()
			ifc := make([]interface{}, len(data))
			for i, str := range data {
				op := pairs[str]
				if op != "" {
					obj := el.FieldByName(op).Addr().Interface()
					ifc[i] =
						obj
				}else {
					var obj interface{}
					ifc[i] =  &obj
				}
			}
			err = rows.Scan(ifc...)
			if err != nil {
				b.err = err
				return b
			}
			if stc {
				vf.Set(reflect.Append(vf, val.Elem()))
			}else{
				vf.Set(reflect.Append(vf, val))
			}
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
		rows, err := b.query()
		if err != nil {
			b.err = err
			return b
		}
		for rows.Next() {
			data, err := rows.Columns()
			if err != nil {
				rows.Close()
				b.err = err
				return b
			}
			ifc := make([]interface{}, len(data))
			for i, str := range data {
				op := pairs[str]
				if op != "" {
					obj := val.FieldByName(op).Addr().Interface()
					ifc[i] =
						obj
				}else {
					var obj interface{}
					ifc[i] =  &obj
				}
			}
			err = rows.Scan(ifc...)
			if err != nil {
				b.err = err
				return b
			}

			return b
		}
	}
	return b
}
func (b *QueryBuilder) GetError() error {
	return b.err
}
func (b *QueryBuilder) Run() Builder {

	if b.db == nil && b.tx == nil {
		panic("db driver not defined")
	}

	var a sql.Result
	var err error

	if b.db != nil {
		a, err = b.db.Exec(b.Query())
	} else {
		a, err = b.tx.Exec(b.Query())
	}
	if err != nil {
		b.err = err
		return b
	}
	b.lastInsertedId, err = a.LastInsertId()
	if err != nil {
		b.err = err
		return b
	}
	b.rowsAffected, err = a.RowsAffected()
	b.err = err
	return b
}
