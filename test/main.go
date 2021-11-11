package main

import (
	"github.com/skip2/go-qrcode"
	"goapi/pkg/logger"
)

func init() {
}

func main() {
	url := "https://baidu.com"
	png, _ := qrcode.Encode(url, qrcode.Medium, 256)
	logger.Info(png)
}
