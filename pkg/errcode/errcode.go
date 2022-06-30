package errcode

import (
	"fmt"
	"net/http"
)

type Error struct {
	code int `json:"code"`
	msg string `json:"msg"`
	details []string `json:"details"`
}

var codes = map[int]string{} // 将所有错误码对应起来，方便判断新加的code是否已存在

func NewError(code int, msg string) *Error{
	_, ok := codes[code]
	if ok {
		panic(fmt.Sprintf("错误码 %d 已存在，请更换", code))
	}

	codes[code] = msg
	return &Error{code: code, msg: msg}
}


func(e *Error) Error() string{
	return fmt.Sprintf("错误码：%d, 错误信息：%s", e.Code(), e.Msg())
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Msg() string {
	return e.msg
}

func (e *Error) Msgf(args []interface{} )  string  {
	return fmt.Sprintf(e.msg, args...)
}


func (e *Error) Details() []string {
	return e.details
}

// 将详情添加到原Error，产生新的Error
func (e *Error) WithDetails(details ...string) *Error {
	newError := *e
	newError.details = []string{}
	for _, d := range details{
		newError.details = append(newError.details, d)
	}
	return &newError
}


// 转化为http的响应代码
func (e *Error) StatusCode() int {
	switch e.Code() {
	case Success.Code():
		return http.StatusOK
	case ServerError.Code():
		return http.StatusInternalServerError
	case InvalidParams.Code():
		return http.StatusBadRequest
	case UnauthorizedAuthNotExist.Code():
		fallthrough
	case UnauthorizedTokenError.Code():
		fallthrough
	case UnauthorizedTokenGenerate.Code():
		fallthrough
	case UnauthorizedTokenTimeout.Code():
		return http.StatusUnauthorized
	case TooManyRequests.Code():
		return http.StatusTooManyRequests
	}

	return http.StatusInternalServerError
}









