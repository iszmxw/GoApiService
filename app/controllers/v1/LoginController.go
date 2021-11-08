package v1

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"goapi/app/models"
	"goapi/app/requests"
	"goapi/app/response"
	"goapi/pkg/config"
	"goapi/pkg/echo"
	"goapi/pkg/email"
	"goapi/pkg/helpers"
	"goapi/pkg/logger"
	"goapi/pkg/mysql"
	"goapi/pkg/redis"
	"goapi/pkg/request"
	"goapi/pkg/share_code"
	"goapi/pkg/validator"
	"strconv"
	"time"
)

type LoginController struct {
	BaseController
}

// LoginHandler 登录接口
// @Summary 登录接口
// @Description 提交注册的邮箱和密码即可登录
// @Tags 登录接口
// @Accept multipart/form-data
// @Produce application/json
// @Param param formData requests.UserLogin false "请求参数"
// @Success 200 {object} response._LoginHandler
// @Router /v1/api/user/login [post]
func (h *LoginController) LoginHandler(c *gin.Context) {
	// 初始化数据模型结构体
	var params requests.UserLogin
	// 绑定接收的 json 数据到结构体中
	_ = c.Bind(&params)
	// 数据验证
	vErr := validator.Validate.Struct(params)
	if vErr != nil {
		msg := validator.Lang(c.Request.Header.Get("Language")).Translate(vErr, params, c.Request.Header.Get("Language"))
		echo.Error(c, "", msg)
		return
	}

	db := new(models.User)
	user := new(response.User)
	db.GetOne(map[string]interface{}{"email": params.Email}, user)
	if user.Id <= 0 {
		echo.Error(c, "UserIsNotExist", "")
		return
	}
	endMd5 := helpers.Md5(params.Password)
	if endMd5 != user.Password {
		echo.Error(c, "PwError", "")
		return
	}
	// 检测邀请码信息，初始化一系列
	share_code.GetShareCode(user)
	info, errs := json.Marshal(user) //转换成JSON返回的是byte[]
	if errs != nil {
		fmt.Println(errs.Error())
		echo.Error(c, "ParsingError", "")
		return
	}
	token := helpers.GetUUID()
	loginKey := "login:user:" + helpers.IntToString(user.Id)
	oldToken, _ := redis.Get(loginKey)
	if len(oldToken) > 0 {
		redis.Delete(oldToken)
	}
	_, err := redis.Add(token, string(info), 60*60*2) // 缓存两个小时过期
	_, err1 := redis.Add(loginKey, token, 60*60*2)    // 缓存两个小时过期
	if err != nil || err1 != nil {
		logger.Info("服务异常，登录失败！")
		logger.Error(err)
		logger.Error(err1)
		echo.Error(c, "LoginFailed", "")
		return
	}
	// 返回结果
	result := gin.H{
		"token": token,
		"uid":   user.Id,
	}
	echo.Success(c, result, "", "")
}

// SendEmailRegisterHandler 发送注册邮件
// @Summary 发送注册邮件
// @Tags 发送注册邮件
// @Description 发送注册邮件
// @Accept multipart/form-data
// @Produce application/json
// @Param param formData requests.UserEmail false "请求参数"
// @Success 200 {object} response._OK
// @Router /v1/api/user/send_email_register [post]
func (h *LoginController) SendEmailRegisterHandler(c *gin.Context) {
	// 初始化数据模型结构体
	var params requests.SendEmailRegister
	// 绑定接收的 json 数据到结构体中
	_ = c.Bind(&params)
	// 数据验证
	vErr := validator.Validate.Struct(params)
	if vErr != nil {
		msg := validator.Lang(c.Request.Header.Get("Language")).Translate(vErr, params, c.Request.Header.Get("Language"))
		echo.Error(c, "", msg)
		return
	}
	// 检查用户是否注册
	UserIsExist := new(response.User)
	models.User{}.GetOne(map[string]interface{}{"email": params.Email}, UserIsExist)
	if UserIsExist.Id > 0 {
		echo.Error(c, "UserIsExist", "")
		return
	}
	Code := helpers.Rand(6)
	key := "RegisterCode:" + params.Email
	_, err1 := redis.Add(key, Code, 60*60*24)
	if err1 != nil {
		logger.Error(err1)
		echo.Error(c, "SendEmail", "")
		return
	}
	err := email.SendEmail("Register Send Code", Code, params.Email)
	if err != nil {
		logger.Error(err)
		echo.Error(c, "SendEmail", err.Error())
		return
	}
	echo.Success(c, "ok", "", "")
}

// VerifyRegisterHandler 验证注册
func (h *LoginController) VerifyRegisterHandler(c *gin.Context) {
	// 初始化数据模型结构体
	var params requests.UserRegister
	parentUser := new(response.User)
	_ = c.Bind(&params)
	// 数据验证
	vErr := validator.Validate.Struct(params)
	if vErr != nil {
		msg := validator.Lang(c.Request.Header.Get("Language")).Translate(vErr, params, c.Request.Header.Get("Language"))
		echo.Error(c, "ValidatorError", msg)
		return
	}
	// 检测验邮箱证码
	Code, codeErr := redis.Get("RegisterCode:" + params.Email)
	if codeErr != nil {
		echo.Error(c, "VerCodeErr", "")
		return
	}
	if Code != params.Code {
		echo.Error(c, "VerCodeErr", "")
		return
	}

	// 检查用户是否注册
	UserIsExist := new(response.User)
	models.User{}.GetOne(map[string]interface{}{"email": params.Email}, UserIsExist)
	if UserIsExist.Id > 0 {
		echo.Error(c, "UserIsExist", "")
		return
	}

	var GlobalsTypes response.GlobalsTypes
	mysql.DB.Model(models.Globals{}).Where(map[string]interface{}{"fields": "invitation_code"}).Find(&GlobalsTypes)
	if GlobalsTypes.Id <= 0 {
		logger.Info("邀请码开关未设置")
	}

	// 开启了邀请码
	if GlobalsTypes.Value == "1" {
		sysShareCode := config.Env("SHARE_CODE").(string)
		fmt.Println(fmt.Sprintf("sysShareCode: %v,ShareCode: %v", sysShareCode, params.ShareCode))
		// 邀请码是否存在
		if sysShareCode != params.ShareCode {
			models.User{}.GetOne(map[string]interface{}{"share_code": params.ShareCode}, parentUser)
			if parentUser.Id <= 0 {
				echo.Error(c, "ShareCodeIsExist", "")
				return
			}
		}
	}

	timeNow := time.Now()
	user := models.User{}
	user.Language = "1" // 默认繁体
	user.IsAgent = "0"  // 默认为零，暂时不是代理
	user.Email = params.Email
	user.Nickname = params.Email
	user.UserLevel = 1 // 默认为1
	user.Password = helpers.Md5(params.Password)
	user.LoginTime = timeNow
	user.LastLoginIp = c.ClientIP()
	user.Status = "0" // 正常
	user.AgentDividend = "0"

	// 开始事务
	DB := mysql.DB.Debug().Begin()
	err := DB.Model(&user).Create(&user).Error
	if err != nil {
		// 遇到错误时回滚事务
		DB.Rollback()
		// todo 日志记录
		echo.Error(c, "LangSetUp", "")
		return
	}

	// 上级信息
	if parentUser.Id > 0 {
		if parentUser.UserLevel == 1 {
			user.UserPath = strconv.Itoa(parentUser.Id) + "," + strconv.Itoa(user.Id)
		} else {
			user.UserPath = parentUser.UserPath + "," + strconv.Itoa(user.Id)
		}
		user.ParentId = parentUser.Id
		user.UserLevel = parentUser.UserLevel + 1
	}

	err = DB.Model(&user).Updates(&user).Error
	if err != nil {
		// 遇到错误时回滚事务
		DB.Rollback()
		echo.Error(c, "AddError", "")
		return
	}

	// 检测邀请码信息，初始化一系列
	UserCode := new(response.User)
	DB.Model(models.User{}).Where(map[string]interface{}{"email": params.Email}).Find(&UserCode)
	if len(UserCode.ShareCode) <= 0 {
		var initShareCode int
		// 邀请码不存在，初始化邀请码
		maxShareCode := new(response.User)
		// 查找系统最大的邀请码
		DB.Model(models.User{}).Where("share_code <> ?", "").Order("share_code desc").Find(maxShareCode)
		if len(maxShareCode.ShareCode) <= 0 {
			init, err := strconv.Atoi(config.Env("INIT_SHARE").(string))
			if err != nil {
				fmt.Println("邀请码初始值获取失败")
			}
			initShareCode = init + UserCode.Id
		} else {
			autoShareCode, err := strconv.Atoi(maxShareCode.ShareCode)
			if err != nil {
				fmt.Println("获取最大邀邀请码失败")
			}
			initShareCode = autoShareCode + 1
		}
		err := DB.Model(models.User{}).Where(map[string]interface{}{"id": UserCode.Id}).Update("share_code", initShareCode).Error
		if err != nil {
			DB.Rollback()
			// todo 日志记录
			fmt.Println("邀请码初始化失败")
			//panic("邀请码初始化失败")
		}
	}
	// 检测邀请码信息，初始化一系列

	var TradingPair []response.TradingPair
	DB.Model(models.TradingPair{}).Find(&TradingPair)
	// 初始化钱包数据
	var UsersWalletArr []models.UsersWallet
	var UsersWallet models.UsersWallet
	Type := map[int]int{1: 1, 2: 2}
	if len(TradingPair) <= 0 {
		DB.Rollback()
		echo.Error(c, "SysTradingPairIsExist", "")
		return
	}
	for _, value := range TradingPair {
		for _, v := range Type {
			UsersWallet.UserId = user.Id
			UsersWallet.Type = v                     // 钱包类型：1现货 2合约
			UsersWallet.TradingPairId = value.Id     // 有多少种交易对就有多少种钱包
			UsersWallet.TradingPairName = value.Name // 交易对名称
			UsersWallet.Address = helpers.GetUUID()
			UsersWallet.Status = 0      // 0正常 1锁定
			UsersWallet.Available = 0   // 可用
			UsersWallet.Freeze = 0      // 冻结
			UsersWallet.TotalAssets = 0 // 总资产
			UsersWalletArr = append(UsersWalletArr, UsersWallet)
		}
	}
	// 生成用戶钱包
	InitUsersWallet := DB.Model(models.UsersWallet{}).Create(&UsersWalletArr).Error
	if InitUsersWallet != nil {
		fmt.Println(InitUsersWallet.Error())
		// 遇到错误时回滚事务
		DB.Rollback()
		echo.Error(c, "AddError", "")
		return
	}

	// 否则，提交事务
	DB.Commit()
	echo.Success(c, gin.H{
		"id":       user.Id,
		"language": user.Language,
		"email":    user.Email,
	}, "", "")
}

// SendEmailRetrieveHandler 找回密码邮件
func (h *LoginController) SendEmailRetrieveHandler(c *gin.Context) {
	// 初始化数据模型结构体
	var params requests.UserEmail
	// 绑定接收的 json 数据到结构体中
	_ = c.Bind(&params)
	// 数据验证
	vErr := validator.Validate.Struct(params)
	if vErr != nil {
		msg := validator.Lang(c.Request.Header.Get("Language")).Translate(vErr, params, c.Request.Header.Get("Language"))
		echo.Error(c, "", msg)
		return
	}

	// 发送找回邮件前检查用户是否注册
	UserIsExist := new(response.User)
	models.User{}.GetOne(map[string]interface{}{"email": params.Email}, UserIsExist)
	if UserIsExist.Id <= 0 {
		echo.Error(c, "UserIsExist", "")
		return
	}

	Code := helpers.Rand(6)
	key := "RetrieveCode:" + params.Email
	_, err1 := redis.Add(key, Code, 60*60*24)
	if err1 != nil {
		echo.Error(c, "SendEmail", "")
		return
	}
	err := email.SendEmail("Retrieve Password Send Code", Code, params.Email)
	if err != nil {
		echo.Error(c, "SendEmail", err.Error())
		return
	}
	echo.Success(c, "ok", "", "")
}

// ResetVerifyHandler 重置密码前验证邮件Code
func (h *LoginController) ResetVerifyHandler(c *gin.Context) {
	// 初始化数据模型结构体
	var (
		params requests.ResetVerify
	)
	// 绑定接收的 json 数据到结构体中
	_ = c.Bind(&params)
	//// 数据验证
	vErr := validator.Validate.Struct(params)
	if vErr != nil {
		msg := validator.Lang(c.Request.Header.Get("Language")).Translate(vErr, params, c.Request.Header.Get("Language"))
		echo.Error(c, "", msg)
		return
	}
	// 检测验邮箱证码
	Code, codeErr := redis.Get("RetrieveCode:" + params.Email)
	if codeErr != nil {
		echo.Error(c, "VerCodeErr", "")
		return
	}
	if Code != params.Code {
		echo.Error(c, "VerCodeErr", "")
		return
	}
	echo.Success(c, "ok", "", "")
}

// ResetPasswordHandler 重置密码
func (h *LoginController) ResetPasswordHandler(c *gin.Context) {

	// 初始化数据模型结构体
	var (
		params requests.ResetPassword
	)
	// 绑定接收的 json 数据到结构体中
	_ = c.Bind(&params)
	//// 数据验证
	vErr := validator.Validate.Struct(params)
	if vErr != nil {
		msg := validator.Lang(c.Request.Header.Get("Language")).Translate(vErr, params, c.Request.Header.Get("Language"))
		echo.Error(c, "", msg)
		return
	}
	// 检测验邮箱证码
	Code, codeErr := redis.Get("RetrieveCode:" + params.Email)
	if codeErr != nil {
		echo.Error(c, "VerCodeErr", "")
		return
	}
	if Code != params.Code {
		echo.Error(c, "VerCodeErr", "")
		return
	}
	// 检查用户是否注册
	UserIsExist := new(response.User)
	models.User{}.GetOne(map[string]interface{}{"email": params.Email}, UserIsExist)
	if UserIsExist.Id <= 0 {
		echo.Error(c, "UserIsExist", "")
		return
	}
	// 开始事务
	DB := mysql.DB.Debug().Begin()
	user := new(response.User)
	// 修改密码
	user.Password = helpers.Md5(params.Password)
	err := DB.Model(&models.User{}).Where(map[string]interface{}{"email": params.Email}).Updates(&user).Error
	if err != nil {
		// 遇到错误时回滚事务
		DB.Rollback()
		echo.Error(c, "ResetPassword", "")
		return
	}
	// 否则，提交事务
	DB.Commit()
	echo.Success(c, "ok", "", "")
}

// LogoutHandler 退出
func (h *LoginController) LogoutHandler(c *gin.Context) {
	tokenString := request.GetParam(c, "token")
	ok := redis.Delete(tokenString)
	fmt.Println("退出登录", ok)
	echo.Success(c, "", "ok", "")
}
