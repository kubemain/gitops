package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type QueryParams struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Keyword  string `form:"keyword"`
	OrderBy  string `form:"order_by"`
	Order    string `form:"order"`
}

func (q *QueryParams) GetPage() int {
	if q.Page <= 0 {
		return 1
	}
	return q.Page
}

func (q *QueryParams) GetPageSize() int {
	if q.PageSize <= 0 {
		return 10
	}
	if q.PageSize > 100 {
		return 100
	}
	return q.PageSize
}

func (q *QueryParams) GetOffset() int {
	return (q.GetPage() - 1) * q.GetPageSize()
}

func (q *QueryParams) GetOrder() string {
	if q.OrderBy == "" {
		return "id desc"
	}
	if q.Order == "" || (q.Order != "asc" && q.Order != "desc") {
		q.Order = "desc"
	}
	return q.OrderBy + " " + q.Order
}
