package ti

// 泰语

type ErrorCode struct {
	ParsingError                        string `code:"10000" msg:"ข้อผิดพลาดในการแยกวิเคราะห์ข้อมูล"`
	LoginError                          string `code:"10001" msg:"การเข้าสู่ระบบล้มเหลว หรือโทเค็นของคุณหมดอายุ"`
	PwError                             string `code:"10002" msg:"รหัสผ่านผิด"`
	HtmlParsingError                    string `code:"10003" msg:"การแยกวิเคราะห์ HTML ล้มเหลว"`
	SendEmail                           string `code:"10004" msg:"ส่งอีเมลไม่สำเร็จ"`
	ValidatorError                      string `code:"10005" msg:"โปรดตรวจสอบพารามิเตอร์ ข้อผิดพลาดในการส่งพารามิเตอร์"`
	AddError                            string `code:"10006" msg:"เพิ่มข้อมูลไม่สำเร็จ"`
	UserIsExist                         string `code:"10007" msg:"ผู้ใช้รายนี้ลงทะเบียนแล้ว"`
	ShareCodeIsExist                    string `code:"10008" msg:"รหัสเชิญไม่มีอยู่"`
	VerCodeErr                          string `code:"10009" msg:"รหัสยืนยันอีเมลไม่ถูกต้อง"`
	ResetPassword                       string `code:"10010" msg:"เรียกรหัสผ่านไม่สำเร็จ"`
	LangSetUp                           string `code:"10011" msg:"การตั้งค่าภาษาล้มเหลว"`
	PayPasswordSetup                    string `code:"10012" msg:"ตั้งรหัสผ่านการชำระเงินไม่สำเร็จ"`
	CurrencyIsExist                     string `code:"10013" msg:"สกุลเงินไม่มีอยู่"`
	SysCurrencyIsExist                  string `code:"10014" msg:"สกุลเงินของระบบปัจจุบันไม่ได้ตั้งค่าไว้ โปรดติดต่อเจ้าหน้าที่เพื่อตั้งค่าสกุลเงิน"`
	InsufficientBalance                 string `code:"10015" msg:"ยอดเงินคงเหลือไม่เพียงพอ"`
	Percentage                          string `code:"10016" msg:"ข้อมูลเปอร์เซ็นต์/ระดับเสียงไม่ถูกต้อง"`
	LimitPrice                          string `code:"10017" msg:"พารามิเตอร์ราคาจำกัดไม่สามารถเว้นว่างได้"`
	UserIsNotExist                      string `code:"10018" msg:"ผู้ใช้ไม่มีอยู่"`
	ContractIsNotExist                  string `code:"10019" msg:"ยังไม่มีการสอบถามข้อมูลสัญญา โปรดติดต่อผู้ดูแลระบบเพื่อกำหนดค่าข้อมูลสัญญาของสกุลเงินปัจจุบัน"`
	SysTradingPairIsExist               string `code:"10020" msg:"ไม่มีคู่การซื้อขายที่กำหนดค่าไว้ในระบบ"`
	TradingPairIsNotExist               string `code:"10021" msg:"คู่ซื้อขายไม่มีอยู่จริง"`
	withdrawalFeesIsNotExist            string `code:"10022" msg:"ระบบไม่ได้กำหนดค่าให้มีค่าธรรมเนียมการถอนและไม่รองรับการถอนชั่วคราว โปรดติดต่อเจ้าหน้าที่"`
	UsersWalletIsNotExist               string `code:"10023" msg:"ไม่มีธุรกรรมอยู่ในกระเป๋าเงิน"`
	ApplyBuySetupIsNotExist             string `code:"10024" msg:"สกุลเงินที่ซื้อไม่มีอยู่"`
	PasswordEditError                   string `code:"10025" msg:"รหัสผ่านเข้าสู่ระบบ แก้ไขไม่สำเร็จ"`
	OperationFailed                     string `code:"10026" msg:"การดำเนินการล้มเหลว"`
	FeeOptionContractIsError            string `code:"10027" msg:"ยังไม่ได้ตั้งค่าค่าธรรมเนียมสัญญาซื้อขายตัวเลือก โปรดติดต่อการกำหนดค่าการดำเนินการที่นี่"`
	WithdrawalFeesIsError               string `code:"10028" msg:"ค่าธรรมเนียมการถอนที่ผิดปกติ"`
	LoginFailed                         string `code:"10029" msg:"ข้อยกเว้นของบริการ การเข้าสู่ระบบล้มเหลว!"`
	MultipleIsError                     string `code:"10030" msg:"ค่าหลายค่าไม่ถูกต้อง ไม่ใช่ค่าที่ระบบกำหนด!"`
	ContractIsNotCorrect                string `code:"10031" msg:"ตั้งค่าข้อมูลสัญญาไม่ถูกต้อง โปรดติดต่อผู้ดูแลระบบเพื่อตั้งค่าข้อมูลสัญญาที่เกี่ยวข้องสำหรับคู่ธุรกรรม"`
	UsersWalletLock                     string `code:"10032" msg:"ขออภัย กระเป๋าเงินของคุณถูกล็อคชั่วคราว"`
	LiquidationUnsuccessful             string `code:"10033" msg:"ตำแหน่งปิด ล้มเหลว"`
	OptionContractSecondsNotCorrect     string `code:"10034" msg:"จำนวนวินาทีที่ไม่ถูกต้องในสัญญาออปชั่น"`
	OptionContractStatus                string `code:"10035" msg:"สัญญาตัวเลือกนี้ยังไม่ได้เปิดใช้งาน!"`
	FeePerpetualContractIsError         string `code:"10036" msg:"ยังไม่ได้กำหนดค่าธรรมเนียมสัญญาถาวร โปรดติดต่อการกำหนดค่าการดำเนินการ!"`
	MinAmountIsNotExist                 string `code:"10037" msg:"ไม่มีการกำหนดจำนวนการถอนขั้นต่ำ โปรดติดต่อ Operation Configuration!"`
	OptionContractProfitRatioNotCorrect string `code:"10038" msg:"อัตราผลตอบแทนสัญญาออปชั่นไม่ถูกต้อง"`
	ApplyBuySetupStatusIsNotExist       string `code:"10039" msg:"สกุลเงินที่ซื้อยังไม่เปิด"`
	CurrencyTransactionIsExist          string `code:"10040" msg:"ระบบไม่อนุญาตให้เปลี่ยนธุรกรรมในสกุลเงินชั่วคราว"`
	LimitPriceErr                       string `code:"10041" msg:"ข้อผิดพลาดพารามิเตอร์ราคาจำกัด"`
	EntrustNumErr                       string `code:"10042" msg:"จำนวนการสั่งซื้อ ข้อผิดพลาดของพารามิเตอร์"`
	CurrencyTypeIsNotAllowed            string `code:"10043" msg:"ธุรกรรมประเภทนี้ไม่ได้รับอนุญาตชั่วคราวสำหรับคู่ธุรกรรมนี้ โปรดติดต่อการดำเนินการในกรณีที่จำเป็น"`
	UserIsLock                          string `code:"10044" msg:"สถานะผู้ใช้ถูกล็อค"`
	SearchTimeErr                       string `code:"10045" msg:"รูปแบบเวลาค้นหาไม่ถูกต้อง eg: 2006-01-02"`
	OptionContractMinimum               string `code:"10046" msg:"ปริมาณธุรกรรมปัจจุบันต่ำกว่าปริมาณการใช้มากที่สุด"`
}
