package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware(c *gin.Context) {
	// @warn fmt.Print处理日志不太好
	fmt.Println("============请求开始============")
	// 记录请求开始时间
	start := time.Now()

	//@warn 这个是做什么的
	// 创建用于传递错误信息的通道
	errorstrCh := make(chan string)
	errstrSliceCh := make(chan []string)
	c.Set("errorstrCh", errorstrCh)
	defer func() {
		close(errorstrCh)
		// 从错误信息切片通道中获取错误信息切片
		errstrSlice := <-errstrSliceCh
		close(errstrSliceCh)
		// 遍历并打印错误信息
		fmt.Println("错误：")
		for _, errstr := range errstrSlice {
			fmt.Println(errstr)
		}
	}()

	// 启动一个goroutine来收集错误信息
	go func(errorstrCh2 chan string, errstrSliceCh2 chan []string) {
		// 初始化一个错误信息切片
		var errorSlice []string
		for true {
			select {
			case errstr, ok := <-errorstrCh2:
				errorSlice = append(errorSlice, errstr)
				// 如果通道关闭，将错误信息切片发送到通道并返回
				if !ok {
					errstrSliceCh2 <- errorSlice
					//fmt.Println("若多次打印，则内存泄露")
					return
				}
				break
			default:
				break
			}
		}
	}(errorstrCh, errstrSliceCh)

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
