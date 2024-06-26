package dao

import (
	"errors"
	"github.com/coocood/freecache"
	"github.com/gomodule/redigo/redis"
	"go_xorm_mysql_redis/pojo"
	"go_xorm_mysql_redis/types"
	"go_xorm_mysql_redis/utils"
	redis_distributed_lock "go_xorm_mysql_redis/utils/redis-distributed-lock"
	"strconv"
	"time"
	"xorm.io/xorm"
)

// 删除商品信息
func DeleteItemDao(db *xorm.Engine, rdb redis.Conn, cache *freecache.Cache, distributedLock *redis_distributed_lock.DistributedLock, itemId int, appLocal string) (time.Time, int, error) {
	// 尝试获取分布式锁(强一致性，同时防止脏读)
	value, err := distributedLock.Lock(types.LockKey)
	if err != nil {
		panic(pojo.PanicData{Code: types.ErrLock, Error: errors.New("加锁失败")})
	}
	defer func() {
		// 尝试释放分布式锁
		err = distributedLock.Unlock(types.LockKey, value)
		if err != nil {
			panic(pojo.PanicData{Code: types.ErrUnlock, Error: errors.New("释放锁失败")})
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
		panic(pojo.PanicData{Code: types.ErrMySQLDeleteData, Error: errors.New("删除数据失败")})
	}

	//再删除缓存
	key := types.GetItemKey(itemId)
	_, err = rdb.Do("DEL", key)
	if err != nil {
		panic(pojo.PanicData{Code: types.ErrRedisSetData, Error: errors.New("删除缓存失败")})
	}

	//删除本地缓存
	itemIdStr := strconv.Itoa(itemId)
	_ = cache.Del([]byte(itemIdStr))

	return utils.GetTime(appLocal), types.Success, nil
}
