package dao

import (
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	"go_xorm_mysql_redis/pojo"
	"go_xorm_mysql_redis/types"
	"go_xorm_mysql_redis/utils"
	redis_distributed_lock "go_xorm_mysql_redis/utils/redis-distributed-lock"
	"log"
	"xorm.io/xorm"
)

func GetItemDao(ch chan string, db *xorm.Engine, rdb redis.Conn, distributedLock *redis_distributed_lock.DistributedLock, itemId int) (pojo.ResponseData, int, error) {
	// 尝试获取分布式锁(强一致性，防止脏读)
	value, err := distributedLock.Lock(ch, types.LockKey)
	if err != nil {
		err2 := utils.LogError(ch, errors.New("加锁失败"))
		return pojo.ResponseData{}, types.LockInvalidRequest, errors.Join(err, err2)
	}
	defer func() {
		// 尝试释放分布式锁
		err = distributedLock.Unlock(ch, types.LockKey, value)
		if err != nil {
			panic(utils.LogError(ch, errors.New("解锁失败")))
		}
	}()
	//先从缓存redis获取
	key := types.GetItemKey(itemId)
	reply, err := rdb.Do("GET", key)
	if err != nil {
		return pojo.ResponseData{}, types.ErrRedisGetData, utils.LogError(ch, err)
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
		return pojo.ResponseData{}, types.ErrMySQLGetData, utils.LogError(ch, err)
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
			return storeInfo, types.ErrJSONConversion, utils.LogError(ch, errors.New("json转换错误"))
		}
		// 将商品信息存入Redis缓存
		key := types.GetItemKey(item.ItemID)
		_, err = rdb.Do("SET", key, string(data))
		if err != nil {
			return storeInfo, types.ErrRedisSetData, utils.LogError(ch, errors.New("缓存数据失败"))
		}
		// 设置Redis缓存的过期时间为3600秒（1小时）
		_, err = rdb.Do("EXPIRE", key, 3600)
		if err != nil {
			return storeInfo, types.ErrRedisExpireTime, utils.LogError(ch, errors.New("设置过期时间失败"))
		}
		return storeInfo, types.Success, nil
	} else {
		return pojo.ResponseData{}, types.ErrMySQLDataNotFound, utils.LogError(ch, errors.New("没有找到该商品"))
	}
}
