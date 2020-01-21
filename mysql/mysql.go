package mysql

import (
	"database/sql"
	go_dal "github.com/LTNB/go-dal"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"strconv"
)

/**
 * @author LTNB (baolam0307@gmail.com)
 * @since
 *
 */

type Helper struct {
}

func (mysqlHelper Helper) GetDatabase() *sql.DB {
	return go_dal.GetDatabase()
}

func (mysqlHelper Helper) TypeMapping(data interface{}, field reflect.Value) {
	if data != nil {
		dataAsString := string(data.([]byte))
		switch field.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			temp, _ :=strconv.Atoi(dataAsString)
			field.SetInt(int64(temp))
		case reflect.String:
			field.SetString(dataAsString)
		case reflect.Float32, reflect.Float64:
			temp ,_:= strconv.ParseFloat(dataAsString, 10)
			field.SetFloat(temp)
		case reflect.Bool:
			if dataAsString == "1" {
				field.SetBool(true)
			} else {
				field.SetBool(false)
			}
		case reflect.Struct:
			field.Set(reflect.ValueOf(data))
		}
	}
}
