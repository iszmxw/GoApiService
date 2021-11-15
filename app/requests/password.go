package requests

// Password 忘记密码修改密码
type Password struct {
	Password   string `json:"password" form:"password" validate:"required,min=6"`                        // 登录密码
	RePassword string `json:"re_password" form:"re_password" validate:"required,min=6,eqfield=Password"` // 确认登录密码
}

// Pw 登录用户修改密码
type Pw struct {
	Password   string `json:"password" form:"password" validate:"required,min=6"`                        // 登录密码
	RePassword string `json:"re_password" form:"re_password" validate:"required,min=6,eqfield=Password"` // 确认登录密码
}
