package response

import (
	"goapi/pkg/helpers"
	"gorm.io/gorm"
)

// ApplyBuySetup 申购币种设置
type ApplyBuySetup struct {
	Id            uint               `json:"id"`             //主键id
	Name          string             `json:"name"`           //币种名称
	IssuePrice    float64            `json:"issue_price"`    //发行价 1 = 多少个USDT
	EstimatedTime helpers.TimeNormal `json:"estimated_time"` //预计上线时间
	StartTime     helpers.TimeNormal `json:"start_time"`     //开始申购时间
	EndTime       helpers.TimeNormal `json:"end_time"`       //结束申购时间
	Code          string             `json:"code"`           //申购码
	CodeStatus    int8               `json:"code_status"`    //申购码开关  0 关闭 1 开启
	Detail        string             `json:"detail"`         //项目详情
	Status        int8               `json:"status"`         //币种状态  0 关闭 1 开启
	CreatedAt     helpers.TimeNormal `json:"created_at"`     //创建时间
	UpdatedAt     helpers.TimeNormal `json:"updated_at"`     //更新时间
	DeletedAt     gorm.DeletedAt     `json:"deleted_at"`     //删除时间，为 null 则是没删除
}
