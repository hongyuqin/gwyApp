package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
	"github.com/sirupsen/logrus"
	"gwyApp/common"
	"gwyApp/models"
	"gwyApp/pkg/app"
	"gwyApp/pkg/e"
	"gwyApp/service/user_service"
	"net/http"
)

func GetUser(c *gin.Context) {
	appG := app.Gin{C: c}
	openId := c.Query("accessToken")
	user, err := models.SelectUserByOpenId(openId)
	if err != nil {
		logrus.Info("获取用户错误 :" + err.Error())
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, user)
}

type Resp struct {
	SessionKey string `json:"session_key"`
	OpenId     string `json:"openid"`
}

func GetAccessToken(c *gin.Context) {
	appG := app.Gin{C: c}
	jsCode := c.Query("code")
	logrus.Info("jsCode is :", jsCode)
	//body, _ := util.GetOpenId(jsCode)
	body := `{"openid":"abc","session_key":"xx"}`
	logrus.Info("body is :", body)

	resp := &Resp{}
	if err := json.Unmarshal([]byte(body), resp); err != nil {
		logrus.Error("marshal error :", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
	}
	//1.保存用户
	/*exist, err := gredis.SIsmember(common.OPENID_SET, resp.OpenId)
	if err != nil {
		logrus.Error("redis error :", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}*/
	exist := false
	if !exist {
		//不存在该用户 插入redis,更新数据库
		/*if err := gredis.SAdd(common.OPENID_SET, resp.OpenId); err != nil {
			logrus.Error("redis error :", err)
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
			return
		}*/
		if err := models.AddUser(models.User{
			OpenId: resp.OpenId,
		}); err != nil {
			logrus.Error("addUser error :", err)
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
			return
		}
		if err := models.AddSetting(models.Setting{
			OpenId:       resp.OpenId,
			Region:       common.REGION_SZ,
			ExamType:     common.EXAM_TYPE_CIVIL,
			DailyNeedNum: 100,
		}); err != nil {
			logrus.Error("addSetting error :", err)
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
			return
		}

	}
	//2.返回openId
	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"accessToken": resp.OpenId,
	})
	return
}

func HomePage(c *gin.Context) {
	appG := app.Gin{C: c}
	openId := c.Query("accessToken")
	homeDetail, err := user_service.Home(openId)
	if err != nil {
		logrus.Error("Home error :", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, homeDetail)
}

func Plan(c *gin.Context) {
	appG := app.Gin{C: c}
	var decoder = schema.NewDecoder()
	req := &user_service.PlanReq{}
	if err := decoder.Decode(req, c.Request.URL.Query()); err != nil {
		logrus.Error("decode error :", err)
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	err := user_service.Plan(req)
	if err != nil {
		logrus.Error("Plan error :", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func GetPlan(c *gin.Context) {
	appG := app.Gin{C: c}
	openId := c.Query("accessToken")
	setting, err := models.SelectSettingByOpenId(openId)
	if err != nil || setting == nil {
		logrus.Error("SelectSettingByOpenId openId error :", openId, err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	data := make(map[string]interface{})
	data["region"] = setting.Region
	data["exam_type"] = setting.ExamType
	data["daily_need_num"] = setting.DailyNeedNum
	appG.Response(http.StatusOK, e.SUCCESS, data)
}
