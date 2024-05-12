package code

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy
	CodeUnAuth
	CodeInvalidToken
	CodeAlreadyJoined
	CodeGroupNotExist
	CodeNotJoinGroup
)

var msgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "请求参数错误",
	CodeUserExist:       "用户名已存在",
	CodeUserNotExist:    "用户名不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务繁忙",
	CodeUnAuth:          "需要登录",
	CodeInvalidToken:    "无效的token",
	CodeAlreadyJoined:   "已经加入该群",
	CodeGroupNotExist:   "群不存在",
	CodeNotJoinGroup:    "未加入该群",
}

func (c ResCode) Msg() string {
	msg, ok := msgMap[c]
	if !ok {
		return msgMap[CodeServerBusy]
	}
	return msg
}
