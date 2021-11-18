package ko

// 韩语

type ErrorCode struct {
	ParsingError                        string `code:"10000" msg:"데이터 분석 오류"`
	LoginError                          string `code:"10001" msg:"로그인이 실패했거나 토큰이 만료되었습니다"`
	PwError                             string `code:"10002" msg:"잘못된 비밀번호"`
	HtmlParsingError                    string `code:"10003" msg:"HTML 분석에 실패했습니다"`
	SendEmail                           string `code:"10004" msg:"메일 배달 실패"`
	ValidatorError                      string `code:"10005" msg:"매개 변수, 파라미터 전송 오류를 확인하십시오"`
	AddError                            string `code:"10006" msg:"데이터 추가 실패"`
	UserIsExist                         string `code:"10007" msg:"이 사용자가 등록되었습니다"`
	ShareCodeIsExist                    string `code:"10008" msg:"초대 코드가 없습니다"`
	VerCodeErr                          string `code:"10009" msg:"사서함 확인 코드가 잘못되었습니다"`
	ResetPassword                       string `code:"10010" msg:"암호 실패를 검색하십시오"`
	LangSetUp                           string `code:"10011" msg:"언어 설정 실패"`
	PayPasswordSetup                    string `code:"10012" msg:"지불 암호 설정 실패"`
	CurrencyIsExist                     string `code:"10013" msg:"통화가 없습니다"`
	SysCurrencyIsExist                  string `code:"10014" msg:"현재 시스템 통화가 설정되어 있지 않으므로 직원에게 문의하여 통화를 설정하십시오."`
	InsufficientBalance                 string `code:"10015" msg:"잔액 불충분"`
	Percentage                          string `code:"10016" msg:"백분율 / 트랜잭션 데이터 오류"`
	LimitPrice                          string `code:"10017" msg:"가격 제한 매개 변수는 비어있을 수 없습니다"`
	UserIsNotExist                      string `code:"10018" msg:"사용자가 존재하지 않습니다"`
	ContractIsNotExist                  string `code:"10019" msg:"정보가 포함되어 있지 않으므로 관리자에게 문의하여 현재 통화의 계약 정보를 구성하십시오."`
	SysTradingPairIsExist               string `code:"10020" msg:"시스템이 구성되어 있지 않습니다"`
	TradingPairIsNotExist               string `code:"10021" msg:"거래는 존재하지 않습니다"`
	withdrawalFeesIsNotExist            string `code:"10022" msg:"이 시스템은 복구 수수료로 구성되지 않으며 현금은 지원되지 않습니다. 직원에게 문의하십시오."`
	UsersWalletIsNotExist               string `code:"10023" msg:"거래가 지갑에 존재하지 않습니다"`
	ApplyBuySetupIsNotExist             string `code:"10024" msg:"통화 적용이 존재하지 않습니다"`
	PasswordEditError                   string `code:"10025" msg:"로그인 암호 수정 실패"`
	OperationFailed                     string `code:"10026" msg:"작업 실패"`
	FeeOptionContractIsError            string `code:"10027" msg:"옵션 거래 계약 수수료가 설정되지 않았습니다. 여기에서 작업 구성에 문의하십시오."`
	WithdrawalFeesIsError               string `code:"10028" msg:"비정상적인 출금 수수료"`
	LoginFailed                         string `code:"10029" msg:"서비스 예외, 로그인 실패!"`
	MultipleIsError                     string `code:"10030" msg:"시스템에서 설정한 값이 아닌 다중 값이 잘못되었습니다!"`
	ContractIsNotCorrect                string `code:"10031" msg:"계약 정보가 올바르게 설정되지 않았습니다. 거래 쌍에 대한 해당 계약 정보를 설정하려면 관리자에게 문의하십시오."`
	UsersWalletLock                     string `code:"10032" msg:"죄송합니다. 지갑이 일시적으로 잠겨 있습니다."`
	LiquidationUnsuccessful             string `code:"10033" msg:"실패"`
	OptionContractSecondsNotCorrect     string `code:"10034" msg:"옵션 계약의 잘못된 시간(초)"`
	OptionContractStatus                string `code:"10035" msg:"이 옵션 계약은 아직 활성화되지 않았습니다!"`
	FeePerpetualContractIsError         string `code:"10036" msg:"무기한 계약 수수료가 설정되지 않았습니다. 운영 구성에 문의하십시오!"`
	MinAmountIsNotExist                 string `code:"10037" msg:"최소 출금 금액이 설정되지 않았습니다. Operation Configuration에 문의하십시오!"`
	OptionContractProfitRatioNotCorrect string `code:"10038" msg:"옵션 계약 수익률이 잘못되었습니다."`
	ApplyBuySetupStatusIsNotExist       string `code:"10039" msg:"구독 통화가 아직 열리지 않았습니다."`
	CurrencyTransactionIsExist          string `code:"10040" msg:"시스템은 일시적으로 통화의 거래를 변경하는 것을 허용하지 않습니다."`
	LimitPriceErr                       string `code:"10041" msg:"제한 가격 매개변수 오류"`
	EntrustNumErr                       string `code:"10042" msg:"주문 금액 매개변수 오류"`
	CurrencyTypeIsNotAllowed            string `code:"10043" msg:"이 유형의 거래는 이 거래 쌍에 대해 일시적으로 허용되지 않습니다. 필요한 경우 작업에 문의하십시오."`
	UserIsLock                          string `code:"10044" msg:"사용자 상태가 잠겨 있습니다."`
	SearchTimeErr                       string `code:"10045" msg:"검색 시간 형식이 잘못되었습니다(예: 2006-01-02)."`
	OptionContractMinimum               string `code:"10046" msg:"현재 거래량이 가장 많이 소비된 것보다 적습니다."`
}
