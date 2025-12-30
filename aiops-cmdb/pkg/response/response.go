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
}

type PageData struct {
	List       interface{} `json:"list"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

func buildResponse(code int, message string, data interface{}) Response {
	return Response{
		Code:      code,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, buildResponse(0, "success", data))
}

func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, buildResponse(0, message, data))
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, buildResponse(code, message, nil))
}

func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, buildResponse(http.StatusBadRequest, message, nil))
}

func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, buildResponse(http.StatusUnauthorized, message, nil))
}

func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, buildResponse(http.StatusForbidden, message, nil))
}

func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, buildResponse(http.StatusNotFound, message, nil))
}

func InternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, buildResponse(http.StatusInternalServerError, message, nil))
}

func PageSuccess(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	data := PageData{
		List:       list,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	Success(c, data)
}
