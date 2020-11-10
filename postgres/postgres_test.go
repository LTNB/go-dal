package postgres

import (
	goDal "github.com/LTNB/go-dal"
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

type BaseBo struct {
	Id string `json:"id" primary:"id" sql:"id"`
}

type Auditor struct {
	Date time.Time `json:"date" sql:"date"`
}
type AccountMock struct {
	BaseBo   `promoted:"true"`
	Auditor  `promoted:"true"`
	Email    string `json:"email" sql:"email"`
	FullName string `json:"full_name" sql:"full_name"`
	Role     string `json:"role" sql:"role"`
	Active   bool   `json:"active" sql:"active"`
}

func setup() {
	conf := goDal.Config{
		DriverName:     "postgres",
		DataSourceName: "postgres://postgres:123456@localhost:5432/template?sslmode=disable&client_encoding=UTF-8",
		MaxOpenConns:   5,
		MaxLifeTime:    1 * time.Minute,
		MaxIdleConns:   5,
	}

	conf.Init()
	aHelper := sql.Helper{
		TableName:      "account",
		Bo:             AccountMock{},
		DefaultTagName: "json",
		Handler:        Helper{},
	}
	accountHelper = aHelper
}

func TestConnection(t *testing.T) {
	db := goDal.GetDatabase()
	err := db.Ping()
	assert.Nil(t, err, "connected")

}

func TestGetOne(t *testing.T) {
	account := AccountMock{
		BaseBo:   BaseBo{Id: "1"},
		Auditor:  Auditor{Date: time.Now()},
		Email:    "baolam0307@gmail.com",
		FullName: "Ta Ngoc Bao Lam",
		Role:     "admin",
		Active:   true,
	}
	_, err := accountHelper.Create(account)
	assert.Nil(t, err, "nothing")
	bo := AccountMock{}
	bo.Id = "1"
	err = accountHelper.GetOne(&bo)
	assert.Nil(t, err, "nothing")
	assert.Equal(t, "baolam0307@gmail.com", bo.Email, "success")
	_, err = accountHelper.Delete(bo)
	assert.Nil(t, err, "nothing")

}

//
func TestGetOneByTag(t *testing.T) {
	account := AccountMock{
		BaseBo:   BaseBo{Id: "1"},
		Auditor:  Auditor{Date: time.Now()},
		Email:    "baolam0307@gmail.com",
		FullName: "Ta Ngoc Bao Lam",
		Role:     "admin",
		Active:   true,
	}
	_, err := accountHelper.Create(account)
	assert.Nil(t, err, "nothing")
	bo := AccountMock{}
	bo.Id = "1"
	err = accountHelper.GetOneByTag(&bo, "json")
	assert.Nil(t, err, "nothing")
	assert.Equal(t, "baolam0307@gmail.com", bo.Email, "success")
	_, err = accountHelper.Delete(bo)
	assert.Nil(t, err, "nothing")
}

func TestGetOneByConditions(t *testing.T) {
	account := AccountMock{
		BaseBo:   BaseBo{Id: "1"},
		Auditor:  Auditor{Date: time.Now()},
		Email:    "baolam0307@gmail.com",
		FullName: "Ta Ngoc Bao Lam",
		Role:     "admin",
		Active:   true,
	}
	_, err := accountHelper.Create(account)
	assert.Nil(t, err, "nothing")
	bo := AccountMock{}
	conditions := make(map[string]interface{})
	conditions["role"] = "admin"
	err = accountHelper.GetOneByConditions(&bo, conditions, "json")
	assert.Nil(t, err, "nothing")
	assert.Equal(t, "baolam0307@gmail.com", bo.Email, "success")
	_, err = accountHelper.Delete(bo)
	assert.Nil(t, err, "nothing")
}

func TestGetAsMap(t *testing.T) {
	account := AccountMock{
		BaseBo:   BaseBo{Id: "1"},
		Auditor:  Auditor{Date: time.Now()},
		Email:    "baolam0307@gmail.com",
		FullName: "Ta Ngoc Bao Lam",
		Role:     "admin",
		Active:   true,
	}
	_, err := accountHelper.Create(account)
	assert.Nil(t, err, "nothing")
	bo := AccountMock{}
	bo.Id = "1"
	result, err := accountHelper.GetOneAsMap(&bo)
	assert.Equal(t, "baolam0307@gmail.com", result["email"], "success")
	assert.Nil(t, err, "success")
	_, err = accountHelper.Delete(bo)
	assert.Nil(t, err, "nothing")
}

func TestGetAll(t *testing.T) {
	account := AccountMock{
		BaseBo:   BaseBo{Id: "1"},
		Auditor:  Auditor{Date: time.Now()},
		Email:    "baolam0307@gmail.com",
		FullName: "Ta Ngoc Bao Lam",
		Role:     "admin",
		Active:   true,
	}
	account1 := AccountMock{
		BaseBo:   BaseBo{Id: "2"},
		Auditor:  Auditor{Date: time.Now()},
		Email:    "lamtnb@gmail.com",
		FullName: "lamtnb",
		Role:     "user",
		Active:   false,
	}
	_, err := accountHelper.Create(account)
	assert.Nil(t, err, "nothing")
	_, err =accountHelper.Create(account1)
	assert.Nil(t, err, "nothing")
	result, err := accountHelper.GetAll()
	assert.Nil(t, err, "nothing")

	assert.Equal(t, 2, len(result), "success")
	_, err = accountHelper.Delete(account)
	assert.Nil(t, err, "nothing")
	_, err = accountHelper.Delete(account1)
	assert.Nil(t, err, "nothing")
}

func TestGetAllByTag(t *testing.T) {
	account := AccountMock{
		BaseBo:   BaseBo{Id: "1"},
		Auditor:  Auditor{Date: time.Now()},
		Email:    "baolam0307@gmail.com",
		FullName: "Ta Ngoc Bao Lam",
		Role:     "admin",
		Active:   true,
	}
	account1 := AccountMock{
		BaseBo:   BaseBo{Id: "2"},
		Auditor:  Auditor{Date: time.Now()},
		Email:    "lamtnb@gmail.com",
		FullName: "lamtnb",
		Role:     "user",
		Active:   false,
	}
	_, err :=accountHelper.Create(account)
	assert.Nil(t, err, "nothing")
	_, err = accountHelper.Create(account1)
	assert.Nil(t, err, "nothing")
	result, err := accountHelper.GetAllByTag("json")

	assert.Equal(t, 2, len(result), "success")
	assert.Nil(t, err, "nothing")
	_, err = accountHelper.Delete(account)
	assert.Nil(t, err, "nothing")
	_, err = accountHelper.Delete(account1)
	assert.Nil(t, err, "nothing")
}

func TestGetAllAsMap(t *testing.T) {
	account := AccountMock{
		BaseBo:   BaseBo{Id: "1"},
		Auditor:  Auditor{Date: time.Now()},
		Email:    "baolam0307@gmail.com",
		FullName: "Ta Ngoc Bao Lam",
		Role:     "admin",
		Active:   true,
	}
	account1 := AccountMock{
		BaseBo:   BaseBo{Id: "2"},
		Auditor:  Auditor{Date: time.Now()},
		Email:    "lamtnb@gmail.com",
		FullName: "lamtnb",
		Role:     "user",
		Active:   false,
	}
	_, err := accountHelper.Create(account)
	assert.Nil(t, err, "nothing")
	_, err = accountHelper.Create(account1)
	assert.Nil(t, err, "nothing")
	result, err := accountHelper.GetAllAsMap()
	assert.Equal(t, "baolam0307@gmail.com", result[0]["email"], "success")
	assert.Equal(t, "lamtnb@gmail.com", result[1]["email"], "success")
	assert.Nil(t, err, "success")
	_, err = accountHelper.Delete(account)
	assert.Nil(t, err, "nothing")
	_, err = accountHelper.Delete(account1)
	assert.Nil(t, err, "nothing")
}

func TestGetByConditions(t *testing.T) {
	account := AccountMock{
		BaseBo:   BaseBo{Id: "1"},
		Auditor:  Auditor{Date: time.Now()},
		Email:    "baolam0307@gmail.com",
		FullName: "Ta Ngoc Bao Lam",
		Role:     "admin",
		Active:   true,
	}
	account1 := AccountMock{
		BaseBo:   BaseBo{Id: "2"},
		Auditor:  Auditor{Date: time.Now()},
		Email:    "lamtnb@gmail.com",
		FullName: "lamtnb",
		Role:     "user",
		Active:   false,
	}
	_, err := accountHelper.Create(account)
	assert.Nil(t, err, "nothing")
	_, err = accountHelper.Create(account1)
	assert.Nil(t, err, "nothing")

	conditions := make(map[string]interface{})
	conditions["role"] = "admin"
	orderBy := make(map[string]string)
	orderBy["full_name"] = "ASC"
	limit := 1
	offset := 0
	result, err := accountHelper.GetByConditions(conditions, orderBy, limit, offset, "")

	assert.Equal(t, 1, len(result), "success")
	assert.Nil(t, err, "nothing")
	_,err = accountHelper.Delete(account)
	_, err = accountHelper.Delete(account1)
	assert.Nil(t, err, "nothing")
}

func TestGetByConditionsAsMap(t *testing.T) {
	account := AccountMock{
		BaseBo:   BaseBo{Id: "1"},
		Auditor:  Auditor{Date: time.Now()},
		Email:    "baolam0307@gmail.com",
		FullName: "Ta Ngoc Bao Lam",
		Role:     "admin",
		Active:   true,
	}
	account1 := AccountMock{
		BaseBo:   BaseBo{Id: "2"},
		Auditor:  Auditor{Date: time.Now()},
		Email:    "lamtnb@gmail.com",
		FullName: "lamtnb",
		Role:     "user",
		Active:   false,
	}
	_, err := accountHelper.Create(account)
	_, err = accountHelper.Create(account1)

	orderBy := make(map[string]string)
	orderBy["full_name"] = "ASC"
	limit := 1
	offset := 0
	result, err := accountHelper.GetByConditionsAsMap(nil, orderBy, limit, offset, "")

	assert.Equal(t, 1, len(result), "success")
	assert.Equal(t, "baolam0307@gmail.com", result[0]["email"], "success")
	assert.Nil(t, err, "nothing")
	_, err = accountHelper.Delete(account)
	assert.Nil(t, err, "nothing")
	_, err = accountHelper.Delete(account1)
	assert.Nil(t, err, "nothing")
}

func TestCreateAndDelete(t *testing.T) {
	account := AccountMock{
		BaseBo:   BaseBo{Id: "1"},
		Auditor:  Auditor{Date: time.Now()},
		Email:    "baolam0307@gmail.com",
		FullName: "Ta Ngoc Bao Lam",
		Role:     "UnixDate",
		Active:   true,
	}
	_, err := accountHelper.Create(account)
	assert.Nil(t, err, "err must be nil")
	conditions := make(map[string]interface{})
	conditions["email"] = "baolam0307@gmail.com"
	_, err = accountHelper.DeleteByConditions(conditions)
	assert.Nil(t, err, "err must be nil")
}

/**
 * create batch object
 */
func TestCreateBatch(t *testing.T) {
	account := AccountMock{
		BaseBo:   BaseBo{Id: "1"},
		Auditor:  Auditor{Date: time.Now()},
		Email:    "baolam0307@gmail.com",
		FullName: "Ta Ngoc Bao Lam",
		Role:     "UnixDate",
		Active:   true,
	}
	account2 := AccountMock{
		BaseBo:   BaseBo{Id: "2"},
		Auditor:  Auditor{Date: time.Now()},
		Email:    "lamtnb@gmail.com",
		FullName: "LTNB",
		Role:     "UnixDate",
		Active:   true,
	}
	boList := []interface{}{account, account2}
	effected, err := accountHelper.CreateBatch(boList)
	assert.Nil(t, err, "err must be nil")
	assert.Equal(t, effected, int64(2))

	//delete data test
	conditions := make(map[string]interface{})
	conditions["active"] = true
	_, err = accountHelper.DeleteByConditions(conditions)
}

func TestCreateByTagAndDelete(t *testing.T) {
	account := AccountMock{
		BaseBo:   BaseBo{Id: "1"},
		Auditor:  Auditor{Date: time.Now()},
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

func TestUpdate(t *testing.T) {
	account := AccountMock{
		BaseBo:   BaseBo{Id: "1"},
		Auditor:  Auditor{Date: time.Now()},
		Email:    "baolam0307@gmail.com",
		FullName: "Ta Ngoc Bao Lam",
		Role:     "admin",
		Active:   true,
	}
	_, err := accountHelper.Create(account)
	assert.Nil(t, err, "nothing")
	account = AccountMock{
		BaseBo:   BaseBo{Id: "1"},
		Auditor:  Auditor{Date: time.Now()},
		Email:    "lamtnb@scommerce.asia",
		FullName: "Ta Ngoc Bao Lam",
		Role:     "admin",
		Active:   true,
	}
	_, err = accountHelper.Update(account)
	assert.Nil(t, err, "nothing")
	account1 := AccountMock{}
	account1.Id = "1"
	err = accountHelper.GetOne(&account1)
	assert.Nil(t, err, "nothing")
	assert.Equal(t, "lamtnb@scommerce.asia", account1.Email, "success")
	_, err = accountHelper.Delete(account1)
	assert.Nil(t, err, "nothing")
}

func TestUpdateByTag(t *testing.T) {
	account := AccountMock{
		BaseBo:   BaseBo{Id: "1"},
		Auditor:  Auditor{Date: time.Now()},
		Email:    "baolam0307@gmail.com",
		FullName: "Ta Ngoc Bao Lam",
		Role:     "admin",
		Active:   true,
	}
	_, err := accountHelper.Create(account)
	assert.Nil(t, err, "nothing")
	account = AccountMock{
		BaseBo:   BaseBo{Id: "1"},
		Auditor:  Auditor{Date: time.Now()},
		Email:    "lamtnb@scommerce.asia",
		FullName: "Ta Ngoc Bao Lam",
		Role:     "admin",
		Active:   true,
	}
	_, err = accountHelper.UpdateByTag(account, "sql")
	assert.Nil(t, err, "nothing")
	account1 := AccountMock{}
	account1.Id = "1"
	err = accountHelper.GetOne(&account1)
	assert.Nil(t, err, "nothing")
	assert.Equal(t, "lamtnb@scommerce.asia", account1.Email, "success")
	_, err = accountHelper.Delete(account1)
	assert.Nil(t, err, "nothing")
}

func TestMain(m *testing.M) {
	setup()
	r := m.Run()
	//destroy()
	os.Exit(r)
}
