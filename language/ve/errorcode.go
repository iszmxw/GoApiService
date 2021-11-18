package ve

// 越南语

type ErrorCode struct {
	ParsingError                        string `code:"10000" msg:"Lỗi phân tích cú pháp dữ liệu"`
	LoginError                          string `code:"10001" msg:"Đăng nhập không thành công hoặc mã thông báo của bạn đã hết hạn"`
	PwError                             string `code:"10002" msg:"sai mật khẩu"`
	HtmlParsingError                    string `code:"10003" msg:"Phân tích cú pháp html không thành công"`
	SendEmail                           string `code:"10004" msg:"Gửi thư không thành công"`
	ValidatorError                      string `code:"10005" msg:"Vui lòng kiểm tra thông số, lỗi truyền thông số"`
	AddError                            string `code:"10006" msg:"Thêm dữ liệu không thành công"`
	UserIsExist                         string `code:"10007" msg:"Người dùng này đã được đăng ký"`
	ShareCodeIsExist                    string `code:"10008" msg:"Mã lời mời không tồn tại"`
	VerCodeErr                          string `code:"10009" msg:"Mã xác minh email không chính xác"`
	ResetPassword                       string `code:"10010" msg:"Lấy lại mật khẩu, không thành công"`
	LangSetUp                           string `code:"10011" msg:"Cài đặt ngôn ngữ không thành công"`
	PayPasswordSetup                    string `code:"10012" msg:"Đặt mật khẩu thanh toán không thành công"`
	CurrencyIsExist                     string `code:"10013" msg:"Tiền tệ không tồn tại"`
	SysCurrencyIsExist                  string `code:"10014" msg:"Đơn vị tiền tệ của hệ thống hiện tại chưa được đặt, vui lòng liên hệ với nhân viên để đặt đơn vị tiền tệ"`
	InsufficientBalance                 string `code:"10015" msg:"Không đủ số dư khả dụng"`
	Percentage                          string `code:"10016" msg:"Dữ liệu phần trăm / khối lượng không chính xác"`
	LimitPrice                          string `code:"10017" msg:"Thông số giá giới hạn không được để trống"`
	UserIsNotExist                      string `code:"10018" msg:"người dùng không tồn tại"`
	ContractIsNotExist                  string `code:"10019" msg:"Thông tin hợp đồng chưa được truy vấn, vui lòng liên hệ với quản trị viên để cấu hình thông tin hợp đồng của đơn vị tiền tệ hiện tại"`
	SysTradingPairIsExist               string `code:"10020" msg:"Không có cặp giao dịch nào được cấu hình trong hệ thống"`
	TradingPairIsNotExist               string `code:"10021" msg:"Cặp giao dịch không tồn tại"`
	withdrawalFeesIsNotExist            string `code:"10022" msg:"Hệ thống không được cấu hình với phí rút tiền, và việc rút tiền tạm thời không được hỗ trợ, vui lòng liên hệ với nhân viên"`
	UsersWalletIsNotExist               string `code:"10023" msg:"Giao dịch không tồn tại trong ví"`
	ApplyBuySetupIsNotExist             string `code:"10024" msg:"Đơn vị tiền tệ mua hàng không tồn tại"`
	PasswordEditError                   string `code:"10025" msg:"Sửa đổi mật khẩu đăng nhập không thành công"`
	OperationFailed                     string `code:"10026" msg:"Hành động không thành công"`
	FeeOptionContractIsError            string `code:"10027" msg:"Phí hợp đồng giao dịch quyền chọn chưa được thiết lập, vui lòng liên hệ cấu hình hoạt động tại đây"`
	WithdrawalFeesIsError               string `code:"10028" msg:"Phí rút tiền bất thường"`
	LoginFailed                         string `code:"10029" msg:"Dịch vụ ngoại lệ, đăng nhập không thành công!"`
	MultipleIsError                     string `code:"10030" msg:"Giá trị bội là bất hợp pháp, không phải giá trị do hệ thống đặt!"`
	ContractIsNotCorrect                string `code:"10031" msg:"Thông tin hợp đồng đặt chưa đúng, vui lòng liên hệ quản trị viên để đặt thông tin hợp đồng tương ứng cho cặp giao dịch"`
	UsersWalletLock                     string `code:"10032" msg:"Xin lỗi, ví của bạn tạm thời bị khóa"`
	LiquidationUnsuccessful             string `code:"10033" msg:"Không thành công"`
	OptionContractSecondsNotCorrect     string `code:"10034" msg:"Số giây không chính xác trong hợp đồng quyền chọn"`
	OptionContractStatus                string `code:"10035" msg:"Hợp đồng tùy chọn này chưa được kích hoạt!"`
	FeePerpetualContractIsError         string `code:"10036" msg:"Phí hợp đồng vĩnh viễn chưa được thiết lập, vui lòng liên hệ với cấu hình hoạt động!"`
	MinAmountIsNotExist                 string `code:"10037" msg:"Không có số tiền rút tối thiểu đã được thiết lập, vui lòng liên hệ với Cấu hình hoạt động!"`
	OptionContractProfitRatioNotCorrect string `code:"10038" msg:"Tỷ lệ hoàn vốn của hợp đồng quyền chọn không chính xác"`
	ApplyBuySetupStatusIsNotExist       string `code:"10039" msg:"Đơn vị tiền tệ mua vẫn chưa mở"`
	CurrencyTransactionIsExist          string `code:"10040" msg:"Hệ thống tạm thời không cho phép thay đổi giao dịch bằng đơn vị tiền tệ"`
	LimitPriceErr                       string `code:"10041" msg:"Lỗi thông số giá giới hạn"`
	EntrustNumErr                       string `code:"10042" msg:"Đặt hàng, thông số số lượng sai"`
	CurrencyTypeIsNotAllowed            string `code:"10043" msg:"Loại giao dịch này tạm thời không được phép cho cặp giao dịch này, vui lòng liên hệ với bộ phận vận hành nếu cần thiết"`
	UserIsLock                          string `code:"10044" msg:"Trạng thái người dùng bị khóa"`
	SearchTimeErr                       string `code:"10045" msg:"Định dạng thời gian tìm kiếm không chính xác eg: 2006-01-02"`
	OptionContractMinimum               string `code:"10046" msg:"Lượng giao dịch hiện tại thấp hơn lượng tiêu thụ nhiều nhất"`
}
