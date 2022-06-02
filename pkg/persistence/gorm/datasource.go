package gorm

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"gorm.io/gorm"
)

type DataSourcePersistenceGormImpl struct {
	db *gorm.DB
}

func (d DataSourcePersistenceGormImpl) GetDadaSourceById(ctx context.Context, id shared.CommonId) (*models.DataSource, error) {
	db := d.db.WithContext(ctx)
	return getById[models.DataSource](db, id)
}

func (d DataSourcePersistenceGormImpl) GetDadaSources(ctx context.Context) ([]*models.DataSource, error) {
	db := d.db.WithContext(ctx)
	return listAll[models.DataSource](db)
}
