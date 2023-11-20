package http

import (
	"github.com/gin-gonic/gin"
	"ref-message-hub/common/ecode"
	gin2 "ref-message-hub/common/gin"
	"ref-message-hub/common/http"
	"ref-message-hub/common/referror"
	"ref-message-hub/model"
	"ref-message-hub/service"
)

func sendMessage(c *gin.Context) {
	var (
		param model.MessageParam
		err   error
	)
	if err = gin2.ShouldBind(c, &param); err != nil {
		http.ReturnError(c, &referror.Error{Code: ecode.RequestErr, Message: err.Error()})
		return
	}
	tg, slack, email, err := service.MessageService.SendMessage(&param)
	if err != nil {
		http.ReturnError(c, err)
		return
	}
	c.JSON(ecode.OK, http.Resp(ecode.OK, model.MessageResponse{
		Telegram: tg,
		Slack:    slack,
		Email:    email,
	}))
}
