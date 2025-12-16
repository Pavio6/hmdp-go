package model

import "time"

// VoucherOrder mirrors tb_voucher_order.
type VoucherOrder struct {
	ID         int64      `gorm:"column:id;primaryKey" json:"id"`
	UserID     int64      `gorm:"column:user_id" json:"userId"`
	VoucherID  int64      `gorm:"column:voucher_id" json:"voucherId"`
	PayType    int        `gorm:"column:pay_type" json:"payType"`
	Status     int        `gorm:"column:status" json:"status"`
	CreateTime time.Time  `gorm:"column:create_time" json:"createTime"`
	PayTime    *time.Time `gorm:"column:pay_time" json:"payTime"`
	UseTime    *time.Time `gorm:"column:use_time" json:"useTime"`
	RefundTime *time.Time `gorm:"column:refund_time" json:"refundTime"`
	UpdateTime time.Time  `gorm:"column:update_time" json:"updateTime"`
}

func (VoucherOrder) TableName() string { return "tb_voucher_order" }
