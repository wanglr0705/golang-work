package router

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go_xorm_mysql_redis/dao"
	"go_xorm_mysql_redis/pojo"
	"go_xorm_mysql_redis/types"
	"go_xorm_mysql_redis/utils"
	"strconv"
)

// 获取商品信息
func (l *Router) HandlerGetItem(c *gin.Context) {
	// 从URL参数中获取itemId
	itemId := c.Param("itemId")
	itemIdInt, err := strconv.Atoi(itemId)

	// 从上下文中获取errorstrCh
	errorstrChAny, _ := c.Get("errorstrCh")
	errorstrCh := errorstrChAny.(chan string)
	if err != nil {
		err2 := utils.LogError(errorstrCh, errors.New("无效的itemId"))
		utils.ResponseError(c, types.ErrInvalidItemID, err2)
		return
	}

	// 调用DAO层获取商品信息
	responseData, code, err := dao.GetItemDao(errorstrCh, l.Db, l.Rdb, l.Cache, l.DistributedLock, itemIdInt)
	if err != nil {
		utils.ResponseError(c, code, utils.LogError(errorstrCh, err))
	} else {
		utils.ResponseSuccess(c, pojo.GetItemResp{StoreInfo: responseData})
	}
}
