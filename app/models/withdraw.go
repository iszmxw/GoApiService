package models

import (
	"goapi/app/response"
	"goapi/pkg/config"
	"goapi/pkg/helpers"
	"goapi/pkg/logger"
	"goapi/pkg/mysql"
)

// Withdraw 用户提现订单记录表
type Withdraw struct {
	Id              int                `json:"id"`
	UserId          int                `json:"user_id"`           //用户id
	Email           string             `json:"email"`             //用户邮箱
	Address         string             `json:"address"`           //提币地址
	TradingPairId   string             `json:"trading_pair_id"`   //提现的交易对id
	TradingPairName string             `json:"trading_pair_name"` //提现的交易对name
	Type            int             `json:"type"`              //币链类型 1-OMNI 2-ERC20 3-TRC20
	WithdrawNum     float64            `json:"withdraw_num"`      //提现数量
	HandlingFee     float64            `json:"handling_fee"`      //手续费
	ActuallyArrived float64            `json:"actually_arrived"`  //实际到账
	Status          string             `json:"status"`            //状态：0-未确认：1-已确认 2-已拒绝
	Remark          string             `json:"remark"`            //备注
	CreatedAt       helpers.TimeNormal `json:"created_at"`        //创建时间
	UpdatedAt       helpers.TimeNormal `json:"updated_at"`        //更新时间
	DeletedAt       helpers.TimeNormal `json:"deleted_at"`        //删除时间，为 null 则是没删除
}

func (m *Withdraw) TableName() string {
	prefix := config.GetString("database.mysql.prefix")
	table := "withdraw"
	return prefix + table
}

// GetPaginate 获取分页数据，返回错误
func (m *Withdraw) GetPaginate(where map[string]interface{}, orderBy interface{}, Page, Limit int64, lists *PageList) {
	var (
		result []response.Withdraw
	)
	// 获取表名
	tableName := m.TableName()
	table := mysql.DB.Debug().Table(Prefix(tableName))
	table = table.Where(where)
	table.Count(&lists.Total)
	// 设置分页参数
	lists.CurrentPage = Page
	lists.PageSize = Limit
	InitPageList(lists)
	// order by
	if len(orderBy.(string)) > 0 {
		table = table.Order(orderBy)
	}
	table = table.Offset(int(lists.Offset))
	table = table.Limit(int(lists.PageSize))
	if err := table.Scan(&result).Error; err != nil {
		// 记录错误
		logger.Error(err)
	} else {
		lists.Data = result
	}
}
