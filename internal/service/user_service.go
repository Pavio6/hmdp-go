package service

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"hmdp-backend/internal/model"
)

// UserService mirrors IUserService usage.
type UserService interface {
	FindByID(ctx context.Context, id int64) (*model.User, error)
}

type userService struct {
	db *gorm.DB
}

func (s *userService) FindByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	err := s.db.WithContext(ctx).First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
