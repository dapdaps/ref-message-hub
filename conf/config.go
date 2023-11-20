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
	Debug        bool
	Timeout      int64
	Log          *log.Config
	Port         int
	AllowOrigins []string
	Slack        *SlackConfig
	Telegram     *TelegramConfig
	Email        *EmailConfig
	Product      map[string]map[string][]string
	Levels       map[string]string
}

type SlackConfig struct {
	Channel map[string]string
	Users   map[string]string
}

type TelegramConfig struct {
	BotToken string
	Channel  map[string]int64
	Users    map[string]string
}

type EmailConfig struct {
	Sender   string
	Password string
	Host     string
	Users    map[string]string
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
	Conf.Levels = map[string]string{}
	Conf.Levels["critial"] = ""
	Conf.Levels["high"] = ""
	Conf.Levels["medium"] = ""
	return
}
