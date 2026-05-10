package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ApiResponse 统一响应结构，匹配前端 axios 拦截器的 code !== 200 判断
type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Success 成功响应 (code=200)
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, ApiResponse{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// Error 业务错误响应（自定义 code）
func Error(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, ApiResponse{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// BadRequest 参数错误 (400)
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, 400, message)
}

// Unauthorized 认证失败 (401)
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, 401, message)
}

// Forbidden 权限不足 (403)
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, 403, message)
}

// NotFound 资源不存在 (404)
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, 404, message)
}

// TooManyRequests 限流 (429)
func TooManyRequests(c *gin.Context, message string) {
	Error(c, http.StatusTooManyRequests, 429, message)
}

// InternalError 服务器内部错误 (500)
func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, 500, message)
}

// PaginatedData 分页数据，匹配前端 PaginatedData<T>
type PaginatedData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}
