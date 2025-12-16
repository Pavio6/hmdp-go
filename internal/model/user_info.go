package model

import "time"

// UserInfo mirrors tb_user_info.
type UserInfo struct {
	ID         int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID     int64      `gorm:"column:user_id" json:"userId"`
	Gender     string     `gorm:"column:gender" json:"gender"`
	City       string     `gorm:"column:city" json:"city"`
	Birthday   *time.Time `gorm:"column:birthday" json:"birthday"`
	Introduce  string     `gorm:"column:introduce" json:"introduce"`
	CreateTime *time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateTime *time.Time `gorm:"column:update_time" json:"updateTime"`
}

func (UserInfo) TableName() string { return "tb_user_info" }
