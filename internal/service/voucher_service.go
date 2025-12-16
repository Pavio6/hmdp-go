package service

import (
	"context"
	"time"

	"gorm.io/gorm"

	"hmdp-backend/internal/model"
)

// VoucherService mirrors IVoucherService.
type VoucherService interface {
	Create(ctx context.Context, voucher *model.Voucher) error
	QueryVoucherOfShop(ctx context.Context, shopID int64) ([]model.Voucher, error)
	AddSeckillVoucher(ctx context.Context, voucher *model.Voucher) error
}

type voucherService struct {
	db         *gorm.DB
	seckillSvc SeckillVoucherService
}

func (s *voucherService) Create(ctx context.Context, voucher *model.Voucher) error {
	return s.db.WithContext(ctx).Create(voucher).Error
}

func (s *voucherService) QueryVoucherOfShop(ctx context.Context, shopID int64) ([]model.Voucher, error) {
	var vouchers []model.Voucher
	query := `
        SELECT v.id, v.shop_id, v.title, v.sub_title, v.rules, v.pay_value,
               v.actual_value, v.type, v.status, v.create_time, v.update_time,
               sv.stock, sv.begin_time, sv.end_time
        FROM tb_voucher v
        LEFT JOIN tb_seckill_voucher sv ON v.id = sv.voucher_id
        WHERE v.shop_id = ? AND v.status = 1`
	err := s.db.WithContext(ctx).Raw(query, shopID).Scan(&vouchers).Error
	return vouchers, err
}

func (s *voucherService) AddSeckillVoucher(ctx context.Context, voucher *model.Voucher) error {
	if err := s.Create(ctx, voucher); err != nil {
		return err
	}
	stock := 0
	if voucher.Stock != nil {
		stock = *voucher.Stock
	}
	var begin time.Time
	if voucher.BeginTime != nil {
		begin = *voucher.BeginTime
	}
	var end time.Time
	if voucher.EndTime != nil {
		end = *voucher.EndTime
	}
	sec := &model.SeckillVoucher{
		VoucherID: voucher.ID,
		Stock:     stock,
		BeginTime: begin,
		EndTime:   end,
	}
	return s.seckillSvc.Create(ctx, sec)
}
