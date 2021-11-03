package response

// GlobalsTypes 全局配置管理 类型获取
type GlobalsTypes struct {
	Id     uint   `json:"id"`     //主键id
	Fields string `json:"fields"` //字段名称
	Value  string `json:"value"`  //字段值
}

// WithdrawalFees 全局配置管理 类型获取
type WithdrawalFees struct {
	Id     uint    `json:"id"`     //主键id
	Fields string  `json:"fields"` //字段名称
	Value  float64 `json:"value"`  //字段值
}
