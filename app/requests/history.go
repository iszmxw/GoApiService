package requests

// History 历史委托
type History struct {
	Page       int    `json:"page" form:"page"`
	Limit      int    `json:"limit" form:"limit"`
	OrderBy    string `json:"orderBy" form:"orderBy"`
	Status     string `json:"status" form:"status" validate:"omitempty,min=0"` // 状态：0 交易中 1 已完成 2 已撤单
	CurrencyId string `json:"currency_id" form:"currency_id"`                  // 查询的币种id
}
