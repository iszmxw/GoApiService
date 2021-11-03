package models

import (
	"goapi/app/response"
	"goapi/pkg/config"
	"goapi/pkg/mysql"
	"gorm.io/gorm"
	"time"
)

// User 用户表
type User struct {
	Id            int            `gorm:"column:id"`             // 主键id
	Language      string         `gorm:"column:language"`       // 语言：1-繁体 2-英文 3-日文
	IsAgent       string         `gorm:"column:is_agent"`       // 是否代理： 0：不是   1：是
	Email         string         `gorm:"column:email"`          // 邮箱
	Nickname      string         `gorm:"column:nickname"`       // 昵称
	Password      string         `gorm:"column:password"`       // 登录密码
	PayPassword   string         `gorm:"column:pay_password"`   // 支付密码
	UserLevel     int            `gorm:"column:user_level"`     // 用户层级
	UserPath      string         `gorm:"column:user_path"`      // 用户关系
	PartnerLevel  int            `gorm:"column:partner_level"`  // 合伙人等级
	AgentDividend string         `gorm:"column:agent_dividend"` // 代理红利
	ShareCode     string         `gorm:"column:share_code"`     // 用户邀请码
	RiskProfit    int            `gorm:"column:risk_profit"`    // 风控概率
	LastLoginIp   string         `gorm:"column:last_login_ip"`  // 登录IP
	Status        string         `gorm:"column:status"`         // 状态： 0正常 1，已锁定
	LockTime      *time.Time     `gorm:"column:lock_time"`      // 锁定时间
	LoginTime     time.Time      `gorm:"column:login_time"`     // 登录时间
	CreatedAt     time.Time      `gorm:"column:created_at"`     // 创建时间
	UpdatedAt     time.Time      `gorm:"column:updated_at"`     // 更新时间
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at"`     // 删除时间，为 null 则是没删除
}

func (m *User) TableName() string {
	prefix := config.GetString("database.mysql.prefix")
	table := "user"
	return prefix + table
}

func (m User) Add(users *User) error {
	DB := mysql.DB.Debug()
	return DB.Create(&users).Error
}

func (m User) Update(user *User) error {
	DB := mysql.DB.Debug()
	return DB.Updates(&user).Error
}

func (m User) GetOne(where map[string]interface{}, users *response.User) {
	DB := mysql.DB.Debug()
	DB.Model(m).Where(where).Find(&users)
}

func (m *User) SelectDelete(where map[string]interface{}, users *User) {
	DB := mysql.DB.Debug()
	DB.Where(where).Delete(users)
}
