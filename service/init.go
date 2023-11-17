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
	httpClient *http.Client
}

func Init() {
	MessageService = &Service{
		httpClient: http.New(time.Duration(conf.Conf.Timeout)*time.Second, 1),
	}
}
