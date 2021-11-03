package abnormal

import (
	"errors"
	"fmt"
	"goapi/pkg/gmail"
	"goapi/pkg/logger"
	"runtime/debug"
)

func Stack(description string) {
	if r := recover(); r != nil {
		// 收集错误堆栈信息 异常
		errInfo := string(debug.Stack())
		err1 := errors.New(fmt.Sprintf(": %v\n", description))
		err2 := errors.New(fmt.Sprintf("Recover: %v\n", r))
		err3 := errors.New(fmt.Sprintf("详细堆栈错误信息: %v\n", errInfo))
		logger.Error(err1)
		logger.Error(err2)
		logger.Error(err3)
		err := gmail.New().Send(description, errInfo, "543619552@qq.com")
		if err != nil {
			logger.Error(err)
		}
	}
}
