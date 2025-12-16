package dto

// Result matches the Java Result DTO to keep API parity.
type Result struct {
	Success  bool        `json:"success"`
	ErrorMsg string      `json:"errorMsg"`
	Data     interface{} `json:"data"`
	Total    *int64      `json:"total"`
}

// Ok returns a successful response without payload.
func Ok() Result {
	return Result{Success: true}
}

// OkWithData returns a successful response with data payload.
func OkWithData(data interface{}) Result {
	return Result{Success: true, Data: data}
}

// OkWithPage returns a paginated success response.
func OkWithPage(data interface{}, total int64) Result {
	return Result{Success: true, Data: data, Total: &total}
}

// Fail returns a failure response.
func Fail(msg string) Result {
	return Result{Success: false, ErrorMsg: msg}
}

// LoginForm mirrors LoginFormDTO from Java.
type LoginForm struct {
	Phone    string `json:"phone"`
	Code     string `json:"code"`
	Password string `json:"password"`
}

// UserDTO mirrors the Java UserDTO.
type UserDTO struct {
	ID       int64  `json:"id"`
	NickName string `json:"nickName"`
	Icon     string `json:"icon"`
}

// ScrollResult is a helper response for scroll pagination.
type ScrollResult struct {
	List    interface{} `json:"list"`
	MinTime int64       `json:"minTime"`
	Offset  int         `json:"offset"`
}
