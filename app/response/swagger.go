package response

type _LoginHandler struct {
	Code   string `json:"code"`
	Msg    string `json:"msg"`
	ReqId  string `json:"reqId"`
	Result struct {
		Token string `json:"token"` // 登录获取的token
		Uid   int    `json:"uid"`   // 用户id
	} `json:"result"`
	Success bool `json:"success"`
}

type _OK struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	ReqId   string `json:"reqId"`
	Result  string `json:"result"`
	Success bool   `json:"success"`
}
