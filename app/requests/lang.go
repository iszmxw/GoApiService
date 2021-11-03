package requests

// Lang 语言
type Lang struct {
	Language string `json:"language" form:"language" validate:"required"` // 语言
}
