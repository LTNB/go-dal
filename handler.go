package go_dal

import (
	"database/sql"
	"reflect"
)

/**
 * @author LTNB (baolam0307@gmail.com)
 * @since
 *
 */
type IDatabaseHelper interface {
	GetDatabase() *sql.DB
	TypeMapping(data interface{}, field reflect.Value)
}