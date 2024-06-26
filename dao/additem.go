package dao

import (
	"encoding/json"
	"errors"
	"github.com/coocood/freecache"
	"github.com/gomodule/redigo/redis"
	"go_xorm_mysql_redis/pojo"
	"go_xorm_mysql_redis/types"
	redis_distributed_lock "go_xorm_mysql_redis/utils/redis-distributed-lock"
	"strconv"
	"time"
	"xorm.io/xorm"
)

// 添加商品信息
func AddItemDao(db *xorm.Engine, rdb redis.Conn, cache *freecache.Cache, distributedLock *redis_distributed_lock.DistributedLock, addItemReq pojo.AddItemReq) (pojo.ResponseData, int, error) {
	// 尝试获取分布式锁(防止重复插入)
	value, err := distributedLock.Lock(types.LockKey)
	if err != nil {
		panic(pojo.PanicData{Code: types.ErrLock, Error: errors.New("加锁失败")})
	}
	defer func() {
		// 尝试释放分布式锁
		err = distributedLock.Unlock(types.LockKey, value)
		// 如果解锁失败，将错误发送到通道
		if err != nil {
			panic(pojo.PanicData{Code: types.ErrUnlock, Error: err})
		}
	}()

	//增加mysql数据
	now := time.Now()
	item := pojo.Item{
		Name:      addItemReq.Name,
		Price:     addItemReq.Price,
		IsActive:  1,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: now,
	}
	_, err = db.Insert(&item)
	if err != nil {
		panic(pojo.PanicData{Code: types.ErrMySQLSetData, Error: err})
	}

	//增加缓存数据
	data, err := json.Marshal(&item)
	if err != nil {
		panic(pojo.PanicData{Code: types.ErrJSONConversion, Error: errors.New("json转换错误")})
	}
	// 将商品信息存入Redis缓存
	key := types.GetItemKey(item.ItemID)
	_, err = rdb.Do("SET", key, string(data))
	if err != nil {
		panic(pojo.PanicData{Code: types.ErrRedisSetData, Error: errors.New("缓存数据失败")})
	}
	// 设置Redis缓存的过期时间为3600秒（1小时）
	_, err = rdb.Do("EXPIRE", key, 3600)
	if err != nil {
		panic(pojo.PanicData{Code: types.ErrRedisExpireTime, Error: errors.New("设置过期时间失败")})
	}

	//写入本地缓存
	itemId := strconv.Itoa(item.ItemID)
	err = cache.Set([]byte(itemId), data, types.LocalCacheExpirationTime)
	if err != nil {
		panic(pojo.PanicData{Code: types.ErrCacheSetData, Error: errors.New("本地缓存写入失败")})
	}

	return pojo.ResponseData{
		ItemID: item.ItemID,
		Name:   item.Name,
		Price:  item.Price,
	}, types.Success, nil
}
