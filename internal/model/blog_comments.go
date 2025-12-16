package model

import "time"

// BlogComments mirrors tb_blog_comments.
type BlogComments struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	BlogID      int64     `gorm:"column:blog_id" json:"blogId"`
	UserID      int64     `gorm:"column:user_id" json:"userId"`
	ParentID    *int64    `gorm:"column:parent_id" json:"parentId"`
	ReplyUserID *int64    `gorm:"column:answer_id" json:"answerId"`
	Content     string    `gorm:"column:content" json:"content"`
	Liked       int       `gorm:"column:liked" json:"liked"`
	Status      int       `gorm:"column:status" json:"status"`
	CreateTime  time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateTime  time.Time `gorm:"column:update_time" json:"updateTime"`
}

func (BlogComments) TableName() string { return "tb_blog_comments" }
