package ecode

const (
	OK           = 0
	RequestErr   = 400
	Unauthorized = 401
	Forbidden    = 403
	ServerErr    = 500

	RefUnknownError = 100000 // 未知错误
	SlackError      = 200001 // slack发送错误
	TelegramError   = 200002 // telegram发送错误
)
