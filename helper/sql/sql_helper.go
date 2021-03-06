package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	goDal "github.com/LTNB/go-dal"
	helper2 "github.com/LTNB/go-dal/helper"
	"reflect"
	"strings"
	"time"
)

/**
 * @author LTNB (baolam0307@gmail.com)
 * @since
 *
 */
type Helper struct {
	Handler        goDal.IDatabaseHelper
	TableName      string
	Bo             interface{}
	DefaultTagName string
}

func (helper Helper) GetDatabase() *sql.DB {
	return helper.Handler.GetDatabase()
}

//impl database_helper

func (helper Helper) GetOne(bo interface{}) error {
	rows, err := getOneRow(bo, helper.TableName, helper.Handler.GetDatabase())
	if rows == nil || err != nil {
		goto END
	}
	rows.Next()
	helper.rowsToStruct(rows, bo, "json")
	err = rows.Close()
END:
	return err
}

func (helper Helper) GetOneByTag(bo interface{}, tagName string) error {
	rows, err := getOneRow(bo, helper.TableName, helper.Handler.GetDatabase())
	if err != nil {
		goto END
	}
	rows.Next()
	helper.rowsToStruct(rows, bo, tagName)
	defer func() {
		err = rows.Close()
	}()
END:
	return err
}

func (helper Helper) GetOneByConditions(bo interface{}, conditions map[string]interface{}, tagName string) error {
	var err error
	var rows *sql.Rows
	selectBuilder := SelectQueryBuilder{
		SelectFields: nil,
		From:         []string{helper.TableName},
		WhereClause: WhereClauseBuilder{
			Pair: conditions,
		},
		Limit: 1,
	}
	sqlQuery, err := selectBuilder.BuildSelectQuery()
	if err != nil {
		goto END
	}
	rows, err = queryWithContext(sqlQuery, helper.Handler.GetDatabase())
	if err != nil {
		goto END
	}
	rows.Next()
	helper.rowsToStruct(rows, bo, tagName)
	defer func() {
		err = rows.Close()
	}()
END:
	return err
}

func (helper Helper) GetOneAsMap(bo interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	rows, err := getOneRow(bo, helper.TableName, helper.Handler.GetDatabase())
	if err != nil {
		goto END
	}
	rows.Next()
	result, err = rowToMap(rows)
	defer func() {
		err = rows.Close()
	}()
END:
	return result, err

}

func getPrimaryKeysValues(boType reflect.Type, boValue reflect.Value, result map[string]interface{}) map[string]interface{} {
	primaryTagName := "primary"
	for i := 0; i < boType.NumField(); i++ {
		field := boType.Field(i)
		if field.Type.Kind() == reflect.Struct {
			result = getPrimaryKeysValues(field.Type, boValue.Field(i), result)
		} else {
			primaryTag := strings.Split(field.Tag.Get(primaryTagName), ",")[0]
			if primaryTag != "" {
				result[primaryTag] = boValue.Field(i).Interface()
			}
		}

	}
	return result
}

//bo = ptr
func getOneRow(bo interface{}, tableName string, db *sql.DB) (*sql.Rows, error) {
	primaryKeys := getPrimaryKeysValues(reflect.TypeOf(bo).Elem(), reflect.ValueOf(bo).Elem(), make(map[string]interface{}, 0))
	selectBuilder := SelectQueryBuilder{
		SelectFields: nil,
		From:         []string{tableName},
		WhereClause: WhereClauseBuilder{
			Pair: primaryKeys,
		},
		Limit: 1,
	}
	sqlQuery, err := selectBuilder.BuildSelectQuery()
	if err != nil {
		goto END
	}
END:
	return queryWithContext(sqlQuery, db)
}

func (helper Helper) GetAll() ([]interface{}, error) {
	return helper.GetAllByTag(helper.DefaultTagName)
}

func (helper Helper) GetAllByTag(tagName string) ([]interface{}, error) {
	rows, err := getAllRows(helper.TableName, nil, nil, 0, -1, helper.Handler.GetDatabase())
	var result []interface{}
	temp := make(map[string]interface{})
	if err != nil {
		goto END
	}
	defer func() {
		err = rows.Close()
	}()
	result = make([]interface{}, 0)
	temp = make(map[string]interface{})
	for rows.Next() {
		temp, _ = rowToMap(rows)
		i := reflect.New(reflect.TypeOf(helper.Bo))
		helper.mapToStruct(temp, tagName, reflect.ValueOf(i.Interface()).Elem())
		result = append(result, i.Interface())
	}
END:
	return result, err
}

func (helper Helper) GetAllAsMap() ([]map[string]interface{}, error) {
	rows, err := getAllRows(helper.TableName, nil, nil, 0, -1, helper.Handler.GetDatabase())
	var temp map[string]interface{}
	var result []map[string]interface{}
	if err != nil {
		goto END
	}
	defer func() {
		err = rows.Close()
	}()
	result = make([]map[string]interface{}, 0)
	temp = make(map[string]interface{})
	for rows.Next() {
		temp, err = rowToMap(rows)
		result = append(result, temp)
	}
END:
	return result, err
}

func (helper Helper) GetByConditions(conditions map[string]interface{}, orderBy map[string]string, limit, offset int, tagName string) ([]interface{}, error) {
	if tagName == "" {
		tagName = helper.DefaultTagName
	}
	rows, err := getAllRows(helper.TableName, conditions, orderBy, limit, offset, helper.Handler.GetDatabase())
	var result []interface{}
	var temp map[string]interface{}
	if err != nil {
		goto END
	}
	defer func() {
		err = rows.Close()
	}()

	result = make([]interface{}, 0)
	temp = make(map[string]interface{})
	for rows.Next() {
		temp, _ = rowToMap(rows)
		i := reflect.New(reflect.TypeOf(helper.Bo))
		helper.mapToStruct(temp, tagName, reflect.ValueOf(i.Interface()).Elem())
		result = append(result, i.Interface())
	}
END:
	return result, err
}

func (helper Helper) GetByConditionsAsMap(conditions map[string]interface{}, orderBy map[string]string, limit, offset int, tagName string) ([]map[string]interface{}, error) {
	var temp map[string]interface{}
	var result []map[string]interface{}
	if tagName == "" {
		tagName = helper.DefaultTagName
	}
	rows, err := getAllRows(helper.TableName, conditions, orderBy, limit, offset, helper.Handler.GetDatabase())
	if err != nil {
		goto END
	}
	defer func() {
		err = rows.Close()
	}()
	result = make([]map[string]interface{}, 0)
	temp = make(map[string]interface{})
	for rows.Next() {
		temp, err = rowToMap(rows)
		result = append(result, temp)
	}
END:
	return result, err
}

func (helper Helper) Create(bo interface{}) (int64, error) {
	result, err := helper.createByTag(bo, helper.TableName, helper.Handler.GetDatabase(), helper.DefaultTagName)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

/**
 * create batch boList
 */
func (helper Helper) CreateBatch(boList []interface{}) (int64, error) {
	result, err := helper.createBatchByTag(boList, helper.TableName, helper.Handler.GetDatabase(), helper.DefaultTagName)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

/**
 * create by field tag ( default json tag)
 */
func (helper Helper) CreateByTag(bo interface{}, tagName string) (int64, error) {
	result, err := helper.createByTag(bo, helper.TableName, helper.Handler.GetDatabase(), tagName)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

/**
 * update data change in bo ( default json tag )
 */
func (helper Helper) Update(bo interface{}) (int64, error) {
	result, err := helper.update(bo, helper.TableName, helper.Handler.GetDatabase(), helper.DefaultTagName)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

/**
 * update data change in bo by tag
 */
func (helper Helper) UpdateByTag(bo interface{}, tagName string) (int64, error) {
	result, err := helper.update(bo, helper.TableName, helper.Handler.GetDatabase(), tagName)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

/**
 * delete by bo condition
 */
func (helper Helper) Delete(bo interface{}) (int64, error) {
	conditions := getPrimaryKeysValues(reflect.TypeOf(bo), reflect.ValueOf(bo), make(map[string]interface{}))
	builder := DeleteBuilder{
		TableName: helper.TableName,
		WhereClause: WhereClauseBuilder{
			Pair: conditions,
		},
	}
	sqlQuery, err := builder.BuildDeleteQuery()
	if err != nil {
		return 0, err
	}
	result, err := execWithContext(sqlQuery, helper.Handler.GetDatabase())
	if err != nil {
		return 0, nil
	}
	return result.RowsAffected()
}

/**
 * delete by map conditions
 */
func (helper Helper) DeleteByConditions(conditions map[string]interface{}) (int64, error) {
	builder := DeleteBuilder{
		TableName: helper.TableName,
		WhereClause: WhereClauseBuilder{
			Pair: conditions,
		},
	}
	sqlQuery, err := builder.BuildDeleteQuery()
	if err != nil {
		return 0, err
	}
	result, err := execWithContext(sqlQuery, helper.Handler.GetDatabase())
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

/**
 * find all by condition
 * @conditions: mapping with table
 * @orderBy: map[key] = "ASC" / "DESC"
 */
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
	sqlQuery, err := builder.BuildSelectQuery()
	if err != nil {
		return nil, err
	}
	return queryWithContext(sqlQuery, db)
}

/**
 * create bo by tag
 */
func (helper Helper) createByTag(bo interface{}, tableName string, db *sql.DB, tagName string) (sql.Result, error) {
	data, err := helper.colsStructMapping(reflect.TypeOf(bo), reflect.ValueOf(bo), make(map[string]interface{}), tagName)
	//check version field existed
	_, has := reflect.TypeOf(bo).FieldByName("Version")
	if has {
		data["version"] = 0
	}
	if err != nil {
		return nil, err
	}
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
	sqlQuery, err := insertBuilder.BuildInsertQuery()
	if err != nil {
		return nil, err
	}
	return execWithContext(sqlQuery, db)
}

/**
 * create list bo by tag
 */
func (helper Helper) createBatchByTag(boList []interface{}, tableName string, db *sql.DB, tagName string) (sql.Result, error) {
	var keys []string
	keyMap := make(map[int]string)
	var batchValue [][]interface{}
	for i, bo := range boList {
		data, err := helper.colsStructMapping(reflect.TypeOf(bo), reflect.ValueOf(bo), make(map[string]interface{}), tagName)
		if err != nil {
			return nil, err
		}
		//check version field existed
		_, has := reflect.TypeOf(bo).FieldByName("Version")
		if has {
			data["version"] = 0
		}
		value := make([]interface{}, len(data))
		count := 0
		for k, v := range data {
			if i == 0 {
				keys = append(keys, k)
				value[count] = v
				keyMap[count] = k
				count++
			} else {
				for x, y := range keyMap {
					if y == k {
						value[x] = v
						break
					}
				}
			}
		}
		batchValue = append(batchValue, value)

	}

	insertBuilder := InsertBuilder{
		TableName: tableName,
		Keys:      keys,
		Values:    batchValue,
	}
	sqlQuery, err := insertBuilder.BuildInsertQuery()
	if err != nil {
		return nil, err
	}
	return execWithContext(sqlQuery, db)
}

/**
 * update bo by primary key and mapping with tagName
 */
func (helper Helper) update(bo interface{}, tableName string, db *sql.DB, tagName string) (sql.Result, error) {
	data, err := helper.colsStructMapping(reflect.TypeOf(bo), reflect.ValueOf(bo), make(map[string]interface{}), tagName)
	if err != nil {
		return nil, err
	}
	pairWhereClause := getPrimaryKeysValues(reflect.TypeOf(bo), reflect.ValueOf(bo), map[string]interface{}{})
	//check version field existed
	if _, has := reflect.TypeOf(bo).FieldByName("Version"); has {
		arr := make([]map[string]interface{}, 1)
		arr[0] = data
		if _, err := helper.checkVersion(&arr, pairWhereClause, tagName); err != nil {
			return nil, err
		}
	}

	builder := UpdateBuilder{
		TableName: tableName,
		WhereClause: WhereClauseBuilder{
			Pair: pairWhereClause,
		},
		Values: data,
	}
	sqlQuery, err := builder.BuildUpdateQuery()
	if err != nil {
		return nil, err
	}
	return execWithContext(sqlQuery, db)
}

/**
 * check version and update new version for bo
 */
func (helper Helper) checkVersion(dataList *[]map[string]interface{}, pairWhereClause map[string]interface{}, tagName string) (*[]map[string]interface{}, error) {
	var err error
	// get all rows by conditions
	rows, err := helper.GetByConditionsAsMap(pairWhereClause, nil, 0, 0, tagName)
	if err != nil {
		return nil, err
	}
	//compare version
	if len(rows) > 0 {
		for _, item := range rows {
			for _, bo := range *dataList {
				versionIP := bo["version"]
				if fmt.Sprintf("%v", item["id"]) == fmt.Sprintf("%v", bo["id"]) {
					if fmt.Sprintf("%v", versionIP) != fmt.Sprintf("%v", item["version"]) {
						return nil, errors.New(fmt.Sprintf("version of {id: %v} is not matched:  %v != %v", item["id"], item["version"], versionIP))
					}
					bo["version"] = item["version"].(int64) + 1
				}
			}
		}
	}
	return dataList, err
}

/**
 * data type from struct to data type sql
 */
func (helper Helper) colsStructMapping(t reflect.Type, v reflect.Value, result map[string]interface{}, tagName string) (map[string]interface{}, error) {
	var err error
	var bo helper2.BaseBo
	numField := v.NumField()
	for i := 0; i < numField; i++ {
		field := t.Field(i)
		value := v.Field(i)
		//get tag name in field
		tagNameStruct := strings.Split(field.Tag.Get(tagName), ",")
		if tagNameStruct == nil || len(tagNameStruct) != 1 {
			err = fmt.Errorf("tag name `%v` is empty or not unique.\n", tagName)
			goto End
		}
		col := tagNameStruct[0]
		//mapping data to map collection
		if field.Type.Kind() == reflect.Struct {
			// struct define by bz
			if strings.Split(field.Tag.Get("promoted"), ",")[0] == "true" {
				if field.Type == reflect.TypeOf(bo) {
					idType := strings.Split(field.Tag.Get("id"), ",")
					if len(idType) != 0 && idType[0] != "" {
						result = helper.initPrimary(idType[0], field, tagName, result)
						continue
					}
				}
				result, err = helper.colsStructMapping(field.Type, value, result, tagName)
			} else {
				// mapping struct define by framework
				switch field.Type {
				case reflect.TypeOf(time.Time{}):
					//call format(input) which return only one response
					res := helper.timeMapping(v.Field(i))
					result[col] = res[0].String()
					break
				}
			}
		} else {
			result[col] = v.Field(i).Interface()
		}
	}
End:
	return result, err
}

/**
 * init primary key
 * uuid :generate uuid value like 17b16d72-5496-4080-be2e-19bb84373a59
 * timestamp: get current millisecond base on system timezone
 * serial: ignore value and auto increase by database ( base on database type )
 */
func (helper Helper) initPrimary(idType string, field reflect.StructField, tagName string, result map[string]interface{}) map[string]interface{} {
	switch idType {
	case "uuid":
		newId := reflect.New(reflect.TypeOf(field))
		newId.Elem().Field(0).Set(reflect.ValueOf(helper2.UUIDGenerate()))
		col := strings.Split(field.Type.Field(0).Tag.Get(tagName), ",")[0]
		result[col] = newId.Elem().Field(0).Interface()
		break
	case "timestamp":
		newId := reflect.New(reflect.TypeOf(field))
		newId.Elem().Field(0).Set(reflect.ValueOf(time.Now().Unix() * 1000))
		col := strings.Split(field.Type.Field(0).Tag.Get(tagName), ",")[0]
		result[col] = newId.Elem().Field(0).Interface()
		break
	case "serial":
		break
	}
	return result
}

/**
 * mapping default time.Time ( RFC3339 )
 */
func (helper Helper) timeMapping(v reflect.Value) []reflect.Value {
	return v.MethodByName("Format").Call([]reflect.Value{reflect.ValueOf(time.RFC3339)})
}

/*
 * convert @rows to @i: struct depend on @tagName
 */
func (helper Helper) rowsToStruct(rows *sql.Rows, i interface{}, tagName string) {
	m, _ := rowToMap(rows)
	if m == nil {
		return
	}
	helper.mapToStruct(m, tagName, reflect.ValueOf(i).Elem())
}

/*
 * query with context
 */
//TODO ==> audit log in here...!
func queryWithContext(sql string, db *sql.DB) (*sql.Rows, error) {
	ctx := context.Background()
	rows, err := db.QueryContext(ctx, sql)
	if err != nil {
		fmt.Printf("execute query %v failed", sql)
	}
	return rows, err
}

/*
 * execute with context
 */
//TODO ==> audit log in here...!
func execWithContext(sql string, db *sql.DB) (sql.Result, error) {
	ctx := context.Background()
	result, err := db.ExecContext(ctx, sql)
	if err != nil {
		fmt.Printf("execute query %v failed", sql)
	}
	return result, err
}

/**
 * convert row to map
 */
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

/**
 * convert map to struct by tagName
 */
func (helper Helper) mapToStruct(source map[string]interface{}, tagName string, target reflect.Value) {
	numField := target.NumField()
	for i := 0; i < numField; i++ {
		field := target.Field(i)
		if !field.CanSet() {
			continue
		}
		if field.Type().Kind() == reflect.Struct && strings.Split(target.Type().Field(i).Tag.Get("promoted"), ",")[0] == "true" {
			helper.mapToStruct(source, tagName, field)
		} else {
			fieldName := strings.Split(target.Type().Field(i).Tag.Get(tagName), ",")[0]
			helper.Handler.TypeMapping(source[fieldName], target.Field(i))
		}
	}
}
