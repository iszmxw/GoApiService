package requests

// 更局交易类型查询系统币种

type SysCurrencyType struct {
	Type string `json:"type" form:"type"` // 1 现货交易的币种 2 期权合约交易的币种 3 永续合约交易的币种
}
