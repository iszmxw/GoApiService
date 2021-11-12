package es

// 西班牙语言

type ErrorCode struct {
	ParsingError                        string `code:"10000" msg:"Error de análisis de datos"`
	LoginError                          string `code:"10001" msg:"El inicio de sesión falló, o su token ha caducado"`
	PwError                             string `code:"10002" msg:"contraseña incorrecta"`
	HtmlParsingError                    string `code:"10003" msg:"Análisis HTML falló"`
	SendEmail                           string `code:"10004" msg:"Entrega de correo fallida"`
	ValidatorError                      string `code:"10005" msg:"Por favor, compruebe los parámetros, el error de transmisión de parámetros"`
	AddError                            string `code:"10006" msg:"Añadiendo datos fallidos"`
	UserIsExist                         string `code:"10007" msg:"Este usuario ha sido registrado."`
	ShareCodeIsExist                    string `code:"10008" msg:"El código de invitación no existe"`
	VerCodeErr                          string `code:"10009" msg:"El código de verificación del buzón es incorrecto"`
	ResetPassword                       string `code:"10010" msg:"Recuperar la falla de la contraseña"`
	LangSetUp                           string `code:"10011" msg:"Falló el ajuste de idioma"`
	PayPasswordSetup                    string `code:"10012" msg:"Falló la configuración de la contraseña de pago"`
	CurrencyIsExist                     string `code:"10013" msg:"La moneda no existe"`
	SysCurrencyIsExist                  string `code:"10014" msg:"La moneda del sistema actual no está establecida, comuníquese con el personal para establecer la moneda"`
	InsufficientBalance                 string `code:"10015" msg:"Saldo insuficiente"`
	Percentage                          string `code:"10016" msg:"Error de datos por ciento / transacción"`
	LimitPrice                          string `code:"10017" msg:"El parámetro de límite de precio no puede estar vacío"`
	UserIsNotExist                      string `code:"10018" msg:"el usuario no existe"`
	ContractIsNotExist                  string `code:"10019" msg:"Contiene información no se pregunta, comuníquese con el administrador para configurar la información del contrato de la moneda actual."`
	SysTradingPairIsExist               string `code:"10020" msg:"El sistema no está configurado"`
	TradingPairIsNotExist               string `code:"10021" msg:"No se existe el comercio"`
	withdrawalFeesIsNotExist            string `code:"10022" msg:"El sistema no está configurado con la tarifa de recuperación, y el efectivo no es compatible. Póngase en contacto con el personal."`
	UsersWalletIsNotExist               string `code:"10023" msg:"La transacción no existe en la billetera."`
	ApplyBuySetupIsNotExist             string `code:"10024" msg:"La aplicación de la moneda no existe."`
	PasswordEditError                   string `code:"10025" msg:"Falló la modificación de la contraseña de inicio de sesión fallida"`
	OperationFailed                     string `code:"10026" msg:"operación fallida"`
	FeeOptionContractIsError            string `code:"10027" msg:"No se ha establecido la tarifa del contrato de transacción de la opción, comuníquese con la configuración de la operación aquí"`
	WithdrawalFeesIsError               string `code:"10028" msg:"Tarifa de retiro anormal"`
	LoginFailed                         string `code:"10029" msg:"¡Excepción de servicio, error de inicio de sesión!"`
	MultipleIsError                     string `code:"10030" msg:"¡El valor múltiple es ilegal, no el valor establecido por el sistema!"`
	ContractIsNotCorrect                string `code:"10031" msg:"La información del contrato no está configurada correctamente, comuníquese con el administrador para configurar la información del contrato correspondiente para el par de transacciones"`
	UsersWalletLock                     string `code:"10032" msg:"Lo sentimos, tu billetera está bloqueada temporalmente"`
	LiquidationUnsuccessful             string `code:"10033" msg:"Fracasado"`
	OptionContractSecondsNotCorrect     string `code:"10034" msg:"Número incorrecto de segundos en el contrato de opciones"`
	OptionContractStatus                string `code:"10035" msg:"¡Este contrato de opción aún no se ha activado!"`
	FeePerpetualContractIsError         string `code:"10036" msg:"No se ha establecido la tarifa del contrato perpetuo, póngase en contacto con la configuración de la operación."`
	MinAmountIsNotExist                 string `code:"10037" msg:"No se ha establecido un monto mínimo de retiro, comuníquese con Configuración de operación."`
	OptionContractProfitRatioNotCorrect string `code:"10038" msg:"La tasa de retorno del contrato de opción es incorrecta"`
	ApplyBuySetupStatusIsNotExist       string `code:"10039" msg:"La moneda de suscripción aún no se ha abierto"`
	CurrencyTransactionIsExist          string `code:"10040" msg:"El sistema no permite temporalmente que se modifiquen las transacciones en la moneda."`
	LimitPriceErr                       string `code:"10041" msg:"Error de parámetro de precio límite"`
	EntrustNumErr                       string `code:"10042" msg:"Error de parámetro de importe de pedido"`
	CurrencyTypeIsNotAllowed            string `code:"10043" msg:"Este tipo de transacción no está permitido temporalmente para este par de transacciones, comuníquese con la operación si es necesario"`
	UserIsLock                          string `code:"10044" msg:"El estado del usuario está bloqueado"`
	SearchTimeErr                       string `code:"10045" msg:"El formato de la hora de búsqueda es incorrecto, por ejemplo: 2006-01-02"`
}
