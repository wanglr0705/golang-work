package dao

import (
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	"go_xorm_mysql_redis/pojo"
	"go_xorm_mysql_redis/types"
	"go_xorm_mysql_redis/utils"
	redis_distributed_lock "go_xorm_mysql_redis/utils/redis-distributed-lock"
	"time"
	"xorm.io/xorm"
)

// 修改商品信息
func UpdateItemDao(ch chan string, db *xorm.Engine, rdb redis.Conn, distributedLock *redis_distributed_lock.DistributedLock, updateItemReq pojo.UpdateItemReq) (pojo.ResponseData, int, error) {
	// 尝试获取分布式锁(防止脏读)
	value, err := distributedLock.Lock(ch, types.LockKey)
	if err != nil {
		return pojo.ResponseData{}, types.LockInvalidRequest, errors.Join(err, utils.LogError(ch, errors.New("加锁失败")))
	}
	defer func() {
		// 尝试释放分布式锁
		err = distributedLock.Unlock(ch, types.LockKey, value)
		if err != nil {
			panic(utils.LogError(ch, err))
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
		return pojo.ResponseData{}, types.ErrMySQLUpdateData, utils.LogError(ch, errors.New("更新mysql数据失败"))
	}
	_, err = db.Where("id = ?", item.ItemID).Get(&item)
	if err != nil {
		return pojo.ResponseData{}, types.ErrMySQLGetData, utils.LogError(ch, errors.New("获取更新的数据失败"))
	}

	//后更新redis数据
	data, err := json.Marshal(&item)
	if err != nil {
		return pojo.ResponseData{}, types.ErrJSONConversion, utils.LogError(ch, errors.New("json数据类型转换失败"))
	}
	key := types.GetItemKey(item.ItemID)
	_, err = rdb.Do("SET", key, string(data))
	if err != nil {
		return pojo.ResponseData{}, types.ErrRedisSetData, utils.LogError(ch, errors.New("更新缓存数据失败"))
	}
	// 设置Redis缓存的过期时间为3600秒（1小时）
	_, err = rdb.Do("EXPIRE", key, 3600)
	if err != nil {
		return pojo.ResponseData{}, types.ErrRedisExpireTime, utils.LogError(ch, errors.New("设置过期时间失败"))
	}

	return pojo.ResponseData{
		ItemID: updateItemReq.ItemID,
		Name:   updateItemReq.Name,
		Price:  updateItemReq.Price,
	}, types.Success, nil
}
