package zh

// 中文

type ErrorCode1 struct {
	ParsingError             string `code:"10000" msg:"数据解析错误"`
	LoginError               string `code:"10001" msg:"登录失败，或者您的token已经过期"`
	PwError                  string `code:"10002" msg:"密码错误"`
	HtmlParsingError         string `code:"10003" msg:"Html解析失败"`
	SendEmail                string `code:"10004" msg:"邮件发送失败"`
	ValidatorError           string `code:"10005" msg:"请检查参数，参数传输错误"`
	AddError                 string `code:"10006" msg:"添加数据失败"`
	UserIsExist              string `code:"10007" msg:"该用户已注册"`
	ShareCodeIsExist         string `code:"10008" msg:"邀请码不存在"`
	VerCodeErr               string `code:"10009" msg:"邮箱验证码不正确"`
	ResetPassword            string `code:"10010" msg:"找回密码失败"`
	LangSetUp                string `code:"10011" msg:"语言设置失败"`
	PayPasswordSetup         string `code:"10012" msg:"支付密码设置失败"`
	CurrencyIsExist          string `code:"10013" msg:"币种不存在"`
	SysCurrencyIsExist       string `code:"10014" msg:"当前系统币种未设置，请联系工作人员设置币种"`
	InsufficientBalance      string `code:"10015" msg:"可用余额不足"`
	Percentage               string `code:"10016" msg:"百分比/成交量数据错误"`
	LimitPrice               string `code:"10017" msg:"限价参数不能为空"`
	UserIsNotExist           string `code:"10018" msg:"用户不存在"`
	ContractIsNotExist       string `code:"10019" msg:"暂未查询到合约信息，请联系管理员配置当前币种的合约信息"`
	SysTradingPairIsExist    string `code:"10020" msg:"系统未配置交易对"`
	TradingPairIsNotExist    string `code:"10021" msg:"交易对不存在"`
	withdrawalFeesIsNotExist string `code:"10022" msg:"系统未配置提现手续费，暂不支持提现，请联系工作人员"`
	UsersWalletIsNotExist    string `code:"10023" msg:"该交易对钱包不存在"`
	ApplyBuySetupIsNotExist  string `code:"10024" msg:"申购币种不存在"`
	PasswordEditError        string `code:"10025" msg:"登录密码修改失败"`
	OperationFailed          string `code:"10026" msg:"操作失败"`
	FeeOptionContractIsError string `code:"10027" msg:"期权交易合约手续费未设置，请里联系运营配置"`
	WithdrawalFeesIsError    string `code:"10027" msg:"提现手续费异常"`
	LoginFailed              string `code:"10028" msg:"服务异常，登录失败！"`
}
