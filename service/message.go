package service

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"net/http"
	"ref-message-hub/common/ecode"
	"ref-message-hub/common/log"
	model2 "ref-message-hub/common/model"
	"ref-message-hub/common/referror"
	"ref-message-hub/conf"
	"ref-message-hub/model"
	"sync"
)

func (s *Service) SendMessage(param *model.MessageParam) (tg bool, slack bool, email bool) {
	var (
		tgError    error
		slackError error
		emailError error
		wg         sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if tgError = s.sendTelegram(param); tgError != nil {
			tg = false
		} else {
			tg = true
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if slackError = s.sendSlackMessage(param); slackError != nil {
			slack = false
		} else {
			slack = true
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if emailError = s.sendEmail(param); emailError != nil {
			email = false
		} else {
			email = true
		}
	}()

	wg.Wait()
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
	if !response.Ok {
		err = &referror.Error{Code: ecode.TelegramError, Message: "failed send to telegram"}
		return
	}
	log.Info("sendTelegram success")
	return
}

func (s *Service) sendSlackMessage(param *model.MessageParam) (err error) {
	var (
		title = "Alert"
	)
	slackWebHook, ok := conf.Conf.SlackWebHooks[param.Slack]
	if !ok {
		return
	}
	if len(param.Title) > 0 {
		title = param.Title
	}
	var slackMessages []model.SlackMessageBlock
	slackMessages = append(slackMessages, model.SlackMessageBlock{
		Type: "header",
		Text: model.SlackMessage{
			Type: "plain_text",
			Text: title,
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
	log.Info("sendSlackMessage success")
	return
}

func (s *Service) sendEmail(param *model.MessageParam) (err error) {
	var (
		title = "Alert"
	)
	if !param.Email {
		return
	}
	if len(param.Title) > 0 {
		title = param.Title
	}
	toAddress, ok := conf.Conf.Email.Receiver[param.Level]
	if !ok {
		return
	}
	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: toAddress,
		},
		Message: &types.Message{
			Body: &types.Body{
				Text: &types.Content{
					Data: aws.String(param.Content),
				},
			},
			Subject: &types.Content{
				Data: aws.String(title),
			},
		},
		Source: aws.String(conf.Conf.Email.Sender),
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	_, err = s.email.Client.SendEmail(ctx, input)
	if err != nil {
		log.Error("sendEmail client.SendEmail error: %v", err)
		return
	}
	log.Info("sendEmail success")
	return
}
