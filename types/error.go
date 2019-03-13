package types

import "fmt"

var (
	ErrSuccess    = NewErrorCode(0, "")
	ErrUnknown    = NewErrorCode(1, "未知错误")
	ErrInternal   = NewErrorCode(2, "内部错误")
	ErrException  = NewErrorCode(3, "内部异常")
	ErrNotSupport = NewErrorCode(4, "不支持的操作")
	ErrExist      = NewErrorCode(5, "已存在")
	ErrNotExist   = NewErrorCode(6, "不存在")

	ErrInput        = NewErrorCode(11, "参数错误")
	ErrInputInvalid = NewErrorCode(12, "参数无效")

	ErrTokenEmpty   = NewErrorCode(101, "缺少凭证")
	ErrTokenInvalid = NewErrorCode(101, "凭证无效")
	ErrTokenIllegal = NewErrorCode(101, "凭证非法")

	ErrLoginCaptchaInvalid           = NewErrorCode(201, "验证码无效")
	ErrLoginAccountNotExit           = NewErrorCode(202, "账号不存在")
	ErrLoginPasswordInvalid          = NewErrorCode(203, "密码不正确")
	ErrLoginAccountOrPasswordInvalid = NewErrorCode(204, "账号或密码不正确")
)

type Error interface {
	Code() int
	Summary() string
	Detail() string
}

func NewError(errCode ErrorCode, detail ...interface{}) Error {
	return &innerError{
		innerErrorCode: innerErrorCode{
			code:    errCode.Code(),
			summary: errCode.Summary(),
		},
		detail: fmt.Sprint(detail...),
	}
}

type innerError struct {
	innerErrorCode
	detail string
}

func (s *innerError) Code() int {
	return s.code
}

func (s *innerError) Summary() string {
	return s.summary
}

func (s *innerError) Detail() string {
	return s.detail
}

func (s *innerError) ErrCode() ErrorCode {
	return &s.innerErrorCode
}

type ErrorCode interface {
	Code() int
	Summary() string
}

func NewErrorCode(code int, summary interface{}) ErrorCode {
	return &innerErrorCode{
		code:    code,
		summary: fmt.Sprint(summary),
	}
}

type innerErrorCode struct {
	code    int
	summary string
}

func (s *innerErrorCode) Code() int {
	return s.code
}

func (s *innerErrorCode) Summary() string {
	return s.summary
}
