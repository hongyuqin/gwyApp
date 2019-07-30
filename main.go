package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gwyApp/models"
	"gwyApp/pkg/gredis"
	"gwyApp/pkg/setting"
	"gwyApp/routers"
)

func init() {
	setting.SetUp()
	models.Setup()
	if err := gredis.Setup(); err != nil {
		logrus.Panic("redis init error")
	}
}
func main() {
	r := routers.InitRouter()

	r.Run(fmt.Sprintf(":%d", setting.ServerSetting.HttpPort))
}
