package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
)

type MySql struct {
	db *gorm.DB
}

var Engine = new(MySql)

func init()  {

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/test?charset=utf8&parseTime=True&loc=Local",
		"root",
		"123456",
		"127.0.0.1",
		"3306",
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             0,
			Colorful:                  true,
			IgnoreRecordNotFoundError: false,
			LogLevel:                  logger.Info,
		}),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	},)

	if err != nil {
		panic(err)
	}

	Engine.db = db
}

func (m *MySql) Driver() *gorm.DB {
	return Engine.db
}