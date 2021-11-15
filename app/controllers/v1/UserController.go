package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goapi/app/models"
	"goapi/app/requests"
	"goapi/app/response"
	"goapi/pkg/echo"
	"goapi/pkg/helpers"
	"goapi/pkg/mysql"
	"goapi/pkg/redis"
	"goapi/pkg/share_code"
	"goapi/pkg/validator"
)

type UserController struct {
	BaseController
}

// UserInfoHandler 获取用户信息
func (h *UserController) UserInfoHandler(c *gin.Context) {
	userId, _ := c.Get("user_id")
	result := new(response.User)
	mysql.DB.Debug().Model(models.User{}).Where("id", userId).Find(&result)
	// 检测邀请码信息，初始化一系列
	share_code.GetShareCode(result)
	echo.Success(c, result, "", "")
}

func (h *UserController) LangSetupHandler(c *gin.Context) { // 初始化数据模型结构体
	userId, _ := c.Get("user_id")
	var (
		params requests.Lang
		user   models.User
	)
	// 绑定接收的 json 数据到结构体中
	_ = c.Bind(&params)
	// 数据验证
	vErr := validator.Validate.Struct(params)
	if vErr != nil {
		msg := validator.Lang(c.Request.Header.Get("Language")).Translate(vErr, params, c.Request.Header.Get("Language"))
		echo.Error(c, "ValidatorError", msg)
		return
	}
	c.Request.Header.Set("Language", params.Language)
	switch params.Language {
	case "zh":
		params.Language = "1"
		break
	case "en":
		params.Language = "2"
		break
	case "jp":
		params.Language = "3"
	}
	// 开始事务
	DB := mysql.DB.Debug().Begin()
	user.Language = params.Language
	err := DB.Model(&models.User{}).Where(map[string]interface{}{"id": fmt.Sprintf("%v", userId)}).Updates(user).Error
	if err != nil {
		// 遇到错误时回滚事务
		DB.Rollback()
		// todo 日志记录
		echo.Error(c, "LangSetUp", "")
		return
	}
	DB.Commit()
	echo.Success(c, c.Request.Header.Get("Language"), "", "")
}

func (h *UserController) PayPasswordSetupHandler(c *gin.Context) {
	var (
		params requests.PayPassword
		user   models.User
	)
	// 绑定接收的 json 数据到结构体中
	_ = c.Bind(&params)
	//// 数据验证
	vErr := validator.Validate.Struct(params)
	if vErr != nil {
		msg := validator.Lang(c.Request.Header.Get("Language")).Translate(vErr, params, c.Request.Header.Get("Language"))
		echo.Error(c, "ValidatorError", msg)
		return
	}
	userInfo, _ := c.Get("user")
	if params.Email != userInfo.(map[string]interface{})["email"].(string) {
		echo.Error(c, "VerCodeErr", "")
		return
	}

	// 开始事务
	DB := mysql.DB.Debug().Begin()
	user.PayPassword = helpers.Md5(params.PayPassword)
	err := DB.Model(&models.User{}).Where(map[string]interface{}{"email": params.Email}).Updates(user).Error
	if err != nil {
		// 遇到错误时回滚事务
		DB.Rollback()
		echo.Error(c, "PayPasswordSetup", "")
		return
	}
	DB.Commit()
	echo.Success(c, "ok", "", "")
}

// EditPasswordHandler 修改密码
func (h *UserController) EditPasswordHandler(c *gin.Context) {
	var (
		params requests.Password // 接收请求参数
		user   models.User       //  用户模型
	)
	// 绑定接收的 json 数据到结构体中
	_ = c.Bind(&params)
	//// 数据验证
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
	userInfo, _ := c.Get("user")
	if params.Email != userInfo.(map[string]interface{})["email"].(string) {
		echo.Error(c, "VerCodeErr", "")
		return
	}
	// 开始事务
	DB := mysql.DB.Debug().Begin()
	user.Password = helpers.Md5(params.Password)
	err := DB.Model(&models.User{}).Where(map[string]interface{}{"email": params.Email}).Updates(user).Error
	if err != nil {
		fmt.Println(err.Error())
		// 遇到错误时回滚事务
		DB.Rollback()
		echo.Error(c, "PasswordEditError", "")
		return
	}
	DB.Commit()
	echo.Success(c, "ok", "", "")
}

func (h *UserController) ShareCodeHandler(c *gin.Context) {
	userId, _ := c.Get("user_id")
	result := new(response.User)
	mysql.DB.Debug().Model(models.User{}).Where("id", userId).Find(&result)

	echo.Success(c, result, "", "")
}
