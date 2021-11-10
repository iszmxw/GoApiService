package ja

// 日语

type ErrorCode struct {
	ParsingError                        string `code:"10000" msg:"データ分析エラー"`
	LoginError                          string `code:"10001" msg:"ログインに失敗しました、またはあなたのトークンの有効期限が切れています"`
	PwError                             string `code:"10002" msg:"間違ったパスワード"`
	HtmlParsingError                    string `code:"10003" msg:"HTML分析に失敗しました"`
	SendEmail                           string `code:"10004" msg:"メールの送信に失敗しました"`
	ValidatorError                      string `code:"10005" msg:"パラメータ、パラメータ送信エラーを確認してください"`
	AddError                            string `code:"10006" msg:"データの追加に失敗しました"`
	UserIsExist                         string `code:"10007" msg:"このユーザーは登録されています"`
	ShareCodeIsExist                    string `code:"10008" msg:"招待コードは存在しません"`
	VerCodeErr                          string `code:"10009" msg:"メールボックスの確認コードが正しくありません"`
	ResetPassword                       string `code:"10010" msg:"パスワード障害を取得します"`
	LangSetUp                           string `code:"10011" msg:"言語設定に失敗しました"`
	PayPasswordSetup                    string `code:"10012" msg:"支払パスワードの設定に失敗しました"`
	CurrencyIsExist                     string `code:"10013" msg:"通貨は存在しません"`
	SysCurrencyIsExist                  string `code:"10014" msg:"現在のシステム通貨が設定されていません、通貨を設定するためにスタッフに連絡してください"`
	InsufficientBalance                 string `code:"10015" msg:"残高不足です"`
	Percentage                          string `code:"10016" msg:"パーセント/トランザクションデータエラー"`
	LimitPrice                          string `code:"10017" msg:"価格制限パラメータを空にすることはできません"`
	UserIsNotExist                      string `code:"10018" msg:"ユーザーは存在しません"`
	ContractIsNotExist                  string `code:"10019" msg:"情報が尋ねられていない場合は、現在の通貨の契約情報を構成するには、管理者に連絡してください"`
	SysTradingPairIsExist               string `code:"10020" msg:"システムは構成されていません"`
	TradingPairIsNotExist               string `code:"10021" msg:"取引は存在しません"`
	withdrawalFeesIsNotExist            string `code:"10022" msg:"システムは回復料手帳では設定されておらず、現金はサポートされていません。スタッフに連絡してください。"`
	UsersWalletIsNotExist               string `code:"10023" msg:"トランザクションは財布に存在しません"`
	ApplyBuySetupIsNotExist             string `code:"10024" msg:"通貨の適用は存在しません"`
	PasswordEditError                   string `code:"10025" msg:"ログインパスワードの変更に失敗しました"`
	OperationFailed                     string `code:"10026" msg:"操作に失敗しました"`
	FeeOptionContractIsError            string `code:"10027" msg:"オプション取引契約手数料が設定されていませんので、こちらの運用構成にお問い合わせください"`
	WithdrawalFeesIsError               string `code:"10027" msg:"異常な引き出し手数料"`
	LoginFailed                         string `code:"10028" msg:"サービス例外、ログインに失敗しました！"`
	MultipleIsError                     string `code:"10029" msg:"複数の値は不正であり、システムによって設定された値ではありません！"`
	ContractIsNotCorrect                string `code:"10030" msg:"契約情報が正しく設定されていません。管理者に連絡して、トランザクションペアに対応する契約情報を設定してください。"`
	UsersWalletLock                     string `code:"10031" msg:"申し訳ありませんが、ウォレットは一時的にロックされています"`
	LiquidationUnsuccessful             string `code:"10032" msg:"失敗"`
	OptionContractSecondsNotCorrect     string `code:"10033" msg:"オプション契約の秒数が正しくありません"`
	OptionContractStatus                string `code:"10034" msg:"このオプション契約はまだ有効化されていません！"`
	FeePerpetualContractIsError         string `code:"10035" msg:"永久契約料は設定されていませんので、運用構成にお問い合わせください！"`
	MinAmountIsNotExist                 string `code:"10036" msg:"最小引き出し額は設定されていません。運用構成にお問い合わせください。"`
	OptionContractProfitRatioNotCorrect string `code:"10037" msg:"オプション契約の収益率が正しくない"`
	ApplyBuySetupStatusIsNotExist       string `code:"10037" msg:"サブスクリプション通貨はまだ開かれていません"`
	CurrencyTransactionIsExist          string `code:"10038" msg:"システムは一時的に通貨での取引を変更することを許可しません"`
	LimitPriceErr                       string `code:"10039" msg:"指値パラメータエラー"`
	EntrustNumErr                       string `code:"10039" msg:"注文金額パラメータエラー"`
}
