package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`
	TraceID   string      `json:"trace_id,omitempty"`
}

func buildResponse(c *gin.Context, code int, message string, data interface{}) Response {
	return Response{
		Code:      code,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Unix(),
		TraceID:   c.GetString("trace_id"),
	}
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, buildResponse(c, 0, "success", data))
}

func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, buildResponse(c, 0, message, data))
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, buildResponse(c, code, message, nil))
}

func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, buildResponse(c, http.StatusUnauthorized, message, nil))
}

func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, buildResponse(c, http.StatusForbidden, message, nil))
}

func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, buildResponse(c, http.StatusNotFound, message, nil))
}

func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, buildResponse(c, http.StatusBadRequest, message, nil))
}

func InternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, buildResponse(c, http.StatusInternalServerError, message, nil))
}

func ServiceUnavailable(c *gin.Context, message string) {
	c.JSON(http.StatusServiceUnavailable, buildResponse(c, http.StatusServiceUnavailable, message, nil))
}

func GatewayTimeout(c *gin.Context, message string) {
	c.JSON(http.StatusGatewayTimeout, buildResponse(c, http.StatusGatewayTimeout, message, nil))
}
