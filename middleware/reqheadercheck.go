package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go_xorm_mysql_redis/pojo"
	"go_xorm_mysql_redis/types"
)

// 请求头校验中间件
func ReqHeaderCheck(c *gin.Context) {
	// 从请求头中获取 "app_local" 字段
	appLocal := c.GetHeader("app_local")

	// 根据 "app_local" 字段的值进行不同的处理
	switch appLocal {
	case "uk": //英国
		//fmt.Println("英国站点：uk")
		break
	case "jp": //日本
		//fmt.Println("日本站点：jp")
		break
	case "ur": //俄罗斯
		//fmt.Println("俄罗斯站点：ur")
		break
	default:
		// 如果 "app_local" 字段不匹配任何已知值，返回错误响应
		c.Abort() // 终止请求处理
		panic(pojo.PanicData{Code: types.ErrMissingAppLocalHeader, Error: errors.New("请正确设置请求头中app_local字段")})
	}

	// 继续处理请求
	c.Next()
}
