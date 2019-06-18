package models

import (
	"fmt"
	"github.com/kataras/iris"
)

const (
	OK              = 0
	RegisterErrCode = -4001 //注册失败
	LoginErrCode    = -4002 //登录失败
	NotLoginCode    = -4003 //未登录
	TokenExpCode    = -4004 //token失效
	ParamErrCode    = -4005 //请求参数异常
	NoFoundErrCode  = -4201 //未知异常
	UnknownErrCode  = -4202 //未知异常

)

const (
	SUCCESS = "success"

	RegisterUserNameEmptyErr  = "Registered username is empty"
	RegisterUserNameFormatErr = "Registered username format error:Number or letter does not limit capitalization 6-30 digits"

	RegisterEmailEmptyErr  = "Registered email is empty"
	RegisterEmailFormatErr = "Registered email format error:Number or letter or symbol 6-30 digits"

	RegisterMobileEmptyErr  = "Registered mobile is empty"
	RegisterMobileFormatErr = "Registered mobile format error"

	RegisterPassWordEmptyErr  = "Registered password is empty"
	RegisterPassWordFormatErr = "Registered password format error:Number or letter or symbol 6-30 digits"

	LoginErrUserNameOrEmailEmptyErr = "Login username or email is empty"
	LoginErrPassWordEmptyErr        = "Login password is empty"
	LoginErrPassWordFormatErr       = "Login password format error:Number + letter + symbol 6-30 digits"

	LoginUserNameFormatErr = "login username format error:Number or letter does not limit capitalization 6-30 digits"
	LoginEmailFormatErr    = "login email format error:Number or letter or symbol 6-30 digits"

	TokenErr        = "Token Error"
	TokenExpiredErr = "Token Expired"
	NoFoundErr      = "NoFound Error"
	UnknownErr      = "Unknown Error"
)

type ProtocolRsp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func (json *ProtocolRsp) ResponseWriter(ctx iris.Context) {
	if _, err := ctx.JSON(json); err != nil {
		fmt.Println(err)
	}
}
