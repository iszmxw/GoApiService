package models

import (
	"goapi/app/response"
	"goapi/pkg/config"
	"goapi/pkg/helpers"
	"goapi/pkg/logger"
	"goapi/pkg/mysql"
	"gorm.io/gorm"
)

type CurrencyTransaction struct {
	Id              int                `gorm:"column:id" json:"id" form:"id"`                                           // 主键id
	UserId          int                `gorm:"column:user_id" json:"user_id" form:"user_id"`                            // 主键id
	OrderNumber     string             `gorm:"column:order_number" json:"order_number" form:"order_number"`             // 订单号
	Email           string             `gorm:"column:email" json:"email" form:"email"`                                  // 邮箱
	CurrencyId      int                `gorm:"column:currency_id" json:"currency_id" form:"currency_id"`                // 币种
	CurrencyName    string             `gorm:"column:currency_name" json:"currency_name" form:"currency_name"`          // 币种名称 例如：BTC/USDT（币种/交易对）
	EntrustNum      string             `gorm:"column:entrust_num" json:"entrust_num" form:"entrust_num"`                // 委托量
	LimitPrice      string             `gorm:"column:limit_price" json:"limit_price" form:"limit_price"`                // 当前限价
	ClinchNum       string             `gorm:"column:clinch_num" json:"clinch_num" form:"clinch_num"`                   // 成交量
	OrderType       string             `gorm:"column:order_type" json:"order_type" form:"order_type"`                   // 挂单类型：1-限价 2-市价
	TransactionType string             `gorm:"column:transaction_type" json:"transaction_type" form:"transaction_type"` // 订单方向：1-买入 2-卖出
	Price           string             `gorm:"column:price" json:"price" form:"price"`                                  // 挂单价格
	Status          string             `gorm:"column:status" json:"status" form:"status"`                               // 状态：0 交易中 1 已完成 2 已撤单
	Remark          string             `gorm:"column:remark" json:"remark" form:"remark"`                               // 备注
	CreatedAt       helpers.TimeNormal `gorm:"column:created_at" json:"created_at" form:"created_at"`                   // 创建时间
	UpdatedAt       helpers.TimeNormal `gorm:"column:updated_at" json:"updated_at" form:"updated_at"`                   // 更新时间
	DeletedAt       gorm.DeletedAt     `gorm:"column:deleted_at" json:"deleted_at" form:"deleted_at"`                   // 删除时间，为 null 则是没删除
}

func (m *CurrencyTransaction) tableName() string {
	prefix := config.GetString("database.mysql.prefix")
	table := "currency_transaction"
	return prefix + table
}

// GetPaginate 获取分页数据，返回错误
func (m *CurrencyTransaction) GetPaginate(where map[string]interface{}, orderBy interface{}, Page, Limit int64, lists *PageList) {
	var (
		result []response.CurrencyTransaction
	)
	// 获取表名
	tableName := m.tableName()
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
