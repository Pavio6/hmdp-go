package service

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"hmdp-backend/internal/model"
)

// UserInfoService mirrors IUserInfoService usage.
type UserInfoService interface {
	FindByID(ctx context.Context, id int64) (*model.UserInfo, error)
}

type userInfoService struct {
	db *gorm.DB
}

func (s *userInfoService) FindByID(ctx context.Context, id int64) (*model.UserInfo, error) {
	var info model.UserInfo
	err := s.db.WithContext(ctx).First(&info, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &info, nil
}
