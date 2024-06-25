package main

import (
	"context"
	"go_xorm_mysql_redis/config"
	local_caching "go_xorm_mysql_redis/local-caching"
	"go_xorm_mysql_redis/mysql"
	"go_xorm_mysql_redis/redis"
	"go_xorm_mysql_redis/router"
	"go_xorm_mysql_redis/utils"
)

func main() {
	cfg := config.Config{}
	// 从YAML文件中获取配置信息
	utils.GetYaml(&cfg)

	//连接mysql
	err := mysql.Mysql(&cfg)
	if err != nil {
		panic(err)
	}

	//连接redis
	err = redis.Redis(&cfg)
	if err != nil {
		panic(err)
	}

	//创建本地缓存
	local_caching.LocalCaching()

	ctx := context.Background()
	//注册路由
	newRouter := router.NewRouter(ctx, &cfg, mysql.Engine, redis.RDB, local_caching.Cache)
	newRouter.RouterRegistration()
}
