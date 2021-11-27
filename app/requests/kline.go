package requests

type QueryHistory struct {
	Symbol string `json:"symbol" form:"symbol"`
	Period string `json:"period" form:"period"`
}
