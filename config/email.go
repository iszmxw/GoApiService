package config

import "goapi/pkg/config"

// email 参数
func init() {
	config.Add("email", config.StrMap{
		"user":     config.Env("EMAIL_USER", ""),
		"password": config.Env("EMAIL_PASSWORD", ""),
	})
}
