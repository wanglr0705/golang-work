package types

var (
	//加锁失败
	LockInvalidRequest = 0001
	//解锁失败
	LockInvalidUnlock = 0002

	// Redis errors
	ErrRedisGetData    = 1001 // 获取缓存数据错误
	ErrRedisSetData    = 1002 // 写入Redis缓存错误
	ErrRedisExpireTime = 1003 // 设置Redis过期时间错误

	// MySQL errors
	ErrMySQLGetData      = 2001 // 从MySQL获取数据错误
	ErrMySQLDataNotFound = 2002 // 找不到数据错误
	ErrMySQLSetData      = 2003 // 写入MySQL错误
	ErrMySQLUpdateData   = 2004 // 更新MySQL数据错误
	ErrMySQLDeleteData   = 2005 // 删除MySQL数据错误

	// General errors
	ErrJSONConversion = 3001 // JSON转换错误

	//request
	ErrInvalidItemID            = 4001 // 无效的itemId
	ErrMissingRequiredParameter = 4002 // 缺少必填参数
	ErrMissingAppLocalHeader    = 4005 // 缺少app_local请求头字段

)
