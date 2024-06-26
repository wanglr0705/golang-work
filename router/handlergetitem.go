package router

import (
	"go_xorm_mysql_redis/dao"
	"go_xorm_mysql_redis/pojo"
	"go_xorm_mysql_redis/types"
	"go_xorm_mysql_redis/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取商品信息
func (l *Router) HandlerGetItem(c *gin.Context) {

	// 从URL参数中获取itemId
	itemId := c.Param("itemId")
	itemIdInt, err := strconv.Atoi(itemId)

	if err != nil {
		panic(pojo.PanicData{types.ErrInvalidItemID, err})
	}

	// 调用DAO层获取商品信息
	responseData, code, err := dao.GetItemDao(l.Db, l.Rdb, l.Cache, l.DistributedLock, itemIdInt)
	if err != nil {
		utils.ResponseError(c, code, err)
	} else {
		utils.ResponseSuccess(c, code, pojo.GetItemResp{StoreInfo: responseData})
	}
}
