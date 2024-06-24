package main

import (
	"context"
	"go_xorm_mysql_redis/mysql"
	"go_xorm_mysql_redis/redis"
	"go_xorm_mysql_redis/router"
)

func main() {
	ctx := context.Background()
	newRouter := router.NewRouter(ctx, mysql.Engine, redis.RDB)
	newRouter.RouterRegistration()
}

// 初始化
func init() {
	//连接mysql
	err := mysql.Mysql()
	if err != nil {
		panic(err)
	}

	//连接redis
	err = redis.Redis()
	if err != nil {
		panic(err)
	}
}
