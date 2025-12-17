package dto

// UserDTO 用户dto
type UserDTO struct {
	ID       int64  `json:"id"`
	NickName string `json:"nickName"`
	Icon     string `json:"icon"`
}
