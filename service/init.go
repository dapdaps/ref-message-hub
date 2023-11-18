package service

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"ref-message-hub/common/http"
	"ref-message-hub/common/log"
	"ref-message-hub/conf"
	"time"
)

var (
	MessageService *Service
)

type Service struct {
	timeout    time.Duration
	httpClient *http.Client
	email      *Email
}

type Email struct {
	Sender string
	Client *ses.Client
}

func Init() (err error) {
	MessageService = &Service{
		timeout:    time.Duration(conf.Conf.Timeout) * time.Second,
		httpClient: http.New(time.Duration(conf.Conf.Timeout)*time.Second, 1),
		email: &Email{
			Sender: conf.Conf.Email.Sender,
		},
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     conf.Conf.Email.AccessID,
				SecretAccessKey: conf.Conf.Email.AccessSecret,
				SessionToken:    "",
			},
		}),
	)
	if err != nil {
		log.Error("Init config.LoadDefaultConfig error: %v", err)
		return
	}
	MessageService.email.Client = ses.NewFromConfig(cfg)
	return
}
