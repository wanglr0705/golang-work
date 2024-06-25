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

// 删除商品信息
func (l *Router) HandlerDeleteItem(c *gin.Context) {
	// 从URL参数中获取itemId
	itemIdStr := c.Param("itemId")
	itemId, err := strconv.Atoi(itemIdStr)

	// 从请求头中获取app_local
	appLocal := c.GetHeader("app_local")

	// 从上下文中获取errorstrCh
	errorstrChAny, _ := c.Get("errorstrCh")
	errorstrCh := errorstrChAny.(chan string)

	//itemId转换失败
	if err != nil {
		utils.ResponseError(c, types.ErrInvalidItemID, utils.LogError(errorstrCh, errors.New("无效的itemId")))
		return
	}

	//删除商品信息
	times, code, err := dao.DeleteItemDao(errorstrCh, l.Db, l.Rdb, l.Cache, l.DistributedLock, itemId, appLocal)
	if err != nil {
		utils.ResponseError(c, code, utils.LogError(errorstrCh, errors.New("删除商品信息失败")))
	}

	// 将删除时间格式化为字符串
	timeStr := times.Format("2006-01-02 15:04:05")
	utils.ResponseSuccess(c, pojo.DeleteItemResp{DeleteTime: timeStr})
}
