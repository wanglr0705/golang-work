package redis_distributed_lock

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
)

func (l *DistributedLock) Lock(key string) (string, error) {
	// 生成一个唯一的锁令牌
	value := GetToken()

	// 定义Lua脚本
	luaScript := `
	local v = redis.call("hget",KEYS[1],ARGV[1])
	if not v then
		redis.call("hset",KEYS[1],ARGV[1],ARGV[2])
		redis.call("expire",KEYS[1],ARGV[3])
		return 1
	else
		return 0
	end
	`
	// 执行Lua脚本，尝试获取锁
	do, err := l.Rdb.Do("EVAL", luaScript, 1, "item", key, value, l.ExpireTime)
	if err != nil {
		fmt.Println("加锁失败:", err)
		return "", err
	}

	// 解析Redis执行Lua脚本的回复
	reply, err := redis.Int(do, err)
	if err != nil {
		fmt.Println("Redis的回复转换为字符串失败:", err)
		return "", err
	}
	if reply == 0 {
		return "", fmt.Errorf("锁已被占用，key:%s\n", key)
	} else {
		fmt.Printf("加锁成功，key:%s,value:%s\n", key, value)
	}

	//异步启动看门狗
	go func() {
		err2 := l.watchdog(key)
		if err2 != nil {
			log.Println("看门狗启动失败:", err2)
		}
	}()

	// 返回生成的锁令牌
	return value, nil
}
