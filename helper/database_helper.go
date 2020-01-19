package helper

import "database/sql"

/**
 * @author LTNB (baolam0307@gmail.com)
 * @since
 *
 */

type IDatabaseHelper interface {
	GetOne(bo interface{}) error
	GetOneByTag(bo interface{}, tagName string) error
	GetOneByConditions(bo interface{}, conditions map[string]interface{}, tagName string) error
	GetOneAsMap(bo interface{}) (map[string]interface{}, error)
	GetAll() ([]interface{}, error)
	GetAllByTag(tagName string) ([]interface{}, error)
	GetAllAsMap() ([]map[string]interface{}, error)
	GetByConditions(conditions map[string]interface{}, orderBy map[string]string, limit, offset int, tagName string) ([]interface{}, error)
	GetByConditionsAsMap(conditions map[string]interface{}, orderBy map[string]string, limit, offset int, tableName, tagName string, db *sql.DB) ([]map[string]interface{}, error)
	Create(bo interface{}) (int64, error)
	CreateByTag(bo interface{}, tagName string) (sql.Result, error) //TODO testing
	Update(bo interface{}) (int64, error) //TODO testing
	UpdateByTag(bo interface{}, tagName string) (sql.Result, error) //TODO testing
	Delete(id string) (int64, error)
}