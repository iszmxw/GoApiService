package gmail

import (
	"crypto/tls"
	"fmt"
	"goapi/pkg/config"
	"net/smtp"
)

type Mail struct {
	user   string
	passwd string
}

// New 初始化用户名和密码
func New() Mail {
	return Mail{
		user:   config.GetString("email.user"),
		passwd: config.GetString("email.password"),
	}
}

// Send 标题 文本 目标邮箱
func (m Mail) Send(title string, text string, toId string) error {
	auth := smtp.PlainAuth("", m.user, m.passwd, "smtp.gmail.com")

	tLsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         "smtp.gmail.com",
	}

	conn, err := tls.Dial("tcp", "smtp.gmail.com:465", tLsConfig)
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, "smtp.gmail.com")
	if err != nil {
		return err
	}

	if err = client.Auth(auth); err != nil {
		return err
	}

	if err = client.Mail(m.user); err != nil {
		return err
	}

	if err = client.Rcpt(toId); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", m.user, toId, title, text)

	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	err1 := client.Quit()
	if err1 != nil {
		return err1
	}

	return nil
}
