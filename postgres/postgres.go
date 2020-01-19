package postgres

import (
	"database/sql"
	godal "github.com/LTNB/go-dal"
	helper "github.com/LTNB/go-dal/helper/sql"
	_ "github.com/lib/pq"
)

/**
 * @author LTNB (baolam0307@gmail.com)
 * @since
 *
 */

type Helper struct {
	db             *sql.DB
	TableName      string
	Bo             interface{}
	DefaultTagName string
}

func (postgresHelper *Helper) Init() {
	if postgresHelper.db == nil {
		postgresHelper.db = godal.GetDatabase()
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

func (postgresHelper Helper) GetOneByConditions(bo interface{}, conditions map[string]interface{}, tagName string) error {
	rows, err := helper.GetOneRowByConditions(conditions, postgresHelper.TableName, godal.GetDatabase())
	rows.Next()
	helper.RowsToStruct(rows, bo, tagName)
	defer rows.Close()
	return err
}

func (postgresHelper Helper) GetOneAsMap(bo interface{}) (map[string]interface{}, error) {
	rows, err := postgresHelper.getOne(bo)
	if err != nil {
		return nil, err
	}
	m, err := helper.RowToMap(rows)
	defer rows.Close()
	return m, err
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
	return helper.GetAllAsMap(postgresHelper.TableName, godal.GetDatabase())
}

func (postgresHelper Helper) GetByConditions(conditions map[string]interface{}, orderBy map[string]string, limit, offset int, tagName string) ([]interface{}, error) {
	if tagName == "" {
		tagName = postgresHelper.DefaultTagName
	}
	return helper.GetByConditions(postgresHelper.Bo, conditions, orderBy, limit, offset, postgresHelper.TableName, tagName, godal.GetDatabase())
}

func (postgresHelper Helper) GetByConditionsAsMap(conditions map[string]interface{}, orderBy map[string]string, limit, offset int, tagName string) ([]map[string]interface{}, error) {
	if tagName == "" {
		tagName = postgresHelper.DefaultTagName
	}
	return helper.GetByConditionsAsMap(conditions, orderBy, limit, offset, postgresHelper.TableName, tagName, godal.GetDatabase())
}

func (postgresHelper Helper) getAllAsInterface(tagName string) ([]interface{}, error) {
	return helper.GetAllAsInterface(postgresHelper.Bo, postgresHelper.TableName, tagName, godal.GetDatabase())
}

//=========Create==========
func (postgresHelper Helper) Create(bo interface{}) (sql.Result, error) {
	db := godal.GetDatabase()
	return helper.Create(bo, postgresHelper.TableName, db)

}
func (postgresHelper Helper) CreateByTag(bo interface{}, tagName string) (sql.Result, error) {
	db := godal.GetDatabase()
	return helper.CreateByTag(bo, postgresHelper.TableName, db, tagName)
}

//=========Update==========
func (postgresHelper Helper) Update(bo interface{}) (sql.Result, error) {
	db := godal.GetDatabase()
	return helper.Update(bo, postgresHelper.TableName, db)
}

func (postgresHelper Helper) UpdateByTag(bo interface{}, tagName string) (sql.Result, error) {
	db := godal.GetDatabase()
	return helper.UpdateByTag(bo, postgresHelper.TableName, db, tagName)
}

//=========Delete==========
func (postgresHelper Helper) DeleteByConditions(conditions map[string]interface{}) (sql.Result, error) {
	db := godal.GetDatabase()
	return helper.DeleteByConditions(conditions, postgresHelper.TableName, db)
}

func (postgresHelper Helper) Delete(bo interface{}) (sql.Result, error) {
	db := godal.GetDatabase()
	return helper.Delete(bo, postgresHelper.TableName, db)
}
