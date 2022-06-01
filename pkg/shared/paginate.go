package shared

import (
	"errors"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type PaginationRequest struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

func NewPaginationRequest(r *http.Request) (*PaginationRequest, error) {
	vars := mux.Vars(r)

	page, ok := vars["page"]
	if !ok {
		page = "1"
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return nil, errors.New("page is not a number")
	}

	size, ok := vars["size"]
	if !ok {
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
