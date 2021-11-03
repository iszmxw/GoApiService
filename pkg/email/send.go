package email

import (
	"bytes"
	sendEmail "github.com/jordan-wright/email"
	"goapi/pkg/config"
	"html/template"
	"net/smtp"
)

// SendEmail 使用第三方库发送邮件
func SendEmail(Subject, fromUser, toUser, Code string) error {
	e := sendEmail.NewEmail()

	e.From = fromUser
	e.To = []string{toUser}
	e.Subject = Subject

	t, err := template.ParseFiles("pkg/email/send.html")
	if err != nil {
		return err
	}
	body := new(bytes.Buffer)
	//作为变量传递给html模板
	err1 := t.Execute(body, struct {
		Email      string
		ActiveCode string
	}{
		Email:      toUser,
		ActiveCode: Code,
	})
	if err1 != nil {
		return err1
	}
	// html形式的消息
	e.HTML = body.Bytes()
	return e.Send("smtp.qq.com:587", smtp.PlainAuth(
		"",
		config.GetString("email.user"),
		config.GetString("email.password"),
		"smtp.qq.com",
	))
}
