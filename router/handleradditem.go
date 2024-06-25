package router

import (
	"go_xorm_mysql_redis/dao"
	"go_xorm_mysql_redis/pojo"
	"go_xorm_mysql_redis/types"
	"go_xorm_mysql_redis/utils"

	"github.com/gin-gonic/gin"
)

// 增加商品信息
func (l *Router) HandlerAddItem(c *gin.Context) {
	var req pojo.AddItemReq

	// 从上下文中获取errorstrCh
	errorstrChAny, _ := c.Get("errorstrCh")
	errorstrCh := errorstrChAny.(chan string)

	// 绑定JSON请求数据到req结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, types.ErrMissingRequiredParameter, utils.LogError(errorstrCh, err))
	}

	// 调用DAO层添加商品信息
	itemInfo, code, err := dao.AddItemDao(errorstrCh, l.Db, l.Rdb, l.Cache, l.DistributedLock, req)
	if err != nil {
		utils.ResponseError(c, code, utils.LogError(errorstrCh, err))
	} else {
		// 添加成功，返回成功响应
		utils.ResponseSuccess(c, pojo.AddItemResp{Item_info: itemInfo})
	}
}
