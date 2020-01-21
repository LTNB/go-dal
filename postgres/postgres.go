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
}

func (postgresHelper Helper) GetDatabase() *sql.DB {
	return go_dal.GetDatabase()
}

func (postgresHelper Helper) TypeMapping(data interface{}, field reflect.Value) {
	if data != nil {
		switch field.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			field.SetInt(data.(int64))
		case reflect.String:
			field.SetString(data.(string))
		case reflect.Float32, reflect.Float64:
			field.SetFloat(data.(float64))
		case reflect.Bool:
			field.SetBool(data.(bool))
		case reflect.Ptr:
			if reflect.ValueOf(data).Kind() == reflect.Bool {
				dataBool := data.(bool)
				field.Set(reflect.ValueOf(&dataBool))
			}
		case reflect.Struct:
			field.Set(reflect.ValueOf(data))
		}
	}
}