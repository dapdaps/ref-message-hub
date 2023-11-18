package conf

import (
	"flag"
	"github.com/BurntSushi/toml"
	"ref-message-hub/common/log"
)

var (
	confPath string
	Conf     = &Config{}
)

type Config struct {
	Debug         bool
	Timeout       int64
	Log           *log.Config
	Port          int
	AllowOrigins  []string
	SlackWebHooks map[string]string
	Telegram      *TelegramConfig
	Email         *EmailConfig
}

type TelegramConfig struct {
	BotToken  string
	ChatGroup map[string]int64
}

type EmailConfig struct {
	Region       string
	Sender       string
	AccessID     string
	AccessSecret string
	Receiver     map[string][]string
}

func init() {
	flag.StringVar(&confPath, "conf", "", "config path")
}

func Init() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	if err != nil {
		log.Error("error decoding [%v]:%v", confPath, err)
		return
	}
	return
}
