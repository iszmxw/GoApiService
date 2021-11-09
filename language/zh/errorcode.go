package zh

// 中文

type ErrorCode struct {
	ParsingError                        string `code:"10000" msg:"數據解析錯誤"`
	LoginError                          string `code:"10001" msg:"登錄失敗，或者您的token已經過期"`
	PwError                             string `code:"10002" msg:"密碼錯誤"`
	HtmlParsingError                    string `code:"10003" msg:"Html解析失敗"`
	SendEmail                           string `code:"10004" msg:"郵件發送失敗"`
	ValidatorError                      string `code:"10005" msg:"請檢查參數，參數傳輸錯誤"`
	AddError                            string `code:"10006" msg:"添加數據失敗"`
	UserIsExist                         string `code:"10007" msg:"該用戶已註冊"`
	ShareCodeIsExist                    string `code:"10008" msg:"邀請碼不存在"`
	VerCodeErr                          string `code:"10009" msg:"郵箱驗證碼不正確"`
	ResetPassword                       string `code:"10010" msg:"找回密碼失敗"`
	LangSetUp                           string `code:"10011" msg:"語言設置失敗"`
	PayPasswordSetup                    string `code:"10012" msg:"支付密碼設置失敗"`
	CurrencyIsExist                     string `code:"10013" msg:"幣種不存在"`
	SysCurrencyIsExist                  string `code:"10014" msg:"當前系統幣種未設置，請聯繫工作人員設置幣種"`
	InsufficientBalance                 string `code:"10015" msg:"可用餘額不足"`
	Percentage                          string `code:"10016" msg:"百分比/成交量數據錯誤"`
	LimitPrice                          string `code:"10017" msg:"限價參數不能為空"`
	UserIsNotExist                      string `code:"10018" msg:"用戶不存在"`
	ContractIsNotExist                  string `code:"10019" msg:"暫未查詢到合約信息，請聯繫管理員配置當前幣種的合約信息"`
	SysTradingPairIsExist               string `code:"10020" msg:"系統未配置交易對"`
	TradingPairIsNotExist               string `code:"10021" msg:"交易對不存在"`
	withdrawalFeesIsNotExist            string `code:"10022" msg:"系統未配置提現手續費，暫不支持提現，請聯繫工作人員"`
	UsersWalletIsNotExist               string `code:"10023" msg:"該交易對錢包不存在"`
	ApplyBuySetupIsNotExist             string `code:"10024" msg:"申購幣種不存在"`
	PasswordEditError                   string `code:"10025" msg:"登錄密碼修改失敗"`
	OperationFailed                     string `code:"10026" msg:"操作失敗"`
	FeeOptionContractIsError            string `code:"10027" msg:"期權交易合約手續費未設置，請里聯系運營配置"`
	WithdrawalFeesIsError               string `code:"10027" msg:"提現手續費異常"`
	LoginFailed                         string `code:"10028" msg:"服務異常，登錄失敗！"`
	MultipleIsError                     string `code:"10029" msg:"倍數值不合法，不是系統設定的值！"`
	ContractIsNotCorrect                string `code:"10030" msg:"合約信息未正確設置，請聯繫管理員設置該交易對相應的合約信息"`
	UsersWalletLock                     string `code:"10031" msg:"對不起，您的錢包暫時已被鎖定"`
	LiquidationUnsuccessful             string `code:"10032" msg:"平倉失敗"`
	OptionContractSecondsNotCorrect     string `code:"10033" msg:"期權合約秒數不正確"`
	OptionContractStatus                string `code:"10034" msg:"該期權合約暫未啟用！"`
	FeePerpetualContractIsError         string `code:"10035" msg:"永續合約手續費未設置，請聯繫運營配置！"`
	MinAmountIsNotExist                 string `code:"10036" msg:"未設置最小提現金額，請聯繫運營配置！"`
	OptionContractProfitRatioNotCorrect string `code:"10037" msg:"期權合約收益率不正確"`
	ApplyBuySetupStatusIsNotExist       string `code:"10037" msg:"申購幣種暫未開啟"`
}
