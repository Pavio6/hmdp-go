package service

import (
	"context"
	"errors"
	"hmdp-backend/internal/mapper"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"hmdp-backend/internal/dto"
	"hmdp-backend/internal/model"
	"hmdp-backend/internal/utils"
)

// UserService 处理登录与验证码相关业务
type UserService struct {
	db  *gorm.DB
	rdb *redis.Client
}

// NewUserService 创建 UserService 实例
func NewUserService(db *gorm.DB, rdb *redis.Client) *UserService {
	return &UserService{db: db, rdb: rdb}
}

func (s *UserService) SendCode(ctx context.Context, phone string) error {
	// 1.校验手机号
	if utils.IsPhoneInvalid(phone) {
		return errors.New("phone is invalid")
	}
	// 2.生成验证码
	code, err := utils.GenerateVerifyCode()
	if err != nil {
		return err
	}
	// 3.将验证码存到redis中
	key := utils.LOGIN_CODE_KEY + phone
	if err := s.rdb.Set(ctx, key, code, time.Duration(utils.LOGIN_CODE_TTL)*time.Minute).Err(); err != nil {
		return err
	}

	// 4.发送验证码
	log.Println("验证码为:", code)
	return nil
}

func (s *UserService) Login(ctx context.Context, loginForm dto.LoginForm) (string, error) {
	var user model.User
	// 1.校验手机号
	if utils.IsPhoneInvalid(loginForm.Phone) {
		return "", errors.New("phone is invalid")
	}
	// 2.校验验证码
	codeKey := utils.LOGIN_CODE_KEY + loginForm.Phone
	cacheCode, err := s.rdb.Get(ctx, codeKey).Result()
	if errors.Is(err, redis.Nil) {
		return "", errors.New("验证码不存在或已过期")
	}
	if err != nil {
		return "", err
	}
	if cacheCode != loginForm.Code {
		return "", errors.New("验证码错误")
	}
	// 验证通过后清理验证码，避免重复使用
	if err := s.rdb.Del(ctx, codeKey).Err(); err != nil && !errors.Is(err, redis.Nil) {
		return "", err
	}
	// 3.根据手机号查询用户
	err = s.db.WithContext(ctx).Where("phone = ?", loginForm.Phone).First(&user).Error
	// 4.用户不存在则创建
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user = model.User{
			Phone:    loginForm.Phone,
			NickName: utils.USER_NICK_NAME_PREFIX + utils.RandomString(10),
		}
		if err := s.db.WithContext(ctx).Create(&user).Error; err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}
	// 5.生成登录令牌并写入Redis
	token := uuid.NewString()
	//userDTO := dto.UserDTO{ID: user.ID, NickName: user.NickName, Icon: user.Icon}
	userDTO := mapper.ToUserDTO(&user)
	tokenKey := utils.LOGIN_USER_KEY + token
	// 将 UserDTO 中的字段完整序列化到 Redis Hash，便于后续统一读取
	data := map[string]string{
		"id":       strconv.FormatInt(userDTO.ID, 10),
		"nickName": userDTO.NickName,
		"icon":     userDTO.Icon,
	}
	if err := s.rdb.HSet(ctx, tokenKey, data).Err(); err != nil {
		return "", err
	}
	// 设置过期时间
	if err := s.rdb.Expire(ctx, tokenKey, time.Duration(utils.LOGIN_USER_TTL)*time.Second).Err(); err != nil {
		return "", err
	}
	// 返回 token
	return token, nil
}

func (s *UserService) FindByID(ctx context.Context, id int64) (*model.User, error) {
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
