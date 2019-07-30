package setting

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"time"
)

var cfg *ini.File

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		logrus.Fatalf("Cfg.MapTo RedisSetting err: %v", err)
	}
}

type MyTextFormatter struct {
}

func (f MyTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timeStr := entry.Time.Format(time.RFC3339)
	str := fmt.Sprintf("%s[%s] %s\n", timeStr, entry.Level.String(), entry.Message)
	return []byte(str), nil
}
func setUpLog() {
	logrus.SetFormatter(&MyTextFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}
func SetUp() {
	var err error
	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}
	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)
	mapTo("redis", RedisSetting)
	mapTo("wechat", WeChatSetting)
	setUpLog()
}

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

var DatabaseSetting = &Database{}

type App struct {
	FileSavePath string
}

var AppSetting = &App{}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

var RedisSetting = &Redis{}

type WeChat struct {
	AppId     string
	AppSecret string
}

var WeChatSetting = &WeChat{}
