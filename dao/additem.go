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
	"time"
	"xorm.io/xorm"
)

// 添加商品信息
func AddItemDao(ch chan string, db *xorm.Engine, rdb redis.Conn, distributedLock *redis_distributed_lock.DistributedLock, addItemReq pojo.AddItemReq) (pojo.ResponseData, int, error) {
	// 尝试获取分布式锁(防止重复插入)
	value, err := distributedLock.Lock(ch, types.LockKey)
	if err != nil {
		err2 := utils.LogError(ch, errors.New("加锁失败"))
		return pojo.ResponseData{}, types.LockInvalidRequest, errors.Join(err, err2)
	}
	defer func() {
		// 尝试释放分布式锁
		err = distributedLock.Unlock(ch, types.LockKey, value)
		// 如果解锁失败，将错误发送到通道
		if err != nil {
			panic(utils.LogError(ch, err))
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
		log.Println("mysql插入数据失败:", utils.LogError(ch, err))
		return pojo.ResponseData{}, types.ErrMySQLSetData, err
	}

	//增加缓存数据
	data, err := json.Marshal(&item)
	if err != nil {
		return pojo.ResponseData{}, types.ErrJSONConversion, utils.LogError(ch, errors.New("json转换错误"))
	}
	// 将商品信息存入Redis缓存
	key := types.GetItemKey(item.ItemID)
	_, err = rdb.Do("SET", key, string(data))
	if err != nil {
		return pojo.ResponseData{}, types.ErrRedisSetData, utils.LogError(ch, errors.New("缓存数据失败"))
	}
	// 设置Redis缓存的过期时间为3600秒（1小时）
	_, err = rdb.Do("EXPIRE", key, 3600)
	if err != nil {
		return pojo.ResponseData{}, types.ErrRedisExpireTime, utils.LogError(ch, errors.New("设置过期时间失败"))
	}

	return pojo.ResponseData{
		ItemID: item.ItemID,
		Name:   item.Name,
		Price:  item.Price,
	}, types.Success, nil
}
