package model

import "time"

// Shop mirrors tb_shop entity.
type Shop struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"column:name" json:"name"`
	TypeID     int64     `gorm:"column:type_id" json:"typeId"`
	Images     string    `gorm:"column:images" json:"images"`
	Area       string    `gorm:"column:area" json:"area"`
	Address    string    `gorm:"column:address" json:"address"`
	X          float64   `gorm:"column:x" json:"x"`
	Y          float64   `gorm:"column:y" json:"y"`
	AvgPrice   int64     `gorm:"column:avg_price" json:"avgPrice"`
	Sold       int       `gorm:"column:sold" json:"sold"`
	Comments   int       `gorm:"column:comments" json:"comments"`
	Score      int       `gorm:"column:score" json:"score"`
	OpenHours  string    `gorm:"column:open_hours" json:"openHours"`
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"`
	Distance   *float64  `gorm:"-" json:"distance,omitempty"`
}

func (Shop) TableName() string { return "tb_shop" }
