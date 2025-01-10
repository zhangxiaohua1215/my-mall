package errcode

import (
	"encoding/json"
	"fmt"
	"path"
	"runtime"
	"strings"
)

type AppError struct {
	code     int    
	msg      string 
	cause    error  
	occurred string // 保存由底层错误导致AppErr发生时的位置
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}

	formattedErr := struct {
		Code     int    `json:"code"`
		Msg      string `json:"msg"`
		Cause    string `json:"cause"`
		Occurred string `json:"occurred"`
	}{
		Code:     e.Code(),
		Msg:      e.Msg(),
		Occurred: e.occurred,
	}

	if e.cause != nil {
		formattedErr.Cause = e.cause.Error()
	}
	errByte, _ := json.Marshal(formattedErr)
	return string(errByte)
}

func (e *AppError) String() string {
	return e.Error()
}

func (e *AppError) Code() int {
	return e.code
}

func (e *AppError) Msg() string {
	return e.msg
}

// WithCause 在逻辑执行中出现错误, 比如dao层返回的数据库查询错误
// 可以在领域层返回预定义的错误前附加上导致错误的基础错误。
// 如果业务模块预定义的错误码比较详细, 可以使用这个方法, 反之错误码定义的比较笼统建议使用Wrap方法包装底层错误生成项目自定义Error
// 并将其记录到日志后再使用预定义错误码返回接口响应
func (e *AppError) WithCause(err error) *AppError {
	e.cause = err
	e.occurred = getAppErrOccurredInfo(2)
	return e
}

func newError(code int, msg string) *AppError {
	if _, duplicated := codes[code]; duplicated {
		panic(fmt.Sprintf("错误码 %d 不能重复, 请检查后更换", code))
	}
	return &AppError{code: code, msg: msg}
}

// Wrap 用于逻辑中包装底层函数返回的error 和 WithCause 一样都是为了记录错误链条
// 该方法生成的error 用于日志记录, 返回响应请使用预定义好的error
func Wrap(msg string, err error) *AppError {
	if err == nil {
		return nil
	}
	appErr := &AppError{code: -1, msg: msg, cause: err}
	appErr.occurred = getAppErrOccurredInfo(2)
	return appErr
}

// getAppErrOccurredInfo 获取项目中调用Wrap或者WithCause方法时的程序位置, 方便排查问题
func getAppErrOccurredInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	// file = path.Base(file)
	// idx := strings.LastIndexByte(file, '/') 
	// if idx != -1 {

	// }
	triggerInfo := fmt.Sprintf("%s:%d", lastTwoParts(file), line)
	return triggerInfo
}

func lastTwoParts(s string) string {
	parts := strings.Split(s, "/")
	if len(parts) <= 2 {
		return s
	}
	return path.Join(parts[len(parts)-2:]...)
}

func lastTwoParts2(s string) string {
	idx := strings.LastIndexByte(s, '/') 
	if idx == -1 {
		return s
	}
	idx = strings.LastIndexByte(s[:idx], '/')
	if idx == -1 {
		return s
	}
	return s[idx+1:]
}
