package helper

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"
)

/**
 * @author LTNB (baolam0307@gmail.com)
 * @since
 *
 */

func GetOneRow(i interface{}, tableName string, db *sql.DB) (*sql.Rows, error) {
	id := reflect.ValueOf(i).Elem().FieldByName("Id")
	sql := fmt.Sprintf("SELECT * FROM %s WHERE id = '%s' LIMIT 1", tableName, id)
	return queryWithContext(sql, db)
}

func GetAllRows(tableName string, db *sql.DB) (*sql.Rows, error) {
	sql := fmt.Sprintf("SELECT * FROM %s", tableName)
	return queryWithContext(sql, db)
}

func Create(bo interface{}, tableName string, db *sql.DB) (sql.Result, error) {
	data := colsStructMapping(reflect.TypeOf(bo).Elem(), reflect.ValueOf(bo).Elem(), make(map[string]interface{}))
	keys, values := insertBuilder(data)
	sql := fmt.Sprintf("INSERT INTO  %s %s VALUES %s ", tableName, keys, values)
	return execWithContext(sql, db)
}

func Update(bo interface{}, tableName string, db *sql.DB) (sql.Result, error) {
	data := colsStructMapping(reflect.TypeOf(bo), reflect.ValueOf(bo), make(map[string]interface{}))
	updateFields := updateBuilder(data)
	sql := fmt.Sprintf("UPDATE  %s SET %s ", tableName, updateFields)
	return execWithContext(sql, db)
}

func Delete(id string, tableName string, db *sql.DB) (sql.Result, error) {
	sql := fmt.Sprintf("DELETE FROM  %s WHERE id ='%s' ", tableName, id)
	return execWithContext(sql, db)
}

func insertBuilder(m map[string]interface{}) (string, string) {
	keys := "("
	values := "("
	for k, v := range m {
		if k == "created_at" {
			v = time.Now().Nanosecond()
		}
		keys = keys + k + ","
		if reflect.TypeOf(v).Kind() == reflect.String {
			values = values + fmt.Sprintf("'%v'", v) + ","
		} else {
			values = values + fmt.Sprintf("%v", v) + ","
		}
	}
	return keys[:(len(keys)-1)] + ")", values[:(len(values)-1)] + ")"
}

// update table set ....
func updateBuilder(m map[string]interface{}) string {
	result := ""
	for k, v := range m {
		if k == "updated_at" {
			v = time.Now().Nanosecond()
		}
		if k == "id" {
			continue
		}
		if reflect.TypeOf(v).Kind() == reflect.String {
			result = result + fmt.Sprintf("%v='%v',", k, v)
		} else {
			result = result + fmt.Sprintf("%v=%v,", k, v)
		}
	}
	result = result[:(len(result) - 1)]
	result = result + fmt.Sprintf(" WHERE id = '%s'", m["id"])
	return result
}

func colsStructMapping(t reflect.Type, v reflect.Value, result map[string]interface{}) map[string]interface{} {
	numField := v.NumField()
	for i := 0; i < numField; i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct {
			result = colsStructMapping(field.Type, v.Field(i), result)
		} else {
			col := strings.Split(field.Tag.Get("sql"), ",")[0]
			switch field.Type.Kind() {
			case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
				result[col] = v.Field(i).Int()
				break
			case reflect.String:
				result[col] = v.Field(i).String()
				break
			case reflect.Bool:
				result[col] = v.Field(i).Bool()
				break
			case reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
				result[col] = v.Field(i).Uint()
				break
			}

		}
	}
	return result
}


func RowsToStruct(rows *sql.Rows, i interface{}, tagName string) {
	m, _ := RowToMap(rows)
	MapToStruct(m, tagName, i)
}

func RowToMap(rows *sql.Rows) (map[string]interface{}, error) {
	cols, _ := rows.Columns()
	columns := make([]interface{}, len(cols))
	columnPointers := make([]interface{}, len(cols))
	for i := range columns {
		columnPointers[i] = &columns[i]
	}

	if err := rows.Scan(columnPointers...); err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	for i, colName := range cols {
		val := columnPointers[i].(*interface{})
		m[colName] = *val
	}
	return m, nil
}

func queryWithContext(sql string, db *sql.DB) (*sql.Rows, error) {
	ctx := context.Background()
	rows, err := db.QueryContext(ctx, sql)
	return rows, err
}

func execWithContext(sql string, db *sql.DB) (sql.Result, error) {
	ctx := context.Background()
	return db.ExecContext(ctx, sql)
}


//TODO: need to use struct, not use pointer
func MapToStruct(source map[string]interface{}, tagName string, target interface{}) {
	val := reflect.ValueOf(target).Elem()
	numField := val.NumField()
	for i := 0; i < numField; i++ {
		field := val.Field(i)
		if !field.CanSet() {
			continue
		}

		if field.Kind() == reflect.Struct {
			for j := 0; j < field.NumField(); j++ {
				fieldName := strings.Split(field.Type().Field(j).Tag.Get(tagName), ",")[0]
				if item, ok := source[fieldName]; ok {
					if field.Field(j).CanSet() {
						typeMapping(item, field.Field(j))
					}
				}
			}
			continue
		}
		fieldName := strings.Split(val.Type().Field(i).Tag.Get(tagName), ",")[0]
		if item, ok := source[fieldName]; ok {
			typeMapping(item, val.Field(i))
		}
	}
}


func typeMapping(item interface{}, field reflect.Value) {
	if item != nil {
		switch field.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			field.SetInt(item.(int64))
		case reflect.String:
			field.SetString(item.(string))
		case reflect.Float32, reflect.Float64:
			field.SetFloat(item.(float64))
		case reflect.Bool:
			field.SetBool(item.(bool))
		case reflect.Ptr:
			if reflect.ValueOf(item).Kind() == reflect.Bool {
				itemBool := item.(bool)
				field.Set(reflect.ValueOf(&itemBool))
			}
		case reflect.Struct:
			field.Set(reflect.ValueOf(item))
		}
	}
}
