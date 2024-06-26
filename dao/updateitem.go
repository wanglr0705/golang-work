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

// 修改商品信息
func UpdateItemDao(db *xorm.Engine, rdb redis.Conn, cache *freecache.Cache, distributedLock *redis_distributed_lock.DistributedLock, updateItemReq pojo.UpdateItemReq) (pojo.ResponseData, int, error) {
	// 尝试获取分布式锁(防止脏读)
	value, err := distributedLock.Lock(types.LockKey)
	if err != nil {
		panic(pojo.PanicData{Code: types.ErrLock, Error: errors.New("加锁失败")})
	}
	defer func() {
		// 尝试释放分布式锁
		err = distributedLock.Unlock(types.LockKey, value)
		if err != nil {
			panic(pojo.PanicData{Code: types.ErrUnlock, Error: errors.New("解锁失败")})
		}
	}()

	//先更新mysql数据
	item := pojo.Item{
		ItemID:    updateItemReq.ItemID,
		Name:      updateItemReq.Name,
		Price:     updateItemReq.Price,
		UpdatedAt: time.Now(),
	}
	//设置更新指定字段数据
	_, err = db.Where("id = ?", item.ItemID).Cols("name", "price", "updated_at").Update(&item)
	if err != nil {
		panic(pojo.PanicData{Code: types.ErrMySQLUpdateData, Error: errors.New("更新mysql数据失败")})
	}
	_, err = db.Where("id = ?", item.ItemID).Get(&item)
	if err != nil {
		panic(pojo.PanicData{Code: types.ErrMySQLGetData, Error: errors.New("获取更新的数据失败")})
	}

	//后更新redis数据
	data, err := json.Marshal(&item)
	if err != nil {
		panic(pojo.PanicData{Code: types.ErrJSONConversion, Error: errors.New("json数据类型转换失败")})
	}
	key := types.GetItemKey(item.ItemID)
	_, err = rdb.Do("SET", key, string(data))
	if err != nil {
		panic(pojo.PanicData{Code: types.ErrRedisSetData, Error: errors.New("更新缓存数据失败")})
	}
	// 设置Redis缓存的过期时间为3600秒（1小时）
	_, err = rdb.Do("EXPIRE", key, 3600)
	if err != nil {
		panic(pojo.PanicData{Code: types.ErrRedisExpireTime, Error: errors.New("设置过期时间失败")})
	}

	//更新本地缓存
	id := strconv.Itoa(updateItemReq.ItemID)
	err = cache.Set([]byte(id), data, types.LocalCacheExpirationTime)
	if err != nil {
		return pojo.ResponseData{}, types.ErrCacheUpdateData, err
	}

	return pojo.ResponseData{
		ItemID: updateItemReq.ItemID,
		Name:   updateItemReq.Name,
		Price:  updateItemReq.Price,
	}, types.Success, nil
}
