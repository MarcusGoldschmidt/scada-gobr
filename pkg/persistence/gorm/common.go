package gorm

import (
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"gorm.io/gorm"
)

func getById[T any](db *gorm.DB, id shared.CommonId) (*T, error) {
	var data *T
	result := db.Model(&data).Where("id = ?", id.String()).First(&data)

	if result.Error != nil {
		return nil, result.Error
	}

	return data, nil
}

func listAll[T any](db *gorm.DB) ([]*T, error) {
	var data []*T
	result := db.Model(&data).Find(&data)

	if result.Error != nil {
		return nil, result.Error
	}

	return data, nil
}

func listPaginate[T any](db *gorm.DB, paginate *shared.PaginationRequest) ([]*T, error) {
	var data []*T
	result := db.Model(&data).Scopes(paginate.GormScope).Find(&data)

	if result.Error != nil {
		return nil, result.Error
	}

	return data, nil
}
