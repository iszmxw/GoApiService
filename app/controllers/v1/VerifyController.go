package v1

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"goapi/app/models"
	"goapi/app/requests"
	"goapi/app/response"
	conf "goapi/pkg/config"
	"goapi/pkg/echo"
	"goapi/pkg/logger"
	"goapi/pkg/mysql"
	"goapi/pkg/validator"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

type VerifyController struct {
	BaseController
}

var ErrorFileToBig = errors.New("每张图片大小不能超过5mb")

// VerifyPrimaryHandle 初级验证
func (h *VerifyController) VerifyPrimaryHandle(c *gin.Context) {
	//获取参数
	var params requests.VerifyParam
	_ = c.Bind(&params)
	//参数验证
	logger.Info(params)
	vErr := validator.Validate.Struct(params)
	if vErr != nil {
		msg := validator.Lang(c.Request.Header.Get("Language")).Translate(vErr, params, c.Request.Header.Get("Language"))
		echo.Error(c, "ValidatorError", msg)
		return
	}
	//逻辑
	//获取token中的用户ID
	userId, _ := c.Get("user_id")
	logger.Info(userId.(int))
	userInfo, _ := c.Get("user")
	Email := userInfo.(map[string]interface{})["email"].(string)
	v := models.Verify{
		UserId:       userId.(int),
		IdentityCard: params.IdentityCard,
		FullName:     params.FullName,
		Email:        Email,
		Status:       0, //前端提交后修改状态为0 后端审核后改为1
	}
	//v1 := models.Verify{}
	//不能修改认证信息
	//查询user_id是否存在
	DB := mysql.DB.Debug().Begin()
	var v1 models.Verify
	DB.Model(models.Verify{}).Where("user_id", userId).Find(&v1)
	//if v1.Status != 0 {
	//	echo.Error(c, "", "用户已存在数据库")
	//	return
	//}
	if v1.UserId != 0 {
		echo.Error(c, "", "已经提交初级验证信息")
		return
	}
	if v1.Status >= 1 {
		echo.Error(c, "", "用户已完成初级认证")
		return
	}
	CreateErr := DB.Model(models.Verify{}).Create(&v).Error
	if CreateErr != nil {
		echo.Error(c, "", "添加用户验证信息失败")
		DB.Rollback()
		return
	}
	DB.Commit()
	echo.Success(c, "ok", "添加初级验证已提交", "")

}

//VerifyAdvancedHandle 高级验证
func (h *VerifyController) VerifyAdvancedHandle(c *gin.Context) {
	////获取图片+参数验证
	//绑定参数
	var p requests.ImgBase64Param
	err := c.Bind(&p)
	if err != nil {
		logger.Error(err)
		return
	}
	//验证参数
	vErr := validator.Validate.Struct(p)
	if vErr != nil {
		msg := validator.Lang(c.Request.Header.Get("Language")).Translate(vErr, p, c.Request.Header.Get("Language"))
		echo.Error(c, "ValidatorError", msg)
		return
	}

	//逻辑处理
	//1.判断是否通过初级验证
	//逻辑
	DB := mysql.DB.Debug().Begin()
	userId, _ := c.Get("user_id")
	var v1 models.Verify
	err = DB.Raw("select * from osx_user_img where user_id =?", userId.(int)).Scan(&v1).Error
	if err != nil {
		logger.Error(err)
		echo.Error(c, "", "查询用户出错")
		return
	}
	if v1.Status == 2 {
		echo.Error(c, "", "用户已完成验证")
		return
	}
	if v1.Status == 0 {
		echo.Error(c, "", "初级验证未通过")
		return
	}
	if v1.Status == -1 {
		filedir := fmt.Sprintf("./resource/photo/%d/", userId)
		Rerr := os.RemoveAll(filedir)
		if Rerr != nil {
			logger.Error(Rerr)
			echo.Error(c, "", "删除图片失败")
			return
		}
		logger.Info("删除失败认证的图片")
	}
	//把四个图片的base64存进一个切片
	imgMap := make(map[string]string, 4)
	imgMap["img_card_front"] = p.ImgCardFront
	imgMap["img_card_behind"] = p.ImgCardBehind
	imgMap["img_bank_front"] = p.ImgBankFront
	imgMap["img_bank_behind"] = p.ImgBankBehind
	//
	filedir := fmt.Sprintf("./resource/photo/%d/", userId)
	Eerr := os.Mkdir(filedir, 0755)
	if Eerr != nil {
		logger.Error(Eerr)
		echo.Error(c, "", "创建文件夹失败")
		return
	}
	for key, value := range imgMap {
		filename := fmt.Sprintf("%s.jpg", key)
		transErr := base64transfer(value, filedir, filename)
		if transErr != nil {
			logger.Error(transErr)
			echo.Error(c, "", "图片转换储存失败")
			return
		}

	}
	//获取地址
	ip := conf.GetString("APP_URL")
	//获取端口
	port := conf.GetString("app.port")
	//api url
	apiUrl := "/v1/api/verify/downloadImg?imgUrl="
	//http://127.0.0.1:80/v1/api/verify/downloadImg?imgUrl=
	//把存进图片的url写进mysql
	v1.ImgCardFront = ip + ":" + port + apiUrl + filedir + "img_card_front.jpg"
	v1.ImgCardBehind = ip + ":" + port + apiUrl + filedir + "img_card_behind.jpg"
	v1.ImgBankFront = ip + ":" + port + apiUrl + filedir + "img_bank_front.jpg"
	v1.ImgBankBehind = ip + ":" + port + apiUrl + filedir + "img_bank_behind.jpg"
	v1.Status = 1
	sErr := DB.Model(models.Verify{}).Where("user_id", userId).Save(&v1).Error
	if sErr != nil {
		logger.Error(err)
		echo.Error(c, "", "更新验证信息失败")
		DB.Rollback()
		return
	}
	DB.Commit()
	echo.Success(c, "ok", "添加高级验证已提交", "")
}

// VerifyDownloadHandle 下载图片
func (h *VerifyController) VerifyDownloadHandle(c *gin.Context) {
	//获取参数
	imgUrl := c.Query("imgUrl")
	logger.Info(imgUrl)
	if imgUrl == "" {
		c.JSON(http.StatusOK, gin.H{
			"msg": "输入参数不能为空",
		})
		return
	}
	fileTmp, errByOpenFile := os.Open(imgUrl)
	if errByOpenFile != nil {
		logger.Error(errByOpenFile)
	}
	defer fileTmp.Close()
	//逻辑
	//c.Header("Content-Type", "application/octet-stream")
	//强制浏览器下载
	//c.Header("Content-Disposition", "attachment; filename="+imgUrl)
	//浏览器下载或预览
	c.Header("Content-Disposition", "inline;filename="+imgUrl)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")
	file, _ := ioutil.ReadFile(imgUrl)
	_, err := c.Writer.WriteString(string(file))
	if err != nil {
		logger.Error(err)
		echo.Error(c, "", "输出图片失败")
		return
	}
	//c.File(imgUrl)
	return
	//返回
}

// UserVerifyStatusHandle 返回用户验证信息
func (h *VerifyController) UserVerifyStatusHandle(c *gin.Context) {
	userId, _ := c.Get("user_id")
	DB := mysql.DB.Debug()
	var Verify response.Verify
	cErr := DB.Model(models.Verify{}).Where("user_id", userId).Find(&Verify).Error
	if cErr != nil {
		logger.Error(cErr)
		echo.Error(c, "", "用户未验证")
		return
	}
	echo.Success(c, Verify, "", "")
}

// PathExists 判断文件路径是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// vfImg 判断图片的格式
func vfImg(filename string) (err error) {
	fileExt := strings.ToLower(path.Ext(filename))
	if fileExt != ".png" && fileExt != ".jpg" && fileExt != ".gif" && fileExt != ".jpeg" {
		str := "上传失败!只允许png,jpg,gif,jpeg文件"
		return errors.New(str)
	}
	return nil
}

/**
 * @Author
 * @Description 通过base64 生成文件且判断文件大小
 * @Date
 * @Param base64 图片base64
 * @Param filepath 路径
 * @return
 **/
func base64transfer(base64str string, filedir string, filename string) error {
	filepath := filedir + filename
	imgbff, _ := base64.StdEncoding.DecodeString(base64str)
	if len(imgbff) > 5120000 {
		return ErrorFileToBig
	}
	err := ioutil.WriteFile(filepath, imgbff, 0666)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
