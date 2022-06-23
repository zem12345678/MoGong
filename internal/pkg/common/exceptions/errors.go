package exceptions

var MsgFlags = map[int]string{
	Success:        "请求成功",
	Err:            "服务器错误",
	ErrPram:        "参数错误",
	ErrTimeout:     "服务器忙碌，请稍后再试",
	InvalidToken:   "验证Token失败",
	ErrTokenFormat: "Token格式错误",
	TokenExpired:   "Token过期",
	LoginFailed:    "登陆失败",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[Err]
}

/*
通用错误error
*/

type (
	// BaseError 基本错误类型
	BaseError struct {
		message string
	}
)

// NewBaseError  初始化基本用户类型
func NewBaseError(message string) *BaseError {
	return &BaseError{message: message}
}

// Error 实现Error
func (e *BaseError) Error() string {

	return e.message
}
