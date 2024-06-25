package local_caching

import "github.com/coocood/freecache"

var Cache *freecache.Cache

func LocalCaching() {
	// 创建一个内存大小为 100 MB 的缓存实例
	cacheSize := 100 * 1024 * 1024
	Cache = freecache.NewCache(cacheSize)
}
