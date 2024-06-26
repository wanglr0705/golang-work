package dao

import (
	"encoding/json"
	"errors"
	"github.com/coocood/freecache"
	"github.com/gomodule/redigo/redis"
	"go_xorm_mysql_redis/pojo"
	"go_xorm_mysql_redis/types"
	redis_distributed_lock "go_xorm_mysql_redis/utils/redis-distributed-lock"
	"log"
	"strconv"
	"xorm.io/xorm"
)

func GetItemDao(db *xorm.Engine, rdb redis.Conn, cache *freecache.Cache, distributedLock *redis_distributed_lock.DistributedLock, itemId int) (pojo.ResponseData, int, error) {
	// 尝试获取分布式锁(强一致性，防止脏读)
	value, err := distributedLock.Lock(types.LockKey)
	if err != nil {
		panic(pojo.PanicData{Code: types.ErrLock, Error: err})
	}
	defer func() {
		// 尝试释放分布式锁
		err = distributedLock.Unlock(types.LockKey, value)
		if err != nil {
			panic(pojo.PanicData{Code: types.ErrUnlock, Error: err})
		}
	}()

	//从本地缓存获取数据
	itemIdStr := strconv.Itoa(itemId)
	v, err := cache.Get([]byte(itemIdStr))
	if err != nil {
		if err == freecache.ErrNotFound { //没有命中缓存
		} else {
			return pojo.ResponseData{}, types.ErrCacheGetData, nil
		}
	} else {
		// 缓存命中
		var item pojo.Item
		err := json.Unmarshal(v, &item)
		if err != nil {
			panic(pojo.PanicData{Code: types.ErrJSONConversion, Error: err})
		}
		return pojo.ResponseData{
			ItemID: item.ItemID,
			Name:   item.Name,
			Price:  item.Price,
		}, types.Success, nil
	}

	//先从缓存redis获取
	key := types.GetItemKey(itemId)
	reply, err := rdb.Do("GET", key)
	if err != nil {
		panic(pojo.PanicData{Code: types.ErrRedisGetData, Error: err})
	}
	//如果缓存中存在内容，那么直接返回缓存中的内容
	if reply != nil {
		var item pojo.Item
		str := string(reply.([]uint8))
		err := json.Unmarshal([]byte(str), &item)
		if err != nil {
			//继续从mysql中读取数据
			log.Println("数据解析失败，err:", err)
		} else {
			return pojo.ResponseData{
				ItemID: item.ItemID,
				Name:   item.Name,
				Price:  item.Price,
			}, types.Success, nil
		}
	}

	//如果缓存没有再从mysql获取，再写入缓存
	var item pojo.Item
	ok, err := db.Where("id = ? AND is_active = ?", itemId, 1).Get(&item)
	if err != nil {
		panic(pojo.PanicData{Code: types.ErrMySQLGetData, Error: err})
	}
	if ok {
		storeInfo := pojo.ResponseData{
			ItemID: item.ItemID,
			Name:   item.Name,
			Price:  item.Price,
		}
		//增加（覆盖）缓存数据
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
		err = cache.Set([]byte(itemIdStr), data, types.LocalCacheExpirationTime)
		if err != nil {
			panic(pojo.PanicData{Code: types.ErrCacheSetData, Error: errors.New("本地缓存写入失败")})
		}
		return storeInfo, types.Success, nil
	} else {
		return pojo.ResponseData{}, types.ErrMySQLDataNotFound, errors.New("没有找到该商品")
	}
}
