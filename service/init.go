package service

import (
	"ref-message-hub/common/http"
	"ref-message-hub/conf"
	"time"
)

var (
	MessageService *Service
)

type Service struct {
	timeout    time.Duration
	httpClient *http.Client
}

func Init() (err error) {
	MessageService = &Service{
		timeout:    time.Duration(conf.Conf.Timeout) * time.Second,
		httpClient: http.New(time.Duration(conf.Conf.Timeout)*time.Second, 1),
	}
	return
}
