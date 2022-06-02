package shared

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type PaginationRequest struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

func NewPaginationRequest(r *http.Request) (*PaginationRequest, error) {
	queryStrings := r.URL.Query()

	page := queryStrings.Get("page")
	if page == "" {
		page = "1"
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return nil, errors.New("page is not a number")
	}

	size := queryStrings.Get("size")
	if size == "" {
		size = "20"
	}

	sizeInt, err := strconv.Atoi(size)
	if err != nil {
		return nil, errors.New("size is not a number")
	}

	return &PaginationRequest{
		Page: pageInt,
		Size: sizeInt,
	}, nil
}

func (p PaginationRequest) Offset() int {
	return (p.Page - 1) * p.Size
}

func (p *PaginationRequest) ToSql() map[string]any {
	return map[string]any{
		"offset":   p.Offset(),
		"pageSize": p.Page,
	}
}

func (p *PaginationRequest) GormScope(db *gorm.DB) *gorm.DB {
	return db.Offset(p.Offset()).Limit(p.Size)
}
