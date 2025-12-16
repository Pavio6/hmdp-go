package service

import (
	"context"

	"gorm.io/gorm"

	"hmdp-backend/internal/model"
)

// SeckillVoucherService mirrors ISeckillVoucherService.
type SeckillVoucherService interface {
	Create(ctx context.Context, voucher *model.SeckillVoucher) error
}

type seckillVoucherService struct {
	db *gorm.DB
}

func (s *seckillVoucherService) Create(ctx context.Context, voucher *model.SeckillVoucher) error {
	return s.db.WithContext(ctx).Create(voucher).Error
}
