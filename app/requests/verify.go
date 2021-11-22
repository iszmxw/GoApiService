package requests

//用户验证 图片上传结构体

type VerifyParam struct {
	IdentityCard string `json:"identity_card" form:"identity_card" validate:"required"` //身份证信息
	FullName     string `json:"full_name" form:"full_name" validate:"required"`         //用户真实姓名
	//ImgUrl       string `json:"img_url" form:"img_url" validate:"omitempty"`
}

type ImgBase64Param struct {
	ImgCardFront  string `json:"img_card_front" form:"img_card_front" validate:"required"`
	ImgCardBehind string `json:"img_card_behind" form:"img_card_behind" validate:"required"`
	ImgBankFront  string `json:"img_bank_front" form:"img_bank_front" validate:"required"`
	ImgBankBehind string `json:"img_bank_behind" form:"img_bank_behind" validate:"required"`
}
