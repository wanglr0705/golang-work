package utils

import (
	"fmt"
	"runtime"
)

func LogError(ch chan string, err error) error {
	var errStr string
	if err != nil {
		pc, file, line, ok := runtime.Caller(1)
		if ok {
			funcName := runtime.FuncForPC(pc).Name()
			errStr = fmt.Sprintf("Error: %v\nAt: %s:%d (function: %s)\n", err, file, line, funcName)
			ch <- errStr
			return err
		} else {
			errStr = fmt.Sprintf("Error: %v\n", err)
			ch <- errStr
		}
	}
	return err
}
