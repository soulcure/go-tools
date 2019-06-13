package models

import (
	"fmt"
	"github.com/kataras/iris"
)

const (
	OK          = 0
	LoginErr    = -4001 //登录失败
	RegisterErr = -4002 //注册失败
	NotLogin    = -4003 //未登录
	TokenExp    = -4004 //token失效
	ParamErr    = -4005 //请求参数异常
	UnknownErr  = -4201 //未知异常

)

var msg = map[int]string{
	OK:          "请求成功",
	LoginErr:    "登录失败",
	RegisterErr: "注册失败",
	NotLogin:    "未登录",
	TokenExp:    "token失效",
	ParamErr:    "请求参数异常",
	UnknownErr:  "未知异常",
}

type ProtocolRsp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func Msg(code int) string {
	str, ok := msg[code]
	if ok {
		return str
	}
	return msg[UnknownErr]
}

func (json *ProtocolRsp) SetCode(code int) {
	json.Code = code
	json.Msg = Msg(code)
}

func (json *ProtocolRsp) ResponseWriter(ctx iris.Context) {
	if _, err := ctx.JSON(json); err != nil {
		fmt.Println(err)
	}
}
