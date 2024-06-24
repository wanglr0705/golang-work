package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var (
	Engine    *xorm.Engine
	userName  string = "root"
	passWord  string = "123456"
	idAddress string = "127.0.0.1"
	port      string = "3306"
	dbName    string = "go-xorm-mysql-redis"
	charset   string = "utf8"
)

func Mysql() error {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", userName, passWord, idAddress, port, dbName, charset)
	engine, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		return err
	}
	Engine = engine
	Engine.SetMaxIdleConns(10)
	Engine.SetMaxOpenConns(100)
	return err
}
