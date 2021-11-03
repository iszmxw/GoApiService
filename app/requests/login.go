package requests

// UserLogin 用户登录
type UserLogin struct {
	Email    string `json:"email" form:"email" validate:"required,email"`    // 邮箱
	Password string `json:"password" form:"password" validate:"required,min=6"` // 登录密码
}
