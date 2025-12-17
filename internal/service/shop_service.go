package service

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"hmdp-backend/internal/model"
)

// ShopService 处理商铺相关业务逻辑
type ShopService struct {
	db *gorm.DB
}

// NewShopService 创建 ShopService 实例
func NewShopService(db *gorm.DB) *ShopService {
	return &ShopService{db: db}
}

func (s *ShopService) GetByID(ctx context.Context, id int64) (*model.Shop, error) {
	var shop model.Shop
	err := s.db.WithContext(ctx).First(&shop, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &shop, nil
}

func (s *ShopService) Create(ctx context.Context, shop *model.Shop) error {
	return s.db.WithContext(ctx).Create(shop).Error
}

func (s *ShopService) Update(ctx context.Context, shop *model.Shop) error {
	return s.db.WithContext(ctx).Save(shop).Error
}

func (s *ShopService) QueryByType(ctx context.Context, typeID int64, page, size int) ([]model.Shop, error) {
	var shops []model.Shop
	offset := (page - 1) * size
	if offset < 0 {
		offset = 0
	}
	err := s.db.WithContext(ctx).
		Where("type_id = ?", typeID).
		Offset(offset).
		Limit(size).
		Order("id ASC").
		Find(&shops).Error
	return shops, err
}

func (s *ShopService) QueryByName(ctx context.Context, name string, page, size int) ([]model.Shop, error) {
	var shops []model.Shop
	offset := (page - 1) * size
	if offset < 0 {
		offset = 0
	}
	query := s.db.WithContext(ctx)
	if name != "" {
		query = query.Where("name LIKE ?", "%%"+name+"%%")
	}
	err := query.Order("id ASC").Offset(offset).Limit(size).Find(&shops).Error
	return shops, err
}
