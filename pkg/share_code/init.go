package share_code

import (
	"fmt"
	"goapi/app/models"
	"goapi/app/response"
	"goapi/pkg/config"
	"goapi/pkg/mysql"
	"strconv"
)

func GetShareCode(currentUser *response.User) {
	if len(currentUser.ShareCode) <= 0 {
		var initShareCode int
		// 邀请码不存在，初始化邀请码
		maxShareCode := new(response.User)
		DB := mysql.DB.Debug().Begin()
		// 查找系统最大的邀请码
		DB.Model(models.User{}).Where("share_code <> ?", "").Order("share_code desc").Find(maxShareCode)
		if len(maxShareCode.ShareCode) <= 0 {
			init, err := strconv.Atoi(config.Env("INIT_SHARE").(string))
			if err != nil {
				fmt.Println("邀请码初始值获取失败")
			}
			initShareCode = init + currentUser.Id
		} else {
			autoShareCode, err := strconv.Atoi(maxShareCode.ShareCode)
			if err != nil {
				fmt.Println("获取最大邀邀请码失败")
			}
			initShareCode = autoShareCode + 1
		}
		err := DB.Model(models.User{}).Where(map[string]interface{}{"id": currentUser.Id}).Update("share_code", initShareCode).Error
		if err != nil {
			DB.Rollback()
			// todo 日志记录
			fmt.Println("邀请码初始化失败")
			//panic("邀请码初始化失败")
		}
		DB.Commit()
	}
}
