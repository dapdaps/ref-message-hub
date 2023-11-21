package service

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/smtp"
	"ref-message-hub/common/ecode"
	"ref-message-hub/common/log"
	model2 "ref-message-hub/common/model"
	"ref-message-hub/common/referror"
	"ref-message-hub/conf"
	"ref-message-hub/model"
	"strings"
	"sync"
)

func (s *Service) SendMessage(param *model.MessageParam) (tg bool, slack bool, email bool, err error) {
	var (
		title      = param.Title
		product    map[string][]string
		users      []string
		tgError    error
		slackError error
		emailError error
		wg         sync.WaitGroup
	)
	if !param.Slack && !param.Telegram && !param.Email {
		err = &referror.Error{Code: ecode.ParamError, Message: "must choose one among Slack, Telegram, and email"}
		return
	}
	if _, ok := conf.Conf.Levels[param.Level]; !ok {
		err = &referror.Error{Code: ecode.ParamError, Message: "illegal level"}
		return
	}
	product, ok := conf.Conf.Product[param.Product]
	if !ok {
		err = &referror.Error{Code: ecode.ParamError, Message: "not find product"}
		return
	}
	users, ok = product[param.Level]
	if !ok {
		err = &referror.Error{Code: ecode.ParamError, Message: "not find users"}
		return
	}
	if len(title) == 0 {
		title = param.Product + " " + param.Type
	}

	if param.Telegram {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if tgError = s.sendTelegram(param, users, title); tgError != nil {
				tg = false
			} else {
				tg = true
			}
		}()
	}

	if param.Slack {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if slackError = s.sendSlackMessage(param, users, title); slackError != nil {
				slack = false
			} else {
				slack = true
			}
		}()
	}

	if param.Email {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if emailError = s.sendEmail(param, users, title); emailError != nil {
				email = false
			} else {
				email = true
			}
		}()
	}

	wg.Wait()
	return
}

func (s *Service) sendTelegram(param *model.MessageParam, users []string, title string) (err error) {
	var (
		text = ""
	)
	tgChannel, ok := conf.Conf.Telegram.Channel["monitor"]
	if !ok {
		return
	}
	text += title + "/n"
	for _, user := range users {
		text += conf.Conf.Telegram.Users[user]
	}
	text += param.Content
	response := &model.TelegramResponse{}
	url := "https://api.telegram.org/bot%s/sendMessage"
	request := &model.TelegramMessage{
		ChatId: tgChannel,
		Text:   text,
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

func (s *Service) sendSlackMessage(param *model.MessageParam, users []string, title string) (err error) {
	var (
		text = ""
	)
	slackWebHook, ok := conf.Conf.Slack.Channel["monitor"]
	if !ok {
		return
	}
	for _, user := range users {
		text += "<@" + conf.Conf.Slack.Users[user] + "> "
	}
	text += "\n" + param.Content
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
			Text: text,
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

// smtp send email
func (s *Service) sendEmail(param *model.MessageParam, users []string, title string) (err error) {
	var (
		to     []string
		client *smtp.Client
	)
	for _, user := range users {
		to = append(to, conf.Conf.Email.Users[user])
	}
	headers := make(map[string]string)
	headers["From"] = conf.Conf.Email.Sender
	headers["To"] = strings.Join(to, ",")
	headers["Subject"] = title

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + param.Content

	auth := smtp.PlainAuth("", conf.Conf.Email.Sender, conf.Conf.Email.Password, conf.Conf.Email.Host)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         conf.Conf.Email.Host,
	}
	conn, err := tls.Dial("tcp", conf.Conf.Email.Host+":465", tlsConfig)
	if err != nil {
		log.Error("sendEmail tls.Dial error: %v", err)
		return
	}
	client, err = smtp.NewClient(conn, conf.Conf.Email.Host)
	if err = client.Auth(auth); err != nil {
		log.Error("sendEmail client.Auth error: %v", err)
		return
	}
	if err = client.Mail(conf.Conf.Email.Sender); err != nil {
		log.Error("sendEmail client.Mail error: %v", err)
		return
	}
	for _, address := range to {
		if err = client.Rcpt(address); err != nil {
			log.Error("sendEmail client.Rcpt error: %v", err)
			return
		}
	}
	w, err := client.Data()
	if err != nil {
		log.Error("sendEmail client.Data error: %v", err)
		return
	}
	if _, err = w.Write([]byte(message)); err != nil {
		log.Error("sendEmail Write error: %v", err)
		return
	}
	if err = w.Close(); err != nil {
		log.Error("sendEmail Close error: %v", err)
		return
	}
	_ = client.Quit()
	return
}

//ses send emial
//func (s *Service) sendEmail(param *model.MessageParam, users []string) (err error) {
//	//toAddress, ok := conf.Conf.Email.Users[param.Level]
//	//if !ok {
//	//	return
//	//}
//	//input := &ses.SendEmailInput{
//	//	Destination: &types.Destination{
//	//		ToAddresses: toAddress,
//	//	},
//	//	Message: &types.Message{
//	//		Body: &types.Body{
//	//			Text: &types.Content{
//	//				Data: aws.String(param.Content),
//	//			},
//	//		},
//	//		Subject: &types.Content{
//	//			Data: aws.String("Alert"),
//	//		},
//	//	},
//	//	Source: aws.String(conf.Conf.Email.Sender),
//	//}
//	//ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
//	//defer cancel()
//	//_, err = s.email.Client.SendEmail(ctx, input)
//	//if err != nil {
//	//	log.Error("sendEmail client.SendEmail error: %v", err)
//	//	return
//	//}
//	//log.Info("sendEmail success")
//	return
//}
