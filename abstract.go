package go_dal

import (
	"database/sql"
	"log"
	"time"
)

/**
 * @author LTNB (baolam0307@gmail.com)
 * @since
 *
 */

type IDatabaseHelper interface {
	GetDatabase() *sql.DB
}

var db *sql.DB

type Config struct {
	DriverName     string
	DataSourceName string
	MaxOpenConns   int
	MaxLifeTime    time.Duration
	MaxIdleConns   int
}

func (dbConf Config) Init() {
	initDefaultValue(&dbConf)
	var err error
	conn, err := sql.Open(dbConf.DriverName, dbConf.DataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	conn.SetMaxOpenConns(dbConf.MaxOpenConns)
	conn.SetConnMaxLifetime(dbConf.MaxLifeTime)
	conn.SetMaxIdleConns(dbConf.MaxIdleConns)
	err = conn.Ping()
	if err != nil {
		log.Fatal(err)
	}
	db = conn
}

func (dbConf Config) InitDefaultConfig(driverName string){
	config := Config{DriverName: driverName}
	initDefaultValue(&config)
	dbConf.Init()
}



func initDefaultValue(database *Config) {
	if database.DataSourceName == "" && database.DriverName == "postgres" {
		database.DataSourceName = "postgres://postgres:postgres@localhost:5432/template?sslmode=disable&client_encoding=UTF-8&stringtype=unspecified"
	}
	if database.MaxOpenConns == 0 {
		database.MaxOpenConns = 4
	}
	if database.MaxLifeTime == 0 {
		database.MaxLifeTime = 30 * time.Second
	}
	if database.MaxIdleConns == 0 {
		database.MaxIdleConns = 4
	}
}

func GetDatabase() *sql.DB{
	return db
}
