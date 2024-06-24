package dao

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"go_xorm_mysql_redis/pojo"
	"go_xorm_mysql_redis/types"
	"go_xorm_mysql_redis/utils"
	redis_distributed_lock "go_xorm_mysql_redis/utils/redis-distributed-lock"
	"time"
	"xorm.io/xorm"
)

// 删除商品信息
func DeleteItemDao(ch chan string, db *xorm.Engine, rdb redis.Conn, distributedLock *redis_distributed_lock.DistributedLock, itemId int) (time.Time, int, error) {
	// 尝试获取分布式锁(强一致性，同时防止脏读)
	value, err := distributedLock.Lock(ch, types.LockKey)
	if err != nil {
		return time.Time{}, types.LockInvalidRequest, errors.Join(err, utils.LogError(ch, errors.New("加锁失败")))
	}
	defer func() {
		// 尝试释放分布式锁
		err = distributedLock.Unlock(ch, types.LockKey, value)
		if err != nil {
			panic(utils.LogError(ch, err))
		}
	}()

	//先删除数据库(为实现幂等性，使用软删除)
	now := time.Now()
	item := pojo.Item{
		ItemID:    itemId,
		IsActive:  0,
		DeletedAt: now,
	}
	_, err = db.Where("id = ?", itemId).Cols("is_active", "deleted_at").Update(&item)
	if err != nil {
		return time.Time{}, types.ErrMySQLDeleteData, errors.Join(err, utils.LogError(ch, errors.New("删除数据失败")))
	}

	//再删除缓存
	key := types.GetItemKey(itemId)
	_, err = rdb.Do("DEL", key)
	if err != nil {
		return time.Time{}, types.ErrRedisSetData, errors.Join(err, utils.LogError(ch, errors.New("删除缓存失败")))
	}

	return now, types.Success, nil
}
