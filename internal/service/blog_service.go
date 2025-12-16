package service

import (
	"context"

	"gorm.io/gorm"

	"hmdp-backend/internal/model"
)

// BlogService mirrors IBlogService usage.
type BlogService interface {
	Create(ctx context.Context, blog *model.Blog) error
	IncrementLike(ctx context.Context, id int64) error
	QueryByUser(ctx context.Context, userID int64, page, size int) ([]model.Blog, error)
	QueryHot(ctx context.Context, page, size int) ([]model.Blog, error)
}

type blogService struct {
	db *gorm.DB
}

func (s *blogService) Create(ctx context.Context, blog *model.Blog) error {
	return s.db.WithContext(ctx).Create(blog).Error
}

func (s *blogService) IncrementLike(ctx context.Context, id int64) error {
	return s.db.WithContext(ctx).
		Model(&model.Blog{}).
		Where("id = ?", id).
		UpdateColumn("liked", gorm.Expr("liked + 1")).
		Error
}

func (s *blogService) QueryByUser(ctx context.Context, userID int64, page, size int) ([]model.Blog, error) {
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

func (s *blogService) QueryHot(ctx context.Context, page, size int) ([]model.Blog, error) {
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
