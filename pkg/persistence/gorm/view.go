package gorm

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ViewPersistenceGormImpl struct {
	db *gorm.DB
}

func (v ViewPersistenceGormImpl) GetViewComponentById(ctx context.Context, u uuid.UUID) (*models.ViewComponent, error) {
	return getById[models.ViewComponent](v.db.WithContext(ctx), u)
}

func NewViewPersistenceGormImpl(db *gorm.DB) *ViewPersistenceGormImpl {
	return &ViewPersistenceGormImpl{db: db}
}

func (v ViewPersistenceGormImpl) AttachViewComponents(ctx context.Context, components ...*models.ViewComponent) error {
	db := v.db.WithContext(ctx)

	for _, component := range components {
		if err := db.Model(&models.ViewComponent{}).Where("id = ?", component.Id.String()).Save(component).Error; err != nil {

			if err == gorm.ErrRecordNotFound {
				db.Create(&component)
			}
		}
	}

	return nil
}

func (v ViewPersistenceGormImpl) GetViewById(ctx context.Context, uuid uuid.UUID) (*models.View, error) {
	db := v.db.WithContext(ctx)

	var data *models.View
	err := db.Model(&models.View{Id: uuid}).Preload("ViewComponents").First(&data).Error

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (v ViewPersistenceGormImpl) GetAllViews(ctx context.Context) ([]*models.View, error) {
	db := v.db.WithContext(ctx)
	return listAll[models.View](db)
}

func (v ViewPersistenceGormImpl) CreateView(ctx context.Context, view *models.View) error {
	db := v.db.WithContext(ctx)
	return db.Create(view).Error
}

func (v ViewPersistenceGormImpl) UpdateView(ctx context.Context, view *models.View) error {
	db := v.db.WithContext(ctx)
	return db.Save(view).Error
}

func (v ViewPersistenceGormImpl) DeleteView(ctx context.Context, uuid uuid.UUID) error {
	db := v.db.WithContext(ctx)
	return db.Delete(&models.View{Id: uuid}).Error
}

func (v ViewPersistenceGormImpl) DeleteViewComponent(ctx context.Context, viewId uuid.UUID, componentId uuid.UUID) error {
	db := v.db.WithContext(ctx)
	return db.Delete(&models.ViewComponent{Id: componentId, ViewId: viewId}).Error
}
