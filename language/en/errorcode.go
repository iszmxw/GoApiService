package en

// 英语

type ErrorCode struct {
	ParsingError                        string `code:"10000" msg:"Data analysis error"`
	LoginError                          string `code:"10001" msg:"Login failed, or your Token has expired"`
	PwError                             string `code:"10002" msg:"wrong password"`
	HtmlParsingError                    string `code:"10003" msg:"HTML analysis failed"`
	SendEmail                           string `code:"10004" msg:"Mail delivery failed"`
	ValidatorError                      string `code:"10005" msg:"Please check the parameters, parameter transmission error"`
	AddError                            string `code:"10006" msg:"Adding data failed"`
	UserIsExist                         string `code:"10007" msg:"This user has been registered"`
	ShareCodeIsExist                    string `code:"10008" msg:"Invitation code does not exist"`
	VerCodeErr                          string `code:"10009" msg:"The mailbox verification code is incorrect"`
	ResetPassword                       string `code:"10010" msg:"Retrieve password failure"`
	LangSetUp                           string `code:"10011" msg:"Language setting failed"`
	PayPasswordSetup                    string `code:"10012" msg:"Payment password setting failed"`
	CurrencyIsExist                     string `code:"10013" msg:"Currency does not exist"`
	SysCurrencyIsExist                  string `code:"10014" msg:"The current system currency is not set, please contact the staff to set the currency"`
	InsufficientBalance                 string `code:"10015" msg:"Insufficient balance"`
	Percentage                          string `code:"10016" msg:"Percent / transaction data error"`
	LimitPrice                          string `code:"10017" msg:"The price limit parameter cannot be empty"`
	UserIsNotExist                      string `code:"10018" msg:"User does not exist"`
	ContractIsNotExist                  string `code:"10019" msg:"Contain information is not inquired, please contact the administrator to configure the contract information of the current currency."`
	SysTradingPairIsExist               string `code:"10020" msg:"System is not configured"`
	TradingPairIsNotExist               string `code:"10021" msg:"Trading is not existed"`
	withdrawalFeesIsNotExist            string `code:"10022" msg:"The system is not configured with the recovery fee, and the cash is not supported. Please contact the staff."`
	UsersWalletIsNotExist               string `code:"10023" msg:"The transaction does not exist on the wallet"`
	ApplyBuySetupIsNotExist             string `code:"10024" msg:"Application of currency does not exist"`
	PasswordEditError                   string `code:"10025" msg:"Login password modification failed"`
	OperationFailed                     string `code:"10026" msg:"operation failed"`
	FeeOptionContractIsError            string `code:"10027" msg:"The option transaction contract fee has not been set, please contact the operation configuration here"`
	WithdrawalFeesIsError               string `code:"10027" msg:"Abnormal withdrawal fee"`
	LoginFailed                         string `code:"10028" msg:"Service exception, login failed!"`
	MultipleIsError                     string `code:"10029" msg:"The multiple value is illegal, not the value set by the system!"`
	ContractIsNotCorrect                string `code:"10030" msg:"The contract information is not set correctly, please contact the administrator to set the corresponding contract information for the transaction pair"`
	UsersWalletLock                     string `code:"10031" msg:"Sorry, your wallet is temporarily locked"`
	LiquidationUnsuccessful             string `code:"10032" msg:"Unsuccessful"`
	OptionContractSecondsNotCorrect     string `code:"10033" msg:"Incorrect number of seconds in options contract"`
	OptionContractStatus                string `code:"10034" msg:"This option contract has not been activated yet!"`
	FeePerpetualContractIsError         string `code:"10035" msg:"The perpetual contract fee has not been set, please contact the operation configuration!"`
	MinAmountIsNotExist                 string `code:"10036" msg:"No minimum withdrawal amount has been set, please contact Operation Configuration!"`
	OptionContractProfitRatioNotCorrect string `code:"10037" msg:"The option contract rate of return is incorrect"`
	ApplyBuySetupStatusIsNotExist       string `code:"10037" msg:"The subscription currency has not been opened yet"`
	CurrencyTransactionIsExist          string `code:"10038" msg:"The system temporarily does not allow transactions in the currency to be changed"`
	LimitPriceErr                       string `code:"10039" msg:"Limit price parameter error"`
	EntrustNumErr                       string `code:"10039" msg:"Order amount parameter error"`
}
