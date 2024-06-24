package redis_distributed_lock

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

// 看门狗
func (l *DistributedLock) watchdog(key string) error {
	luaScript := `
	if redis.call("ttl", KEYS[1]) <= 5 then
		return redis.call('expire',KEYS[1],ARGV[1])
	else
		return 0
	end
	`
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		select {
		case <-l.WatchdogCtx.Done(): //结束看门狗
			fmt.Printf("关闭看门狗，key:%s\n", key)
			return nil
		default:
			//判断是否需要增加锁的过期时间
			reply, err := l.Rdb.Do("EVAL", luaScript, 1, key, l.WatchdogTime)
			if err != nil {
				log.Println("看门狗执行错误：", err)
				break
			}
			v, err := redis.Int64(reply, err)
			if err != nil {
				log.Println("Redis的回复转换为字符串行错误：", err)
			}
			fmt.Println("看门狗执行结果：", v)
			break
		}
	}
	return nil
}
