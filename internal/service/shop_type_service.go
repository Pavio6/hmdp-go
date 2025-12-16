package service

import (
	"context"

	"gorm.io/gorm"

	"hmdp-backend/internal/model"
)

// ShopTypeService mirrors IShopTypeService.
type ShopTypeService interface {
	List(ctx context.Context) ([]model.ShopType, error)
}

type shopTypeService struct {
	db *gorm.DB
}

func (s *shopTypeService) List(ctx context.Context) ([]model.ShopType, error) {
	var types []model.ShopType
	err := s.db.WithContext(ctx).
		Order("sort ASC").
		Find(&types).Error
	return types, err
}
