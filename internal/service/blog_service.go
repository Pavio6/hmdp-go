package service

import (
	"context"

	"gorm.io/gorm"

	"hmdp-backend/internal/model"
)

// BlogService 处理博客相关业务逻辑
type BlogService struct {
	db *gorm.DB
}

// NewBlogService 创建 BlogService 实例
func NewBlogService(db *gorm.DB) *BlogService {
	return &BlogService{db: db}
}

func (s *BlogService) Create(ctx context.Context, blog *model.Blog) error {
	return s.db.WithContext(ctx).Create(blog).Error
}

func (s *BlogService) IncrementLike(ctx context.Context, id int64) error {
	return s.db.WithContext(ctx).
		Model(&model.Blog{}).
		Where("id = ?", id).
		UpdateColumn("liked", gorm.Expr("liked + 1")).
		Error
}

func (s *BlogService) QueryByUser(ctx context.Context, userID int64, page, size int) ([]model.Blog, error) {
	var blogs []model.Blog
	offset := (page - 1) * size
	if offset < 0 {
		offset = 0
	}
	err := s.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("id ASC").
		Offset(offset).
		Limit(size).
		Find(&blogs).Error
	return blogs, err
}

func (s *BlogService) QueryHot(ctx context.Context, page, size int) ([]model.Blog, error) {
	var blogs []model.Blog
	offset := (page - 1) * size
	if offset < 0 {
		offset = 0
	}
	err := s.db.WithContext(ctx).
		Order("liked DESC").
		Offset(offset).
		Limit(size).
		Find(&blogs).Error
	return blogs, err
}
