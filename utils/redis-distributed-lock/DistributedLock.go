package redis_distributed_lock

import (
	"context"
	"github.com/gomodule/redigo/redis"
)

type DistributedLock struct {
	Ctx                context.Context
	WatchdogCtx        context.Context
	WatchdogCancelFunc context.CancelFunc
	Rdb                redis.Conn
	ExpireTime         int //初始过期时间
	WatchdogTime       int //看门狗单次增加的时间
}

func NewDistributedLock(ctx context.Context, rdb redis.Conn, expireTime int, watchdogTime int) *DistributedLock {
	watchdogCtx, cancelFunc := context.WithCancel(ctx)
	return &DistributedLock{
		Ctx:                ctx,
		WatchdogCtx:        watchdogCtx,
		WatchdogCancelFunc: cancelFunc,
		Rdb:                rdb,
		ExpireTime:         expireTime,
		WatchdogTime:       watchdogTime,
	}
}

type DistributedLockIn interface {
	Lock(ch chan string, key string) (string, error)
	Unlock(ch chan string, key string, value string) error
	watchdog(ch chan string, key string) error
}
