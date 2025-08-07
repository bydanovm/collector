package redis

import (
	"runtime"
	"strings"
)

func GetFunctionName() []byte {
	pc, _, _, _ := runtime.Caller(1)
	fullFuncName := runtime.FuncForPC(pc).Name()
	funcName := strings.Split(fullFuncName, "/")
	return []byte(funcName[len(funcName)-1])
}
