package redis_distributed_lock

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
)

func (l *DistributedLock) Unlock(key string, value string) error {

	// 定义Lua脚本
	luaScript := `
	if redis.call("hget", KEYS[1],ARGV[1]) == ARGV[2] then
		redis.call("hdel",KEYS[1],ARGV[1])
		return 1
	else
		return 0
	end
	`

	// 执行Lua脚本
	reply, err := l.Rdb.Do("eval", luaScript, 1, "item", key, value)
	if err != nil {
		log.Println("解锁失败:", err)
		return err
	}

	//关闭看门狗
	l.WatchdogCancelFunc()

	// 将Redis的回复转换为字符串
	v, err := redis.Int64(reply, err)
	if err != nil {
		log.Println("Redis的回复转换为字符串失败:", err)
		return err
	}
	if v == 0 {
		return errors.New(fmt.Sprintf("解锁失败，key:%s,value:%s\n", key, value))
	} else {
		fmt.Printf("解锁成功，key:%s,value:%s\n", key, value)
	}

	return nil
}
