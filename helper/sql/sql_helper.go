package sql

import (
	"context"
	"database/sql"
	"fmt"
	go_dal "github.com/LTNB/go-dal"
	"reflect"
	"strings"
)

/**
 * @author LTNB (baolam0307@gmail.com)
 * @since
 *
 */
type SQLHelper struct {
	Handler        go_dal.IDatabaseHelper
	TableName      string
	Bo             interface{}
	DefaultTagName string
}

func (helper SQLHelper) GetDatabase() *sql.DB{
	return helper.Handler.GetDatabase()
}

//impl database_helper

func (helper SQLHelper) GetOne(bo interface{}) error {
	rows, err := getOneRow(bo, helper.TableName, helper.Handler.GetDatabase())
	rows.Next()
	rowsToStruct(rows, bo, "json")
	defer rows.Close()
	return err
}

func (helper SQLHelper) GetOneByTag(bo interface{}, tagName string) error {
	rows, err := getOneRow(bo, helper.TableName, helper.Handler.GetDatabase())
	rows.Next()
	rowsToStruct(rows, bo, tagName)
	defer rows.Close()
	return err
}

func (helper SQLHelper) GetOneByConditions(bo interface{}, conditions map[string]interface{}, tagName string) error {
	selectBuilder := SelectQueryBuilder{
		SelectFields: nil,
		From:         []string{helper.TableName},
		WhereClause: WhereClauseBuilder{
			Pair: conditions,
		},
		Limit: 1,
	}
	sql, _ := selectBuilder.BuildSelectQuery()
	rows, err := queryWithContext(sql, helper.Handler.GetDatabase())
	rows.Next()
	rowsToStruct(rows, bo, tagName)
	defer rows.Close()
	return err
}

func (helper SQLHelper) GetOneAsMap(bo interface{}) (map[string]interface{}, error) {
	rows, err := getOneRow(bo, helper.TableName, helper.Handler.GetDatabase())
	rows.Next()
	if err != nil {
		return nil, err
	}
	m, err := rowToMap(rows)
	defer rows.Close()
	return m, err

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

func getOneRow(bo interface{}, tableName string, db *sql.DB) (*sql.Rows, error) {
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

func (helper SQLHelper) GetAll() ([]interface{}, error) {
	return helper.GetAllByTag(helper.DefaultTagName)
}

func (helper SQLHelper) GetAllByTag(tagName string) ([]interface{}, error) {
	rows, err := getAllRows(helper.TableName, nil, nil, 0, -1, helper.Handler.GetDatabase())
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	result := make([]interface{}, 0)
	temp := make(map[string]interface{})
	for rows.Next() {
		temp, _ = rowToMap(rows)
		i := reflect.New(reflect.TypeOf(helper.Bo))
		mapToStruct(temp, tagName, i.Interface())
		result = append(result, i.Interface())
	}
	return result, err
}

func (helper SQLHelper) GetAllAsMap() ([]map[string]interface{}, error) {
	rows, err := getAllRows(helper.TableName, nil, nil, 0, -1, helper.Handler.GetDatabase())
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	result := make([]map[string]interface{}, 0)
	temp := make(map[string]interface{})
	for rows.Next() {
		temp, err = rowToMap(rows)
		result = append(result, temp)
	}
	return result, err
}

func (helper SQLHelper) GetByConditions(conditions map[string]interface{}, orderBy map[string]string, limit, offset int, tagName string) ([]interface{}, error) {
	if tagName == "" {
		tagName = helper.DefaultTagName
	}
	rows, err := getAllRows(helper.TableName, conditions, orderBy, limit, offset, helper.Handler.GetDatabase())
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	result := make([]interface{}, 0)
	temp := make(map[string]interface{})
	for rows.Next() {
		temp, _ = rowToMap(rows)
		i := reflect.New(reflect.TypeOf(helper.Bo))
		mapToStruct(temp, tagName, i.Interface())
		result = append(result, i.Interface())
	}
	return result, err
}

func (helper SQLHelper) GetByConditionsAsMap(conditions map[string]interface{}, orderBy map[string]string, limit, offset int, tagName string) ([]map[string]interface{}, error) {
	if tagName == "" {
		tagName = helper.DefaultTagName
	}
	rows, err := getAllRows(helper.TableName, conditions, orderBy, limit, offset, helper.Handler.GetDatabase())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]map[string]interface{}, 0)
	temp := make(map[string]interface{})
	for rows.Next() {
		temp, err = rowToMap(rows)
		result = append(result, temp)
	}
	return result, err
}

func (helper SQLHelper) Create(bo interface{}) (int64, error) {
	result, err := createByTag(bo, helper.TableName, helper.Handler.GetDatabase(), helper.DefaultTagName)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (helper SQLHelper) CreateByTag(bo interface{}, tagName string) (int64, error) {
	result, err := createByTag(bo, helper.TableName, helper.Handler.GetDatabase(), tagName)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (helper SQLHelper) Update(bo interface{}) (int64, error) {
	result, err := update(bo, helper.TableName, helper.Handler.GetDatabase(), helper.DefaultTagName)
	if err != nil {
		return 0, nil
	}
	return result.RowsAffected()
}

func (helper SQLHelper) UpdateByTag(bo interface{}, tagName string) (int64, error) {
	result, err := update(bo, helper.TableName, helper.Handler.GetDatabase(), tagName)
	if err != nil {
		return 0, nil
	}
	return result.RowsAffected()
}

func (helper SQLHelper) Delete(bo interface{}) (int64, error) {
	conditions, err := getPrimaryKeysValues(bo, make(map[string]interface{}))
	builder := DeleteBuilder{
		TableName: helper.TableName,
		WhereClause: WhereClauseBuilder{
			Pair: conditions,
		},
	}
	sql, err := builder.BuildDeleteQuery()
	if err != nil {
		return 0, err
	}
	result, err := execWithContext(sql, helper.Handler.GetDatabase())
	if err != nil {
		return 0, nil
	}
	return result.RowsAffected()
}

func (helper SQLHelper) DeleteByConditions(conditions map[string]interface{}) (int64, error) {
	builder := DeleteBuilder{
		TableName: helper.TableName,
		WhereClause: WhereClauseBuilder{
			Pair: conditions,
		},
	}
	sql, err := builder.BuildDeleteQuery()
	if err != nil {
		return 0, err
	}
	result, err := execWithContext(sql, helper.Handler.GetDatabase())
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
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

func rowsToStruct(rows *sql.Rows, i interface{}, tagName string) {
	m, _ := rowToMap(rows)
	mapToStruct(m, tagName, i)
}

func rowToMap(rows *sql.Rows) (map[string]interface{}, error) {
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
	result, err := db.ExecContext(ctx, sql)
	if err != nil {
		fmt.Printf("execute query %v failed", sql)
	}
	return result, err
}

//TODO: need to use struct, not use pointer
func mapToStruct(source map[string]interface{}, tagName string, target interface{}) {
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
