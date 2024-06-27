package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_xorm_mysql_redis/item"
	"go_xorm_mysql_redis/utils"
	"runtime/debug"
	"time"
)

func LoggerMiddleware(c *gin.Context) {
	// 记录请求开始时间
	start := time.Now()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("严重错误panic信息:", r)
			panicData := r.(pojo.PanicData)
			utils.ResponseError(c, panicData.Code, panicData.Error)
			debug.PrintStack()
		}
	}()

	// 处理请求
	c.Next()

	// 记录请求结束时间，计算请求处理时间
	end := time.Now()
	latency := end.Sub(start)

	// 遍历请求头
	fmt.Println("请求头：")
	for name, values := range c.Request.Header {
		// 遍历每个头字段的值
		for _, value := range values {
			fmt.Printf("%s: %s\n", name, value)
		}
	}

	// 获取请求信息
	method := c.Request.Method
	path := c.Request.URL.Path
	clientIP := c.ClientIP()
	statusCode := c.Writer.Status()

	// 打印日志
	fmt.Printf("[GIN] %v | %3d | %13v | %15s | %-7s %s\n",
		end.Format("2006/01/02 - 15:04:05"),
		statusCode,
		latency,
		clientIP,
		method,
		path,
	)

}
