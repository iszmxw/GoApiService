package main

import (
	"fmt"
	"goapi/config"
	"goapi/pkg/logger"
	"math/big"
)

func init() {
	// 初始化配置信息
	config.Initialize()
	// 定义日志目录
	logger.Init("test")
}

func hexToBigInt(hex string) *big.Int {
	n := new(big.Int)
	n, _ = n.SetString(hex, 16)
	return n
}

func main() {
	fmt.Println(hexToBigInt("12a05f200"))
}
