package service

import "gorm.io/gorm"

// Registry aggregates all service implementations for easy wiring.
type Registry struct {
	Blog           BlogService
	Shop           ShopService
	ShopType       ShopTypeService
	Voucher        VoucherService
	SeckillVoucher SeckillVoucherService
	User           UserService
	UserInfo       UserInfoService
}

// NewRegistry wires all services using the shared DB.
func NewRegistry(db *gorm.DB) *Registry {
	seckillSvc := &seckillVoucherService{db: db}
	return &Registry{
		Blog:           &blogService{db: db},
		Shop:           &shopService{db: db},
		ShopType:       &shopTypeService{db: db},
		Voucher:        &voucherService{db: db, seckillSvc: seckillSvc},
		SeckillVoucher: seckillSvc,
		User:           &userService{db: db},
		UserInfo:       &userInfoService{db: db},
	}
}
