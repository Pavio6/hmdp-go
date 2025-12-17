package mapper

import (
	"hmdp-backend/internal/dto"
	"hmdp-backend/internal/model"
)

// UserDTO 用户dto
type UserDTO struct {
	ID       int64  `json:"id"`
	NickName string `json:"nickName"`
	Icon     string `json:"icon"`
}

// ToUserDTO 将user对象转为userDTO
func ToUserDTO(u *model.User) *dto.UserDTO {
	if u == nil {
		return nil
	}
	return &dto.UserDTO{
		ID:       u.ID,
		Icon:     u.Icon,
		NickName: u.NickName,
	}
}
