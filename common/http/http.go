package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ref-message-hub/common/ecode"
	"ref-message-hub/common/referror"
	"time"
)

func Resp(code int, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"code": code,
		"data": data,
	}
}

func RespWithMsg(code int, message string) map[string]interface{} {
	return map[string]interface{}{
		"code":       code,
		"message":    message,
		"time_stamp": time.Now().UnixNano() / 1e6,
	}
}

func ReturnError(c *gin.Context, err error) {
	var statusCode int
	switch e := err.(type) {
	case *referror.Error:
		switch e.Code {
		case ecode.RequestErr:
			statusCode = http.StatusBadRequest
		case ecode.Unauthorized:
			statusCode = http.StatusUnauthorized
		case ecode.Forbidden:
			statusCode = http.StatusForbidden
		default:
			statusCode = http.StatusOK
		}
		c.JSON(statusCode, err)
	//case *referror.ThirdPartyError:
	//	statusCode = http.StatusInternalServerError
	//	c.JSON(statusCode, Resp(ecode.ExternalError, err))
	default:
		statusCode = http.StatusInternalServerError
		c.JSON(statusCode, RespWithMsg(ecode.ServerErr, err.Error()))
	}
}
