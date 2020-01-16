package postgres

import (
	"fmt"
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

type BaseBo struct {
	Id string `json:"id" sql:"id"`
}

type UserBo struct {
	BaseBo
	Email    string `json:"email" sql:"email"`
	FullName string `json:"full_name" sql:"full_name"`
	Role     string `json:"role" sql:"role"`
	Active   bool   `json:"active" sql:"active"`
}

func TestGetDatabase(t *testing.T) {
	database := go_dal.GetDatabase()
	err := database.Ping()
	assert.Equal(t, err, nil, "Ping Success")
}

func TestInitHelper(t *testing.T) {
	helper := Helper{
		TableName: "account",
	}
	helper.Init()
}

func setup() {
	fmt.Println("run before")
	config := go_dal.Config{
		DriverName:     "postgres",
		DataSourceName: "postgres://lamtnb:Abc123@localhost:5432/template?sslmode=disable&client_encoding=UTF-8",
		MaxIdleConns:   5,
		MaxOpenConns:   5,
		MaxLifeTime:    1 * time.Minute}
	config.Init()
}

func TestMain(m *testing.M) {
	setup()
	r := m.Run()
	os.Exit(r)
}
