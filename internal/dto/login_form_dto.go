package dto

// LoginForm 登录表单
type LoginForm struct {
	Phone    string `json:"phone"`
	Code     string `json:"code"`
	Password string `json:"password"`
}
