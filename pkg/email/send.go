package email

import (
	"goapi/pkg/config"
	"goapi/pkg/email/google"
	"goapi/pkg/email/qq"
)

func SendEmail(title string, text string, toId string) error {
	var err error
	emailType := config.GetString("email.type")
	if emailType == "QQ" {
		err = qq.SendEmail(title, config.GetString("email.qq.user"), toId, text)
		if err != nil {
			return err
		}
	}
	if emailType == "GOOGLE" {
		err = google.New().Send(title, text, toId)
		if err != nil {
			return err
		}
	}
	return nil
}
