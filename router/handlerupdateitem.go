package router

import (
	"go_xorm_mysql_redis/dao"
	"go_xorm_mysql_redis/item"
	"go_xorm_mysql_redis/types"
	"go_xorm_mysql_redis/utils"

	"github.com/gin-gonic/gin"
)

// 修改商品信息
func (l *Router) HandlerUpdateItem(c *gin.Context) {
	var updateItemReq pojo.UpdateItemReq

	// 绑定JSON请求数据到updateItemReq结构体
	if err := c.ShouldBindJSON(&updateItemReq); err != nil {
		panic(pojo.PanicData{Code: types.ErrMissingRequiredParameter, Error: err})
	}

	// 调用DAO层更新商品信息
	storeInfo, code, err := dao.UpdateItemDao(l.Db, l.Rdb, l.Cache, l.DistributedLock, updateItemReq)
	if err != nil {
		utils.ResponseError(c, code, err)
	} else {
		// 更新成功，返回成功响应
		utils.ResponseSuccess(c, code, pojo.UpdateItemResp{StoreInfo: storeInfo})
	}
}
