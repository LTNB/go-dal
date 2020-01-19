package sql

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

/**
 * @author LTNB (baolam0307@gmail.com)
 * @since
 *
 */

func getPrimaryKeys(bo interface{}, result []string) ([]string, error) {
	primaryTagName := "primary"
	var err error
	typ := reflect.TypeOf(bo)
	for i := 0; i < typ.Elem().NumField(); i++ {
		field := typ.Field(i)
		if field.Type.Kind() == reflect.Struct {
			result, err = getPrimaryKeys(field, result)
		}
		if field.Type.Kind() == reflect.Ptr {
			continue
		}
		primaryTag := strings.Split(field.Tag.Get(primaryTagName), ",")[0]
		if primaryTag == "" {
			continue
		}
		result = append(result, primaryTag)
	}
	return result, err
}

func getPrimaryKeysValues(bo interface{}, result map[string]interface{}) (map[string]interface{}, error) {
	primaryTagName := "primary"
	var err error
	typ := reflect.TypeOf(bo)
	val := reflect.ValueOf(bo)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Type.Kind() == reflect.Struct {
			result, err = getPrimaryKeysValues(field, result)
		}
		if field.Type.Kind() == reflect.Ptr {
			continue
		}
		primaryTag := strings.Split(field.Tag.Get(primaryTagName), ",")[0]
		if primaryTag == "" {
			continue
		}
		result[primaryTag] = val.Field(i).Interface()
	}
	return result, err
}

func GetOneRow(bo interface{}, tableName string, db *sql.DB) (*sql.Rows, error) {
	primaryKeys, _ := getPrimaryKeysValues(bo, make(map[string]interface{}, 0))
	selectBuilder := SelectQueryBuilder{
		SelectFields: nil,
		From:         []string{tableName},
		WhereClause: WhereClauseBuilder{
			Pair: primaryKeys,
		},
		Limit: 1,
	}
	sql, _ := selectBuilder.BuildSelectQuery()
	return queryWithContext(sql, db)
}

func GetOneRowByConditions(conditions map[string]interface{}, tableName string, db *sql.DB) (*sql.Rows, error) {
	selectBuilder := SelectQueryBuilder{
		SelectFields: nil,
		From:         []string{tableName},
		WhereClause: WhereClauseBuilder{
			Pair: conditions,
		},
		Limit: 1,
	}
	sql, _ := selectBuilder.BuildSelectQuery()
	return queryWithContext(sql, db)
}

func GetAllAsMap(tableName string, db *sql.DB) ([]map[string]interface{}, error) {
	rows, err := getAllRows(tableName, nil, nil, 0, -1, db)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	result := make([]map[string]interface{}, 0)
	temp := make(map[string]interface{})
	for rows.Next() {
		temp, err = RowToMap(rows)
		result = append(result, temp)
	}
	return result, err
}

func GetAllAsInterface(bo interface{}, tableName, tagName string, db *sql.DB) ([]interface{}, error) {
	rows, err := getAllRows(tableName, nil, nil, 0, -1, db)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	result := make([]interface{}, 0)
	temp := make(map[string]interface{})
	for rows.Next() {
		temp, _ = RowToMap(rows)
		i := reflect.New(reflect.TypeOf(bo))
		MapToStruct(temp, tagName, i.Interface())
		result = append(result, i.Interface())
	}
	return result, err
}

func getAllRows(tableName string, conditions map[string]interface{}, orderBy map[string]string, limit, offset int, db *sql.DB) (*sql.Rows, error) {
	builder := SelectQueryBuilder{
		From: []string{tableName},
		WhereClause: WhereClauseBuilder{
			Pair:              conditions,
			NativeWhereClause: "",
		},
		OrderBy: orderBy,
		Limit:   limit,
		Offset:  offset,
	}
	sql, err := builder.BuildSelectQuery()
	if err != nil {
		return nil, err
	}
	return queryWithContext(sql, db)
}

func GetByConditions(bo interface{}, conditions map[string]interface{}, orderBy map[string]string, limit, offset int, tableName, tagName string, db *sql.DB) ([]interface{}, error) {
	rows, err := getAllRows(tableName, conditions, orderBy, limit, offset, db)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	result := make([]interface{}, 0)
	temp := make(map[string]interface{})
	for rows.Next() {
		temp, _ = RowToMap(rows)
		i := reflect.New(reflect.TypeOf(bo))
		MapToStruct(temp, tagName, i.Interface())
		result = append(result, i.Interface())
	}
	return result, err
}

func GetByConditionsAsMap(conditions map[string]interface{}, orderBy map[string]string, limit, offset int, tableName, tagName string, db *sql.DB) ([]map[string]interface{}, error) {
	rows, err := getAllRows(tableName, conditions, orderBy, limit, offset, db)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	result := make([]map[string]interface{}, 0)
	temp := make(map[string]interface{})
	for rows.Next() {
		temp, err = RowToMap(rows)
		result = append(result, temp)
	}
	return result, err
}

func Create(bo interface{}, tableName string, db *sql.DB) (sql.Result, error) {
	return createByTag(bo, tableName, db, "json")
}

func CreateByTag(bo interface{}, tableName string, db *sql.DB, tagName string) (sql.Result, error) {
	return createByTag(bo, tableName, db, tagName)
}

func createByTag(bo interface{}, tableName string, db *sql.DB, tagName string) (sql.Result, error) {
	data := colsStructMapping(reflect.TypeOf(bo), reflect.ValueOf(bo), make(map[string]interface{}), tagName)
	var keys []string
	var value []interface{}
	for k, v := range data {
		keys = append(keys, k)
		value = append(value, v)
	}
	insertBuilder := InsertBuilder{
		TableName: tableName,
		Keys:      keys,
		Values:    [][]interface{}{value},
	}
	sql, _ := insertBuilder.BuildInsertQuery()
	return execWithContext(sql, db)
}

func update(bo interface{}, tableName string, db *sql.DB, tagName string) (sql.Result, error) {
	data := colsStructMapping(reflect.TypeOf(bo), reflect.ValueOf(bo), make(map[string]interface{}), tagName)
	pairWhereClause, err := getPrimaryKeysValues(bo, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	builder := UpdateBuilder{
		TableName: tableName,
		WhereClause: WhereClauseBuilder{
			Pair: pairWhereClause,
		},
		Values: data,
	}
	sql, err := builder.BuildUpdateQuery()
	if err != nil {
		return nil, err
	}
	return execWithContext(sql, db)
}

func UpdateByTag(bo interface{}, tableName string, db *sql.DB, tagName string) (sql.Result, error) {
	return update(bo, tableName, db, tableName)
}

func Update(bo interface{}, tableName string, db *sql.DB) (sql.Result, error) {
	return update(bo, tableName, db, "json")
}

func DeleteByConditions(conditions map[string]interface{}, tableName string, db *sql.DB) (sql.Result, error) {
	builder := DeleteBuilder{
		TableName: tableName,
		WhereClause: WhereClauseBuilder{
			Pair: conditions,
		},
	}
	sql, err := builder.BuildDeleteQuery()
	if err != nil {
		return nil, err
	}
	return execWithContext(sql, db)
}

func Delete(bo interface{}, tableName string, db *sql.DB) (sql.Result, error) {
	conditions, err := getPrimaryKeysValues(bo, make(map[string]interface{}))
	builder := DeleteBuilder{
		TableName: tableName,
		WhereClause: WhereClauseBuilder{
			Pair: conditions,
		},
	}
	sql, err := builder.BuildDeleteQuery()
	if err != nil {
		return nil, err
	}
	return execWithContext(sql, db)
}

func colsStructMapping(t reflect.Type, v reflect.Value, result map[string]interface{}, tagName string) map[string]interface{} {
	numField := v.NumField()
	for i := 0; i < numField; i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct {
			result = colsStructMapping(field.Type, v.Field(i), result, tagName)
		} else {
			col := strings.Split(field.Tag.Get(tagName), ",")[0]
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
	if err != nil {
		fmt.Printf("execute query %v failed", sql)
	}
	return rows, err
}

func execWithContext(sql string, db *sql.DB) (sql.Result, error) {
	ctx := context.Background()
	result, err :=  db.ExecContext(ctx, sql)
	if err != nil {
		fmt.Printf("execute query %v failed", sql)
	}
	return result, err
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
