package helpers

// 转换器
import (
	"bytes"
	"encoding/json"
	"goapi/pkg/logger"
	"strconv"
)

// Int64ToString 将 int64 转换为 string
func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}

// IntToString 整型转字符串
func IntToString(num int) string {
	return strconv.Itoa(num)
}

// Uint64ToString 将 uint64 转换为 string
func Uint64ToString(num uint64) string {
	return strconv.FormatUint(num, 10)
}

// StringToInt 将字符串转换为 int
func StringToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		logger.Error(err)
	}
	return i
}

// Uint2String 将字符串转换为 int
func Uint2String(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		logger.Error(err)
	}
	return i
}

// Struct2json Struct转换json
func Struct2json(value interface{}) string {
	bs, _ := json.Marshal(value)
	var out bytes.Buffer
	err := json.Indent(&out, bs, "", "\t")
	if err != nil {
		logger.Error(err)
	}
	return out.String()
}
