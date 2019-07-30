package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"gwyApp/pkg/setting"
)

//接受到微信的登录请求后，拿code去访问微信后台拿openId
func GetOpenId(jsCode string) (string, error) {
	reqUrl := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		setting.WeChatSetting.AppId, setting.WeChatSetting.AppSecret, jsCode)

	resp, err := http.Get(reqUrl)
	if err != nil {
		// handle error
		log.Fatal("请求报错了", err)
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		log.Fatal("请求报错了", err)
		return "", err
	}

	fmt.Println(string(body))
	return string(body), nil
}
