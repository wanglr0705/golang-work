package utils

import (
	"github.com/gin-gonic/gin"
)

// 响应体格式
type Body struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// 错误返回
func ResponseError(c *gin.Context, code int, err error) {
	var body Body
	if err != nil {
		body.Code = code
		body.Msg = err.Error()
	}
	c.JSON(200, body)
}

// 正确返回
func ResponseSuccess(c *gin.Context, code int, resp interface{}) {
	var body Body
	body.Code = code
	body.Msg = "成功"
	body.Data = resp
	c.JSON(200, body)
}
