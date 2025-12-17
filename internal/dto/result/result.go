package result

type Result struct {
	Success  bool        `json:"success"`
	ErrorMsg string      `json:"errorMsg"`
	Data     interface{} `json:"data"`
	Total    *int64      `json:"total"`
}

// Ok 返回一个不带数据的成功响应
func Ok() Result {
	return Result{Success: true}
}

// OkWithData 返回携带数据的成功响应
func OkWithData(data interface{}) Result {
	return Result{Success: true, Data: data}
}

// OkWithPage 返回分页成功响应
func OkWithPage(data interface{}, total int64) Result {
	return Result{Success: true, Data: data, Total: &total}
}

// Fail 返回失败响应
func Fail(msg string) Result {
	return Result{Success: false, ErrorMsg: msg}
}

// ScrollResult 滚动分页的辅助响应
type ScrollResult struct {
	List    interface{} `json:"list"`
	MinTime int64       `json:"minTime"`
	Offset  int         `json:"offset"`
}
