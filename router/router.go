package router

import (
	"context"
	"github.com/coocood/freecache"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"go_xorm_mysql_redis/config"
	"go_xorm_mysql_redis/middleware"
	redis_distributed_lock "go_xorm_mysql_redis/utils/redis-distributed-lock"
	"log"
	"xorm.io/xorm"
)

type Router struct {
	Db              *xorm.Engine
	Cfg             *config.Config
	Rdb             redis.Conn
	Cache           *freecache.Cache
	DistributedLock *redis_distributed_lock.DistributedLock
}

func NewRouter(ctx context.Context, cfg *config.Config, db *xorm.Engine, rdb redis.Conn, cache *freecache.Cache) *Router {
	distributedLock := redis_distributed_lock.NewDistributedLock(ctx, rdb, 30, 5)
	return &Router{
		Db:              db,
		Cfg:             cfg,
		Rdb:             rdb,
		Cache:           cache,
		DistributedLock: distributedLock,
	}
}

type RouterInter interface {
	RouterRegistration()
	HandlerAddItem(c *gin.Context)
	HandlerUpdateItem(c *gin.Context)
	HandlerGetItem(c *gin.Context)
	HandlerDeleteItem(c *gin.Context)
}

func (l *Router) RouterRegistration() {
	// 初始化路由。
	r := gin.Default()

	// 使用自定义的日志中间件
	r.Use(middleware.LoggerMiddleware)

	// 创建一个使用 ReqHeaderCheck 中间件的路由组。
	group := r.Group("/", middleware.ReqHeaderCheck)
	{
		//增加商品信息
		group.POST("/item", l.HandlerAddItem)

		//修改商品信息
		group.PUT("/item/:itemId", l.HandlerUpdateItem)

		//查询商品信息
		group.GET("/item/:itemId", l.HandlerGetItem)

		//删除商品信息
		group.DELETE("/item/:itemId", l.HandlerDeleteItem)
	}

	// 在端口 8080 上启动服务器。
	hostPort := l.Cfg.Host + ":" + l.Cfg.Port

	err := r.Run(hostPort)
	if err != nil {
		log.Println(err)
	}
}
