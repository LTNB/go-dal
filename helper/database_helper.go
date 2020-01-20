package helper

import "database/sql"

/**
 * @author LTNB (baolam0307@gmail.com)
 * @since
 *
 */

type IDatabaseHelper interface {
	GetOne(bo interface{}) error //done
	GetOneByTag(bo interface{}, tagName string) error //done
	GetOneByConditions(bo interface{}, conditions map[string]interface{}, tagName string) error //done
	GetOneAsMap(bo interface{}) (map[string]interface{}, error) //done
	GetAll() ([]interface{}, error) //done
	GetAllByTag(tagName string) ([]interface{}, error) //done
	GetAllAsMap() ([]map[string]interface{}, error) //done
	GetByConditions(conditions map[string]interface{}, orderBy map[string]string, limit, offset int, tagName string) ([]interface{}, error) //done
	GetByConditionsAsMap(conditions map[string]interface{}, orderBy map[string]string, limit, offset int, tableName, tagName string, db *sql.DB) ([]map[string]interface{}, error) //done
	Create(bo interface{}) (int64, error) //done
	CreateByTag(bo interface{}, tagName string) (int64, error) //done
	Update(bo interface{}) (int64, error) //done
	UpdateByTag(bo interface{}, tagName string) (int64, error) //done
	Delete(id string) (int64, error) //done
	DeleteByConditions(conditions map[string]interface{}) (int64, error) // done
}