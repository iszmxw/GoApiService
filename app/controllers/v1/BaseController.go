package v1

type BaseController struct {
}

type Group struct {
	BaseController
	LoginController
	IndexController
	UserController
	OrderController
	TradeController
	AssetsStreamController
	OptionContractController
	PerpetualContractController
	CurrencyCurrencyController
	KlineController
}
