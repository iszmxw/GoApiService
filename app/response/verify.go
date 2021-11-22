package response

import (
	"goapi/pkg/helpers"
)

type Verify struct {
	UserId        int                `json:"user_id"`          //用户id从token中拿
	IdentityCard  int                `json:"identity_card" `   //身份证信息
	Status        int                `json:"status" `          //状态0未提交 1初级验证  2高级验证
	Email         string             `json:"email " `          //邮箱
	FullName      string             `json:"full_name" `       //用户真实姓名
	ImgCardBehind string             `json:"img_card_behind" ` //身份证后面图片
	ImgCardFront  string             `json:"img_card_front" `  //身份证前面图片
	ImgBankFront  string             `json:"img_bank_front" `  //银行卡前面
	ImgBankBehind string             `json:"img_bank_behind" ` //银行卡后面
	CreatedAt     helpers.TimeNormal `json:"created_at"`       // 创建时间
	UpdatedAt     helpers.TimeNormal `json:"updated_at"`       // 更新时间
	DeletedAt     helpers.TimeNormal `json:"deleted_at"`       // 删除时间，为 null 则是没删除
}
