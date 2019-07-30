package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
	"github.com/sirupsen/logrus"
	"net/http"
	"gwyApp/models"
	"gwyApp/pkg/app"
	"gwyApp/pkg/e"
	"gwyApp/service/topic_service"
)

func NextTopic(c *gin.Context) {
	appG := app.Gin{C: c}
	var decoder = schema.NewDecoder()
	req := &topic_service.TopicReq{}
	if err := decoder.Decode(req, c.Request.URL.Query()); err != nil {
		logrus.Error("decode error :", err)
		appG.Response(http.StatusOK, e.ERROR, err.Error())
		return
	}

	topic, err := topic_service.NextTopic(req)
	if err != nil {
		logrus.Error("NextTopic error :", err)
		appG.Response(http.StatusOK, e.ERROR, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, topic)
}

func NextWrongTopic(c *gin.Context) {
	appG := app.Gin{C: c}
	var decoder = schema.NewDecoder()
	req := &topic_service.TopicReq{}
	if err := decoder.Decode(req, c.Request.URL.Query()); err != nil {
		logrus.Error("decode error :", err)
		appG.Response(http.StatusOK, e.ERROR, err.Error())
		return
	}
	topic, err := topic_service.NextWrongTopic(req)
	if err != nil {
		logrus.Error("NextWrongTopic error :", err)
		appG.Response(http.StatusOK, e.ERROR_NO_WRONG_TOPIC, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, topic)
}

func NextCollect(c *gin.Context) {
	appG := app.Gin{C: c}
	var decoder = schema.NewDecoder()
	req := &topic_service.TopicReq{}
	if err := decoder.Decode(req, c.Request.URL.Query()); err != nil {
		logrus.Error("decode error :", err)
		appG.Response(http.StatusOK, e.ERROR, err.Error())
		return
	}
	topic, err := topic_service.NextCollect(req)
	if err != nil {
		logrus.Error("NextCollect error :", err)
		appG.Response(http.StatusOK, e.ERROR_NO_COLLECT, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, topic)
}

type CollectReq struct {
	AccessToken string `schema:"accessToken"`
	TopicId     int    `schema:"topic_id"`
}

func Collect(c *gin.Context) {
	appG := app.Gin{C: c}
	var decoder = schema.NewDecoder()
	req := &CollectReq{}
	if err := decoder.Decode(req, c.Request.URL.Query()); err != nil {
		logrus.Error("decode error :", err)
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	if err := topic_service.Collect(req.AccessToken, req.TopicId); err != nil {
		logrus.Error("collect error :", req.AccessToken, req.TopicId)
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func CancelCollect(c *gin.Context) {
	appG := app.Gin{C: c}
	var decoder = schema.NewDecoder()
	req := &CollectReq{}
	if err := decoder.Decode(req, c.Request.URL.Query()); err != nil {
		logrus.Error("decode error :", err)
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	if err := models.DelCollect(req.TopicId, req.AccessToken); err != nil {
		logrus.Error("del collect error :", err)
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type AnswerReq struct {
	AccessToken string `schema:"accessToken"`
	MyAnswer    string `schema:"my_answer"`
	TopicId     int    `schema:"topic_id"`
}

func Answer(c *gin.Context) {
	appG := app.Gin{C: c}
	var decoder = schema.NewDecoder()
	req := &AnswerReq{}
	if err := decoder.Decode(req, c.Request.URL.Query()); err != nil {
		logrus.Error("decode error :", err)
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	resp, err := topic_service.Answer(req.AccessToken, req.TopicId, req.MyAnswer)
	if err != nil {
		logrus.Error("answer error :", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, resp)
}
func GetAnalysis(c *gin.Context) {
	appG := app.Gin{C: c}
	var decoder = schema.NewDecoder()
	req := &AnswerReq{}
	if err := decoder.Decode(req, c.Request.URL.Query()); err != nil {
		logrus.Error("decode error :", err)
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	resp, err := topic_service.GetAnalysis(req.AccessToken, req.TopicId)
	if err != nil {
		logrus.Error("getAnalysis error :", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, resp)
}
