package mysql

import (
	go_dal "github.com/LTNB/go-dal"
	"github.com/LTNB/go-dal/helper/sql"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

/**
 * @author LTNB (baolam0307@gmail.com)
 * @since
 *
 */

var accountHelper sql.Helper

type AccountMock struct {
	Id       string `json:"id" primary:"id" sql:"id"`
	Email    string `json:"email" sql:"email"`
	FullName string `json:"full_name" sql:"full_name"`
	Role     string `json:"role" sql:"role"`
	Active   bool   `json:"active" sql:"active"`
}

func setup() {
	conf := go_dal.Config{
		DriverName:     "mysql",
		DataSourceName: "root:12345678@tcp(127.0.0.1:3306)/template",
		MaxOpenConns:   5,
		MaxLifeTime:    1 * time.Minute,
		MaxIdleConns:   5,
	}

	conf.Init()
	aHelper := sql.Helper{
		TableName: "account",
		Bo:        AccountMock{},
		DefaultTagName: "json",
		Handler: Helper{},
	}
	accountHelper = aHelper
}

func TestConnection(t *testing.T) {
	db := go_dal.GetDatabase()
	err := db.Ping()
	assert.Nil(t, err, "connected")

}

func TestGetOne(t *testing.T) {
	account := AccountMock{
		Id:       "1",
		Email:    "baolam0307@gmail.com",
		FullName: "Ta Ngoc Bao Lam",
		Role:     "admin",
		Active:   true,
	}
	accountHelper.Create(account)
	bo := AccountMock{Id: "1"}
	accountHelper.GetOne(&bo)
	assert.Equal(t, "baolam0307@gmail.com", bo.Email, "success")
	accountHelper.Delete(bo)
}

func TestGetOneByTag(t *testing.T) {
	account := AccountMock{
		Id:       "1",
		Email:    "baolam0307@gmail.com",
		FullName: "Ta Ngoc Bao Lam",
		Role:     "admin",
		Active:   true,
	}
	accountHelper.Create(account)
	bo := AccountMock{Id: "1"}
	accountHelper.GetOneByTag(&bo, "json")
	assert.Equal(t, "baolam0307@gmail.com", bo.Email, "success")
	accountHelper.Delete(bo)
}

func TestGetOneByConditions(t *testing.T) {
	account := AccountMock{
		Id:       "1",
		Email:    "baolam0307@gmail.com",
		FullName: "Ta Ngoc Bao Lam",
		Role:     "admin",
		Active:   true,
	}
	accountHelper.Create(account)
	bo := AccountMock{}
	conditions := make(map[string]interface{})
	conditions["role"] = "admin"
	accountHelper.GetOneByConditions(&bo, conditions, "json")
	assert.Equal(t, "baolam0307@gmail.com", bo.Email, "success")
	accountHelper.Delete(bo)
}

func TestGetAsMap(t *testing.T) {
	account := AccountMock{
		Id:       "1",
		Email:    "baolam0307@gmail.com",
		FullName: "Ta Ngoc Bao Lam",
		Role:     "admin",
		Active:   true,
	}
	accountHelper.Create(account)
	bo := AccountMock{Id: "1"}
	result, err := accountHelper.GetOneAsMap(&bo)
	assert.Equal(t, "baolam0307@gmail.com", string(result["email"].([]byte)), "success")
	assert.Nil(t, err, "success")
	accountHelper.Delete(bo)
}

func TestGetAll(t *testing.T) {
	account := AccountMock{
		Id:       "1",
		Email:    "baolam0307@gmail.com",
		FullName: "Ta Ngoc Bao Lam",
		Role:     "admin",
		Active:   true,
	}
	account1 := AccountMock{
		Id:       "2",
		Email:    "lamtnb@gmail.com",
		FullName: "lamtnb",
		Role:     "user",
		Active:   false,
	}
	accountHelper.Create(account)
	accountHelper.Create(account1)
	result, err := accountHelper.GetAll()

	assert.Equal(t, 2, len(result), "success")
	assert.Nil(t, err, "success")
	accountHelper.Delete(account)
	accountHelper.Delete(account1)
}

func TestGetAllByTag(t *testing.T) {
	account := AccountMock{
		Id:       "1",
		Email:    "baolam0307@gmail.com",
		FullName: "Ta Ngoc Bao Lam",
		Role:     "admin",
		Active:   true,
	}
	account1 := AccountMock{
		Id:       "2",
		Email:    "lamtnb@gmail.com",
		FullName: "lamtnb",
		Role:     "user",
		Active:   false,
	}
	accountHelper.Create(account)
	accountHelper.Create(account1)
	result, err := accountHelper.GetAllByTag("json")

	assert.Equal(t, 2, len(result), "success")
	assert.Nil(t, err, "success")
	accountHelper.Delete(account)
	accountHelper.Delete(account1)
}

func TestGetAllAsMap(t *testing.T) {
	account := AccountMock{
		Id:       "1",
		Email:    "baolam0307@gmail.com",
		FullName: "Ta Ngoc Bao Lam",
		Role:     "admin",
		Active:   true,
	}
	account1 := AccountMock{
		Id:       "2",
		Email:    "lamtnb@gmail.com",
		FullName: "lamtnb",
		Role:     "user",
		Active:   false,
	}
	accountHelper.Create(account)
	accountHelper.Create(account1)
	result, err := accountHelper.GetAllAsMap()
	assert.Equal(t, "baolam0307@gmail.com", string(result[0]["email"].([]byte)), "success")
	assert.Equal(t, "lamtnb@gmail.com", string(result[1]["email"].([]byte)), "success")
	assert.Nil(t, err, "success")
	accountHelper.Delete(account)
	accountHelper.Delete(account1)
}

func TestGetByConditions(t *testing.T) {
	account := AccountMock{
		Id:       "1",
		Email:    "baolam0307@gmail.com",
		FullName: "Ta Ngoc Bao Lam",
		Role:     "admin",
		Active:   true,
	}
	account1 := AccountMock{
		Id:       "2",
		Email:    "lamtnb@gmail.com",
		FullName: "lamtnb",
		Role:     "user",
		Active:   false,
	}
	accountHelper.Create(account)
	accountHelper.Create(account1)
	defer accountHelper.Delete(account)
	defer accountHelper.Delete(account1)

	conditions := make(map[string]interface{})
	conditions["role"] = "admin"
	orderBy := make(map[string]string)
	orderBy["full_name"] = "ASC"
	limit := 1
	offset := 0
	result, err := accountHelper.GetByConditions(conditions, orderBy, limit, offset, "")

	assert.Equal(t, 1, len(result), "success")
	assert.Nil(t, err, "success")
}

func TestGetByConditionsAsMap(t *testing.T) {
	account := AccountMock{
		Id:       "1",
		Email:    "baolam0307@gmail.com",
		FullName: "Ta Ngoc Bao Lam",
		Role:     "admin",
		Active:   true,
	}
	account1 := AccountMock{
		Id:       "2",
		Email:    "lamtnb@gmail.com",
		FullName: "lamtnb",
		Role:     "user",
		Active:   false,
	}
	accountHelper.Create(account)
	accountHelper.Create(account1)
	defer accountHelper.Delete(account)
	defer accountHelper.Delete(account1)

	orderBy := make(map[string]string)
	orderBy["full_name"] = "ASC"
	limit := 1
	offset := 0
	result, err := accountHelper.GetByConditionsAsMap(nil, orderBy, limit, offset, "")

	assert.Equal(t, 1, len(result), "success")
	assert.Equal(t, "lamtnb@gmail.com", string(result[0]["email"].([]byte)) , "success")
	assert.Nil(t, err, "success")
}

func TestCreateAndDelete(t *testing.T) {
	account := AccountMock{
		Id:       "1",
		Email:    "baolam0307@gmail.com",
		FullName: "Ta Ngoc Bao Lam",
		Role:     "admin",
		Active:   true,
	}
	_, err := accountHelper.Create(account)
	assert.Nil(t, err, "err must be nil")
	conditions := make(map[string]interface{})
	conditions["email"] = "baolam0307@gmail.com"
	_, err = accountHelper.DeleteByConditions(conditions)
	assert.Nil(t, err, "err must be nil")
}

func TestCreateByTagAndDelete(t *testing.T) {
	account := AccountMock{
		Id:       "1",
		Email:    "baolam0307@gmail.com",
		FullName: "Ta Ngoc Bao Lam",
		Role:     "admin",
		Active:   true,
	}
	_, err := accountHelper.CreateByTag(account, "json")
	assert.Nil(t, err, "err must be nil")
	conditions := make(map[string]interface{})
	conditions["email"] = "baolam0307@gmail.com"
	_, err = accountHelper.DeleteByConditions(conditions)
	assert.Nil(t, err, "err must be nil")
}

func TestUpdate(t *testing.T){
	account := AccountMock{
		Id:       "1",
		Email:    "baolam0307@gmail.com",
		FullName: "Ta Ngoc Bao Lam",
		Role:     "admin",
		Active:   true,
	}
	accountHelper.Create(account)
	account = AccountMock{
		Id:       "1",
		Email:    "lamtnb@scommerce.asia",
		FullName: "Ta Ngoc Bao Lam",
		Role:     "admin",
		Active:   true,
	}
	_, err := accountHelper.Update(account)
	assert.Nil(t, err, "success")
	account1 := AccountMock{Id:"1"}
	err = accountHelper.GetOne(&account1)
	assert.Nil(t, err, "success")
	assert.Equal(t, "lamtnb@scommerce.asia", account1.Email, "success")
	accountHelper.Delete(account1)
}


func TestUpdateByTag(t *testing.T){
	account := AccountMock{
		Id:       "1",
		Email:    "baolam0307@gmail.com",
		FullName: "Ta Ngoc Bao Lam",
		Role:     "admin",
		Active:   true,
	}
	accountHelper.Create(account)
	account = AccountMock{
		Id:       "1",
		Email:    "lamtnb@scommerce.asia",
		FullName: "Ta Ngoc Bao Lam",
		Role:     "admin",
		Active:   true,
	}
	_, err := accountHelper.UpdateByTag(account, "sql")
	assert.Nil(t, err, "success")
	account1 := AccountMock{Id:"1"}
	err = accountHelper.GetOne(&account1)
	assert.Nil(t, err, "success")
	assert.Equal(t, "lamtnb@scommerce.asia", account1.Email, "success")
	accountHelper.Delete(account1)
}

func TestMain(m *testing.M) {
	setup()
	r := m.Run()
	//destroy()
	os.Exit(r)
}
