package models

// ErrorResponse 通用错误响应结构体
type ErrorResponse struct {
	ErrorMsg  string `json:"error_msg"`         // 错误具体信息
	Details   string `json:"details,omitempty"` // 错误详情（可选字段）
	ErrorCode int    `json:"error_code"`        // 错误码
}
