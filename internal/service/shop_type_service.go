package service

import (
	"context"

	"gorm.io/gorm"

	"hmdp-backend/internal/model"
)

// ShopTypeService 处理店铺类型数据
type ShopTypeService struct {
	db *gorm.DB
}

// NewShopTypeService 创建 ShopTypeService 实例
func NewShopTypeService(db *gorm.DB) *ShopTypeService {
	return &ShopTypeService{db: db}
}

func (s *ShopTypeService) List(ctx context.Context) ([]model.ShopType, error) {
	var types []model.ShopType
	err := s.db.WithContext(ctx).
		Order("sort ASC").
		Find(&types).Error
	return types, err
}
