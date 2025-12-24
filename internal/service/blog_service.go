package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"hmdp-backend/internal/model"
	"hmdp-backend/internal/utils"
)

// BlogService 处理博客相关业务逻辑
type BlogService struct {
	db  *gorm.DB
	rdb *redis.Client
}

// NewBlogService 创建 BlogService 实例
func NewBlogService(db *gorm.DB, rdb *redis.Client) *BlogService {
	return &BlogService{db: db, rdb: rdb}
}

func (s *BlogService) Create(ctx context.Context, blog *model.Blog) error {
	return s.db.WithContext(ctx).Create(blog).Error
}

func (s *BlogService) GetByID(ctx context.Context, id int64) (*model.Blog, error) {
	var blog model.Blog
	err := s.db.WithContext(ctx).First(&blog, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &blog, nil
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

// ToggleLike 点赞/取消点赞；返回 true 表示点赞后状态
func (s *BlogService) ToggleLike(ctx context.Context, blogID, userID int64) (bool, error) {
	key := fmt.Sprintf("%s%d", utils.BLOG_LIKED_KEY, blogID)
	// 判断当前用户是否已点赞（使用 ZSet）
	_, err := s.rdb.ZScore(ctx, key, fmt.Sprint(userID)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return false, err
	}
	if errors.Is(err, redis.Nil) {
		if err := s.db.WithContext(ctx).
			Model(&model.Blog{}).
			Where("id = ?", blogID).
			UpdateColumn("liked", gorm.Expr("liked + 1")).Error; err != nil {
			return false, err
		}
		if err := s.rdb.ZAdd(ctx, key, redis.Z{
			Score:  float64(time.Now().Unix()),
			Member: fmt.Sprint(userID),
		}).Err(); err != nil {
			return false, err
		}
		return true, nil
	}

	// 已点赞，执行取消点赞
	if err := s.db.WithContext(ctx).
		Model(&model.Blog{}).
		Where("id = ? AND liked > 0", blogID).
		UpdateColumn("liked", gorm.Expr("liked - 1")).Error; err != nil {
		return false, err
	}
	if err := s.rdb.ZRem(ctx, key, fmt.Sprint(userID)).Err(); err != nil {
		return false, err
	}
	return false, nil
}

// IsLiked 判断用户是否点赞过
func (s *BlogService) IsLiked(ctx context.Context, blogID, userID int64) (bool, error) {
	key := fmt.Sprintf("%s%d", utils.BLOG_LIKED_KEY, blogID)
	_, err := s.rdb.ZScore(ctx, key, fmt.Sprint(userID)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// TopLikerIDs 返回最早点赞的前 N 个用户ID
func (s *BlogService) TopLikerIDs(ctx context.Context, blogID int64, limit int64) ([]int64, error) {
	key := fmt.Sprintf("%s%d", utils.BLOG_LIKED_KEY, blogID)
	members, err := s.rdb.ZRange(ctx, key, 0, limit-1).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	var ids []int64
	for _, m := range members {
		if v, err := strconv.ParseInt(m, 10, 64); err == nil {
			ids = append(ids, v)
		}
	}
	return ids, nil
}
