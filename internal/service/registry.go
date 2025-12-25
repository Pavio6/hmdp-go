package service

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Registry 聚合全部业务 Service，方便注入 handler
type Registry struct {
	Blog           *BlogService
	Shop           *ShopService
	ShopType       *ShopTypeService
	Voucher        *VoucherService
	SeckillVoucher *SeckillVoucherService
	User           *UserService
	VoucherOrder   *VoucherOrderService
	Follow         *FollowService
}

// NewRegistry 使用共享 DB 与 Redis 构建所有服务
func NewRegistry(db *gorm.DB, rdb *redis.Client, log *zap.Logger) *Registry {
	if log == nil {
		log = zap.NewNop()
	}
	seckillSvc := NewSeckillVoucherService(db)
	followSvc := NewFollowService(db, rdb)
	return &Registry{
		Blog:           NewBlogService(db, rdb, followSvc),
		Shop:           NewShopService(db, rdb, log),
		ShopType:       NewShopTypeService(db, rdb),
		Voucher:        NewVoucherService(db, seckillSvc, rdb),
		SeckillVoucher: seckillSvc,
		User:           NewUserService(db, rdb),
		VoucherOrder:   NewVoucherOrderService(db, rdb),
		Follow:         followSvc,
	}
}
