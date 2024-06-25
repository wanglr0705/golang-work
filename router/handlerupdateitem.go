package router

import (
	"github.com/gin-gonic/gin"
	"go_xorm_mysql_redis/dao"
	"go_xorm_mysql_redis/pojo"
	"go_xorm_mysql_redis/types"
	"go_xorm_mysql_redis/utils"
)

// 修改商品信息
func (l *Router) HandlerUpdateItem(c *gin.Context) {
	var updateItemReq pojo.UpdateItemReq

	// 从上下文中获取errorstrCh
	errorstrChAny, _ := c.Get("errorstrCh")
	errorstrCh := errorstrChAny.(chan string)

	// 绑定JSON请求数据到updateItemReq结构体
	if err := c.ShouldBindJSON(&updateItemReq); err != nil {
		utils.ResponseError(c, types.ErrMissingRequiredParameter, utils.LogError(errorstrCh, err))
	}

	// 调用DAO层更新商品信息
	storeInfo, code, err := dao.UpdateItemDao(errorstrCh, l.Db, l.Rdb, l.Cache, l.DistributedLock, updateItemReq)
	if err != nil {
		utils.ResponseError(c, code, utils.LogError(errorstrCh, err))
	} else {
		// 更新成功，返回成功响应
		utils.ResponseSuccess(c, pojo.UpdateItemResp{StoreInfo: storeInfo})
	}
}
