package response

type RespData struct {
	Data struct {
		Code    string      `json:"code"`
		Success bool        `json:"success"`
		Result  interface{} `json:"result"`
		Msg     string      `json:"msg"`
	} `json:"data"`
}
