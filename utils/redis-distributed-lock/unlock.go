package redis_distributed_lock

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"go_xorm_mysql_redis/pojo"
	"go_xorm_mysql_redis/types"
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
		//log.Println("解锁失败:", utils.LogError(ch, err))
		//return err
		err = fmt.Errorf("解锁失败: %v", err)
		panic(pojo.PanicData{Error: err, Code: types.ErrUnlock})
	}

	//关闭看门狗
	l.WatchdogCancelFunc()

	// 将Redis的回复转换为字符串
	v, err := redis.Int(reply, err)
	if err != nil {
		//log.Println("Redis的回复转换为整数失败:", utils.LogError(ch, err))
		//return err
		err = fmt.Errorf("Redis的回复转换为整数失败: %v", err)
		panic(pojo.PanicData{Error: err, Code: types.ErrJSONConversion})
	}
	if v == 0 {
		//return errors.New(fmt.Sprintf("解锁失败，key:%s,value:%s\n", key, value))
		err = fmt.Errorf("解锁失败，key:%s,value:%s", key, value)
		panic(pojo.PanicData{Error: err, Code: types.ErrUnlock})
	} else {
		fmt.Printf("解锁成功，key:%s,value:%s\n", key, value)
	}

	return nil
}
