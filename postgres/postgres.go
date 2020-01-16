package postgres

import (
	"database/sql"
	go_dal "github.com/LTNB/go-dal"
	"github.com/LTNB/go-dal/helper"
	_ "github.com/lib/pq"
	"reflect"
)

/**
 * @author LTNB (baolam0307@gmail.com)
 * @since
 *
 */

type Helper struct {
	db        *sql.DB
	TableName string
	Bo        interface{}
}

func (postgresHelper *Helper) Init() {
	if postgresHelper.db == nil {
		postgresHelper.db = go_dal.GetDatabase()
	}
}

//=======Get one============
func (postgresHelper Helper) GetOne(bo interface{}) error {
	rows, err := postgresHelper.getOne(bo)
	helper.RowsToStruct(rows, bo, "json")
	defer rows.Close()
	return err
}

func (postgresHelper Helper) GetOneByTag(bo interface{}, tagName string) error {
	rows, err := postgresHelper.getOne(bo)
	helper.RowsToStruct(rows, bo, tagName)
	defer rows.Close()
	return err
}

func (postgresHelper Helper) GetOneAsMap(bo interface{}) (map[string]interface{}, error) {
	rows, err := postgresHelper.getOne(bo)
	if err != nil {
		return nil, err
	}
	m, err :=  helper.RowToMap(rows)
	defer rows.Close()
	return m ,err
}

func (postgresHelper Helper) getOne(bo interface{}) (*sql.Rows, error) {
	rows, err := helper.GetOneRow(bo, postgresHelper.TableName, postgresHelper.db)
	rows.Next()
	return rows, err
}

//=========Get all==========
func (postgresHelper Helper) GetAll() ([]interface{}, error) {
	return postgresHelper.getAllAsInterface("json")
}

func (postgresHelper Helper) GetAllByTag(tagName string) ([]interface{}, error) {
	return postgresHelper.getAllAsInterface(tagName)
}

func (postgresHelper Helper) GetAllAsMap() ([]map[string]interface{}, error) {
	return postgresHelper.getAllAsMap()
}

func (postgresHelper Helper) getAllAsMap() ([]map[string]interface{}, error) {
	rows, err := postgresHelper.getAllRows()
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	result := make([]map[string]interface{}, 0)
	temp := make(map[string]interface{})
	for rows.Next() {
		temp, err = helper.RowToMap(rows)
		result = append(result, temp)
	}
	return result, err
}

func (postgresHelper Helper) getAllAsInterface(tagName string) ([]interface{}, error) {
	rows, err := postgresHelper.getAllRows()
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	result := make([]interface{}, 0)
	temp := make(map[string]interface{})
	for rows.Next() {
		temp, _ = helper.RowToMap(rows)
		i := reflect.New(reflect.TypeOf(postgresHelper.Bo))
		helper.MapToStruct(temp, tagName, i.Interface())
		result = append(result, i.Interface())
	}
	return result, err
}

func (postgresHelper Helper) getAllRows() (*sql.Rows, error) {
	db := go_dal.GetDatabase()
	rows, error := helper.GetAllRows(postgresHelper.TableName, db)
	return rows, error
}

//=========Create==========
func (postgresHelper Helper) Create(bo interface{}) (int64, error) {
	db := go_dal.GetDatabase()
	result, err := helper.Create(bo, postgresHelper.TableName, db)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//=========Update==========
func (postgresHelper Helper) Update(bo interface{}) (int64, error) {
	db := go_dal.GetDatabase()
	result, err := helper.Update(bo, postgresHelper.TableName, db)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//=========Delete==========
func (postgresHelper Helper) Delete(id string) (int64, error) {
	db := go_dal.GetDatabase()
	result, err := helper.Delete(id, postgresHelper.TableName, db)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
