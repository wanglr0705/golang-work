package router

import (
	"errors"

	"go_xorm_mysql_redis/dao"
	"go_xorm_mysql_redis/pojo"
	"go_xorm_mysql_redis/types"
	"go_xorm_mysql_redis/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 删除商品信息
func (l *Router) HandlerDeleteItem(c *gin.Context) {
	// 从URL参数中获取itemId
	itemIdStr := c.Param("itemId")
	itemId, err := strconv.Atoi(itemIdStr)

	// 从请求头中获取app_local
	appLocal := c.GetHeader("app_local")

	//itemId转换失败
	if err != nil {
		panic(pojo.PanicData{Code: types.ErrInvalidItemID, Error: errors.New("无效的itemId")})
	}

	//删除商品信息
	times, code, err := dao.DeleteItemDao(l.Db, l.Rdb, l.Cache, l.DistributedLock, itemId, appLocal)
	if err != nil {
		utils.ResponseError(c, code, errors.New("删除商品信息失败"))
	}

	// 将删除时间格式化为字符串
	timeStr := times.Format("2006-01-02 15:04:05")
	utils.ResponseSuccess(c, code, pojo.DeleteItemResp{DeleteTime: timeStr})
}
