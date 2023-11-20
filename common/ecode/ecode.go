package ecode

const (
	OK           = 0
	RequestErr   = 400
	Unauthorized = 401
	Forbidden    = 403
	ServerErr    = 500

	RefUnknownError = 100000
	ParamError      = 110000 // slack error
	SlackError      = 200001 // slack error
	TelegramError   = 200002 // telegram error
)
