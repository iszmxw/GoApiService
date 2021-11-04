package models

import (
	"goapi/app/response"
	"goapi/pkg/config"
	"goapi/pkg/helpers"
	"goapi/pkg/logger"
	"goapi/pkg/mysql"
	"gorm.io/gorm"
)

type PerpetualContractTransaction struct {
	Id              int                `json:"id"`                //主键id
	UserId          int                `json:"user_id"`           //用户id
	Email           string             `json:"email"`             //用户邮箱
	OrderNumber     string             `json:"order_number"`      //订单号
	CurrencyId      int                `json:"currency_id"`       //币种
	CurrencyName    string             `json:"currency_name"`     //币种名称 例如：BTC/USDT（币种/交易对）
	TradingPairId   int                `json:"trading_pair_id"`   //交易对id
	TradingPairName string             `json:"trading_pair_name"` //交易对名称
	KLineCode       string             `json:"k_line_code"`       //K线图代码
	OrderType       string             `json:"order_type"`        //订单类型：1-限价 2-市价
	LimitPrice      string             `json:"limit_price"`       //当前限价
	TransactionType string             `json:"transaction_type"`  //交易类型：1-开多 2-开空
	EntrustNum      string             `json:"entrust_num"`       //委托量
	EntrustPrice    string             `json:"entrust_price"`     //委托价格
	EnsureAmount    string             `json:"ensure_amount"`     //保证金
	HandleFee       string             `json:"handle_fee"`        //手续费，单位百分比
	Multiple        string             `json:"multiple"`          //倍数值
	Price           string             `json:"price"`             //交易金额
	Income          float64            `json:"income"`            //最终收益
	Status          string             `json:"status"`            //状态：0 交易中 1 已完成 2 已撤单
	Remark          string             `json:"remark"`            //备注
	CreatedAt       helpers.TimeNormal `json:"created_at"`        //创建时间
	UpdatedAt       helpers.TimeNormal `json:"updated_at"`        //更新时间
	DeletedAt       gorm.DeletedAt     `json:"deleted_at"`        //删除时间，为 null 则是没删除
}

func (m *PerpetualContractTransaction) TableName() string {
	prefix := config.GetString("database.mysql.prefix")
	table := "perpetual_contract_transaction"
	return prefix + table
}

// GetPaginate 获取分页数据，返回错误
func (m *PerpetualContractTransaction) GetPaginate(where map[string]interface{}, orderBy interface{}, Page, Limit int64, lists *PageList) {
	var (
		result []response.PerpetualContractTransaction
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
	} else {
		table = table.Order("id desc")
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
