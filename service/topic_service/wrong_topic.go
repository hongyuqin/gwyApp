package topic_service

import (
	"errors"
	"github.com/sirupsen/logrus"
	"gwyApp/common"
	"gwyApp/models"
	"gwyApp/pkg/gredis"
	"strconv"
)

func getBeginWrongTopic(req *TopicReq) (*Topic, error) {
	wrongTopics, err := models.GetWrongTopics(req.AccessToken)
	if err != nil {
		logrus.Error("GetWrongTopics error :", err)
		return nil, err
	}
	_, err = gredis.Delete(common.WRONG_TOPIC_LIST + req.AccessToken)
	if err != nil {
		logrus.Error("delete error :", err)
		return nil, err
	}
	if len(wrongTopics) == 0 {
		logrus.Error("no wrong topics")
		return nil, errors.New("当前无错题")
	}
	for _, wrongTopic := range wrongTopics {
		_, err = gredis.RPush(common.WRONG_TOPIC_LIST+req.AccessToken, strconv.Itoa(wrongTopic.TopicId))
		if err != nil {
			logrus.Error("rpush redis error :", err)
			return nil, err
		}
	}
	return getTopicByIndex(common.WRONG_TOPIC_LIST, req.AccessToken, 0)
}
func NextWrongTopic(req *TopicReq) (*Topic, error) {
	if req.IsBegin {
		return getBeginWrongTopic(req)
	}
	if req.Operate == common.OPERATE_LAST {
		return getTopicByIndex(common.COLLECT_LIST, req.AccessToken, req.CurrentIndex-1)
	}
	if req.Operate == common.OPERATE_NEXT {
		return getTopicByIndex(common.COLLECT_LIST, req.AccessToken, req.CurrentIndex+1)
	}
	return nil, errors.New("no topic")

}
