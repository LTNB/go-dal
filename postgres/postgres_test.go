package postgres

import (
	go_dal "github.com/LTNB/go-dal"
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

var accountHelper Helper

type AccountMock struct {
	Id       string `json:"id" primary:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
	Active   bool   `json:"active"`
}

func setup() {
	conf := go_dal.Config{
		DriverName:     "postgres",
		DataSourceName: "postgres://lamtnb:123456@localhost:5432/template?sslmode=disable&client_encoding=UTF-8",
		MaxOpenConns:   5,
		MaxLifeTime:    1 * time.Minute,
		MaxIdleConns:   5,
	}
	conf.Init()
	aHelper := Helper{
		TableName: "account",
		Bo:        AccountMock{},
	}
	aHelper.Init()
	accountHelper = aHelper
}

func TestConnection(t *testing.T) {
	db := go_dal.GetDatabase()
	err := db.Ping()
	assert.Nil(t, err, "connected")

}

func TestSelectOne(t *testing.T) {
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

func TestSelectOneByConditions(t *testing.T) {
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

func TestSelectOneByTag(t *testing.T) {
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

func TestSelectAsMap(t *testing.T) {
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
	assert.Equal(t, "baolam0307@gmail.com", result["email"], "success")
	assert.Nil(t, err, "success")
	accountHelper.Delete(bo)
}

func TestSelectAllAsMap(t *testing.T) {
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
	assert.Equal(t, "baolam0307@gmail.com", result[0]["email"], "success")
	assert.Equal(t, "lamtnb@gmail.com", result[1]["email"], "success")
	assert.Nil(t, err, "success")
	accountHelper.Delete(account)
	accountHelper.Delete(account1)
}

func TestSelectAll(t *testing.T) {
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

func TestSelectAllByTag(t *testing.T) {
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

func TestSelectByConditions(t *testing.T) {
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

func TestSelectByConditionsAsMap(t *testing.T) {
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
	assert.Equal(t, "lamtnb@gmail.com", result[0]["email"] , "success")
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
	result, err := accountHelper.Create(account)
	assert.Nil(t, err, "err must be nil")
	affected, err := result.RowsAffected()
	assert.Nil(t, err, "err must be nil")
	assert.Equal(t, int(affected), 1, "add one row success")
	conditions := make(map[string]interface{})
	conditions["email"] = "baolam0307@gmail.com"
	result, err = accountHelper.DeleteByConditions(conditions)
	assert.Nil(t, err, "err must be nil")
	affected, err = result.RowsAffected()
	assert.Nil(t, err, "err  must be nil")
	assert.Equal(t, int(affected), 1, "add one row success")
}

func TestMain(m *testing.M) {
	setup()
	r := m.Run()
	//destroy()
	os.Exit(r)
}
