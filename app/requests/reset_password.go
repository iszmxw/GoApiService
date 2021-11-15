package requests

// ResetVerify 重置密码前验证邮件Code
type ResetVerify struct {
	Email string `json:"email" form:"email" validate:"required,email,contains=@"` // 邮箱
	Code  string `json:"code" form:"code" validate:"required,numeric"`            // 验证码
}

// ResetPassword 重置密码
type ResetPassword struct {
	Email      string `json:"email" form:"email" validate:"required,email,contains=@"`          // 邮箱
	Code       string `json:"code" form:"code" validate:"required,numeric"`                     // 验证码
	Password   string `json:"password" form:"password" validate:"min=6"`                        // 密码
	RePassword string `json:"re_password" form:"re_password" validate:"min=6,eqfield=Password"` // 确认密码
}

// EditPassword 修改密码
type EditPassword struct {
	Password   string `json:"password" form:"password" validate:"min=6"`                        // 密码
	RePassword string `json:"re_password" form:"re_password" validate:"min=6,eqfield=Password"` // 确认密码
}
