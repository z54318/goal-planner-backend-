package response

import "github.com/gin-gonic/gin"

// Body 是项目统一的接口返回结构。
type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorBody 用于 OpenAPI 文档中的失败响应结构。
type ErrorBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Success 返回成功响应。
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Body{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Fail 返回失败响应。
func Fail(c *gin.Context, code int, message string) {
	c.JSON(code, Body{
		Code:    code,
		Message: message,
	})
}
