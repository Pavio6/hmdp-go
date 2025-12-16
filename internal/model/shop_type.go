package model

import "time"

// ShopType mirrors tb_shop_type.
type ShopType struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"column:name" json:"name"`
	Icon       string    `gorm:"column:icon" json:"icon"`
	Sort       int       `gorm:"column:sort" json:"sort"`
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"`
}

func (ShopType) TableName() string { return "tb_shop_type" }
