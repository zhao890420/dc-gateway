package common

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DefaultDbPool *gorm.DB

func InitDatabase() error {
	mysqlUrl := config.MustValue("mysql", "mysql_url", "")
	vdb, err := gorm.Open(mysql.Open(mysqlUrl), &gorm.Config{})
	if err != nil {
		DefLogger.Error(fmt.Sprintf("init query with mysql url err: %v", err))
		return err
	}
	DefaultDbPool = vdb
	DefLogger.Info("======finish db init ")
	return nil
}

type Db struct {
}

func (c Db) Destroy() {
}
