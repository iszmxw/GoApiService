package response

import (
	"goapi/pkg/helpers"
	"gorm.io/gorm"
)

// 用户信息表

type User struct {
	Id            int                `json:"id"`             //主键id
	Language      string             `json:"language"`       //语言：1-繁体 2-英文 3-日文 4-韩语 5 西班牙语
	IsAgent       string             `json:"is_agent"`       //是否代理： 0：不是   1：是
	ParentId      int                `json:"parent_id"`      //上级用户id
	Email         string             `json:"email"`          //邮箱
	Nickname      string             `json:"nickname"`       //昵称
	Password      string             `json:"password"`       //登录密码
	PayPassword   string             `json:"pay_password"`   //支付密码
	UserLevel     int                `json:"user_level"`     //用户层级
	UserPath      string             `json:"user_path"`      //用户关系
	AgentDividend string             `json:"agent_dividend"` //代理红利
	ShareCode     string             `json:"share_code"`     //用户邀请码，每个用户唯一
	RiskProfit    int                `json:"risk_profit"`    //风控 0-无 1-盈 2-亏
	LastLoginIp   string             `json:"last_login_ip"`  //登录IP
	Status        string             `json:"status"`         //状态： 0正常 1，已锁定
	LockTime      helpers.TimeNormal `json:"lock_time"`      //锁定时间
	LoginTime     helpers.TimeNormal `json:"login_time"`     //登录时间
	CreatedAt     helpers.TimeNormal `json:"created_at"`     //创建时间
	UpdatedAt     helpers.TimeNormal `json:"updated_at"`     //更新时间
	DeletedAt     gorm.DeletedAt     `json:"deleted_at"`     //删除时间，为 null 则是没删除
}
