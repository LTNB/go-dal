package postgres

import (
	"database/sql"
	go_dal "github.com/LTNB/go-dal"
	_ "github.com/lib/pq"
	"reflect"
)

/**
 * @author LTNB (baolam0307@gmail.com)
 * @since
 *
 */

type Helper struct {
	DatetimeFormat string
}

func (helper Helper) GetDatabase() *sql.DB {
	return go_dal.GetDatabase()
}

//data type from sql ==> datatype from struct
func (helper Helper) TypeMapping(data interface{}, field reflect.Value) {
	field.Set(reflect.ValueOf(data))
}