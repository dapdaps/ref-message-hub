package service

import (
	"fmt"
	"net/http"
	"ref-message-hub/common/ecode"
	"ref-message-hub/common/log"
	model2 "ref-message-hub/common/model"
	"ref-message-hub/common/referror"
	"ref-message-hub/conf"
	"ref-message-hub/model"
)

func (s *Service) SendMessage(param *model.MessageParam) (err error) {
	if err = s.sendTelegram(param); err != nil {
		return
	}
	if err = s.sendSlackMessage(param); err != nil {
		return
	}
	return
}

func (s *Service) sendTelegram(param *model.MessageParam) (err error) {
	tgGroup, ok := conf.Conf.Telegram.ChatGroup[param.TgGroup]
	if !ok {
		return
	}
	response := &model.TelegramResponse{}
	url := "https://api.telegram.org/bot%s/sendMessage"
	request := &model.TelegramMessage{
		ChatId: tgGroup,
		Text:   param.Content,
	}
	err = s.httpClient.PostJSON(fmt.Sprintf(url, conf.Conf.Telegram.BotToken), nil, request, response)
	if err != nil {
		log.Error("sendTelegram error: %v", err)
		return
	}
	if response.Ok {
		return
	}
	err = &referror.Error{Code: ecode.TelegramError, Message: "failed send to telegram"}
	return
}

func (s *Service) sendSlackMessage(param *model.MessageParam) (err error) {
	slackWebHook, ok := conf.Conf.SlackWebHooks[param.Slack]
	if !ok {
		return
	}
	var slackMessages []model.SlackMessageBlock
	slackMessages = append(slackMessages, model.SlackMessageBlock{
		Type: "header",
		Text: model.SlackMessage{
			Type: "plain_text",
			Text: "反馈",
		},
	})
	slackMessages = append(slackMessages, model.SlackMessageBlock{
		Type: "section",
		Text: model.SlackMessage{
			Type: "mrkdwn",
			Text: param.Content,
		},
	})
	request := map[string][]model.SlackMessageBlock{}
	request["blocks"] = slackMessages
	response := &model2.BaseResponse{}
	err = s.httpClient.PostJSON(slackWebHook, nil, request, response)
	if err != nil {
		log.Error("sendSlackMessage error: %v", err)
		return
	}
	if response.HttpStatusCode != http.StatusOK {
		err = &referror.Error{Code: ecode.SlackError, Message: response.HttpBodyText}
		return
	}
	return
}
