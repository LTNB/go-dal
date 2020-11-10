package postgres

import (
	"database/sql"
	goDal "github.com/LTNB/go-dal"
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
	return goDal.GetDatabase()
}

//data type from sql ==> datatype from struct
func (helper Helper) TypeMapping(data interface{}, field reflect.Value) {
	if data == nil {
		return
	}
	field.Set(reflect.ValueOf(data))
}