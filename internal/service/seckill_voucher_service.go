package service

import (
	"context"

	"gorm.io/gorm"

	"hmdp-backend/internal/model"
)

// SeckillVoucherService 处理秒杀券相关业务
type SeckillVoucherService struct {
	db *gorm.DB
}

// NewSeckillVoucherService 创建 SeckillVoucherService 实例
func NewSeckillVoucherService(db *gorm.DB) *SeckillVoucherService {
	return &SeckillVoucherService{db: db}
}

func (s *SeckillVoucherService) Create(ctx context.Context, voucher *model.SeckillVoucher) error {
	return s.db.WithContext(ctx).Create(voucher).Error
}
