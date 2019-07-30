package topic_service

import (
	"errors"
	"github.com/sirupsen/logrus"
	"gwyApp/common"
	"gwyApp/models"
	"gwyApp/pkg/gredis"
	"strconv"
)

type TopicReq struct {
	AccessToken  string `schema:"accessToken"`
	IsBegin      bool   `schema:"is_begin"`
	CurrentIndex int    `schema:"current_index"`
	Operate      string `schema:"operate"`
}

func getBeginCollect(req *TopicReq) (*Topic, error) {
	collects, err := models.GetCollects(req.AccessToken)
	if err != nil {
		logrus.Error("GetCollects error :", err)
		return nil, err
	}
	_, err = gredis.Delete(common.COLLECT_LIST + req.AccessToken)
	if err != nil {
		logrus.Error("delete error :", err)
		return nil, err
	}
	for _, collect := range collects {
		_, err = gredis.RPush(common.COLLECT_LIST+req.AccessToken, strconv.Itoa(collect.TopicId))
		if err != nil {
			logrus.Error("lpush redis error :", err)
			return nil, err
		}
	}
	return getTopicByIndex(common.COLLECT_LIST, req.AccessToken, 0)
}
func NextCollect(req *TopicReq) (*Topic, error) {
	if req.IsBegin {
		return getBeginCollect(req)
	}
	if req.Operate == common.OPERATE_LAST {
		return getTopicByIndex(common.COLLECT_LIST, req.AccessToken, req.CurrentIndex-1)
	}
	if req.Operate == common.OPERATE_NEXT {
		return getTopicByIndex(common.COLLECT_LIST, req.AccessToken, req.CurrentIndex+1)
	}
	return nil, errors.New("no topic")

}
