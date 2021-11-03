package user_wallet

import (
	"goapi/app/models"
	"goapi/app/response"
	"goapi/pkg/mysql"
)

// GetBalance 币种余额
func GetBalance(userId, CurrencyId string) response.UsersWallet {
	where := cmap.New().Items()
	where["user_id"] = userId
	where["currency"] = CurrencyId
	var result response.UsersWallet
	DB := mysql.DB.Debug()
	DB.Model(models.UsersWallet{}).Where(where).Find(&result)
	return result
}
