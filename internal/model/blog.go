package model

import "time"

// Blog mirrors tb_blog.
type Blog struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ShopID     int64     `gorm:"column:shop_id" json:"shopId"`
	UserID     int64     `gorm:"column:user_id" json:"userId"`
	Title      string    `gorm:"column:title" json:"title"`
	Images     string    `gorm:"column:images" json:"images"`
	Content    string    `gorm:"column:content" json:"content"`
	Liked      int       `gorm:"column:liked" json:"liked"`
	Comments   int       `gorm:"column:comments" json:"comments"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	Icon       string    `gorm:"-" json:"icon,omitempty"`
	Name       string    `gorm:"-" json:"name,omitempty"`
	IsLike     *bool     `gorm:"-" json:"isLike,omitempty"`
}

func (Blog) TableName() string { return "tb_blog" }
