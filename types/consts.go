package types

import "strconv"

var (
	//分布式锁的key
	LockKey = "item_lock_key"
	ItemKey = "item_data:"
)

// 这是一个将数据id和ItemKey结合一起来的函数
func GetItemKey(itemID int) string {
	return ItemKey + strconv.Itoa(itemID)
}
