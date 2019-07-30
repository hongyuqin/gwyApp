package topic_service

import (
	"errors"
	"github.com/sirupsen/logrus"
	"gwyApp/common"
	"gwyApp/models"
	"gwyApp/pkg/gredis"
	"gwyApp/pkg/util"
	"strconv"
)

type Topic struct {
	Index     int    `json:"index"`
	TopicId   int    `json:"topic_id"`
	TopicName string `json:"topic_name"`
	Photo     string `json:"photo"`
	OptionA   string `json:"option_a"`
	OptionB   string `json:"option_b"`
	OptionC   string `json:"option_c"`
	OptionD   string `json:"option_d"`
	Collect   bool   `json:"collect"`
}

//提交答案的返回
type AnswerResp struct {
	Analysis
	MyAnswer string `json:"my_answer"`
	Right    bool   `json:"right"`
}

type Analysis struct {
	Answer        string `json:"answer"`
	TopicAnalysis string `json:"topic_analysis"`
	WrongNum      int    `json:"wrong_num"`
	RightNum      int    `json:"right_num"`
}

func getBeginTopic(req *TopicReq) (*Topic, error) {
	allTopicsMap := make(map[int]struct{})
	setting, err := models.SelectSettingByOpenId(req.AccessToken)
	if err != nil {
		logrus.Error("get setting error")
		return nil, err
	}
	region := setting.Region
	elementTypeOne := setting.ElementTypeOne
	elementTypeTwo := setting.ElementTypeTwo
	topics, err := models.GetTopics(&models.Topic{
		Region:         region,
		ElementTypeOne: elementTypeOne,
		ElementTypeTwo: elementTypeTwo,
	})
	if err != nil {
		logrus.Error("GetTopics error :", err)
		return nil, err
	}
	topicIdCol := make([]int, 0)
	for _, topic := range topics {
		topicIdCol = append(topicIdCol, topic.ID)
		allTopicsMap[topic.ID] = struct{}{}
	}
	//2.拿到已做题目的全部id
	doneTopics, err := models.GetDoneTopics(req.AccessToken)
	if err != nil {
		logrus.Error("GetDoneTopics error :", err)
		return nil, err
	}
	for _, topic := range doneTopics {
		delete(allTopicsMap, topic.TopicId)
	}
	_, err = gredis.Delete(common.TOPIC_LIST + req.AccessToken)
	if err != nil {
		logrus.Error("delete error :", err)
		return nil, err
	}
	for topicId := range allTopicsMap {
		_, err = gredis.RPush(common.TOPIC_LIST+req.AccessToken, strconv.Itoa(topicId))
		if err != nil {
			logrus.Error("lpush redis error :", err)
			return nil, err
		}
	}
	return getTopicByIndex(common.TOPIC_LIST, req.AccessToken, 0)
}
func getTopicByIndex(prefix, openId string, index int) (*Topic, error) {
	if index < 0 {
		return nil, errors.New("已经是第一题")
	}
	//查看长度
	len, err := gredis.LLen(prefix + openId)
	if err != nil {
		logrus.Error("LLen error :", err)
		return nil, err
	}
	if index >= len {
		return nil, errors.New("已经是最后一题")
	}
	topicId, err := gredis.LIndex(prefix+openId, index)
	if err != nil {
		logrus.Error("LIndex error :", err)
		return nil, err
	}
	topic, err := models.GetTopic(topicId)
	if err != nil {
		logrus.Error("GetTopic error :", err)
		return nil, err
	}
	collect, err := models.ExistCollect(openId, topicId)
	if err != nil {
		logrus.Error("ExistCollect error :", err)
		return nil, err
	}
	return &Topic{
		Index:     index,
		TopicId:   topic.ID,
		Photo:     topic.Photo,
		TopicName: topic.TopicName,
		OptionA:   topic.OptionA,
		OptionB:   topic.OptionB,
		OptionC:   topic.OptionC,
		OptionD:   topic.OptionD,
		Collect:   collect,
	}, nil
}
func NextTopic(req *TopicReq) (*Topic, error) {
	if req.IsBegin {
		return getBeginTopic(req)
	}
	if req.Operate == common.OPERATE_LAST {
		return getTopicByIndex(common.TOPIC_LIST, req.AccessToken, req.CurrentIndex-1)
	}
	if req.Operate == common.OPERATE_NEXT {
		return getTopicByIndex(common.TOPIC_LIST, req.AccessToken, req.CurrentIndex+1)
	}
	return nil, errors.New("no topic")

}
func Answer(openId string, topicId int, myAnswer string) (*AnswerResp, error) {
	//1.看下答案对不对
	topic, err := models.GetTopic(topicId)
	if err != nil {
		logrus.Error("GetTopic error :", err)
		return nil, err
	}
	if topic == nil {
		logrus.Error("GetTopic empty :", topicId)
		return nil, err
	}
	//2.对的话更新题目的正确数，错的话更新题目错误数
	right := true
	if myAnswer == topic.Answer {
		topic.RightNum = topic.RightNum + 1
		//往用户错题表删除数据
		if err = models.DelWrongTopic(models.WrongTopic{
			OpenId:  openId,
			TopicId: topicId,
		}); err != nil {
			logrus.Error("DelWrongTopic error :", err)
			return nil, err
		}
		//更新今日已学
		if err = gredis.INCR(common.TODAY_PREFIX+openId, util.GetNextDayBegin()); err != nil {
			logrus.Error("incr error :", err)
		}
		//3.往用户做过题表插入数据
		if err = models.AddDoneTopic(models.DoneTopic{
			OpenId:  openId,
			TopicId: topicId,
		}); err != nil {
			logrus.Error("AddDoneTopic error :", err)
			return nil, err
		}
	} else {
		right = false
		topic.WrongNum = topic.WrongNum + 1
		//往用户错题表插入数据
		if err = models.AddWrongTopic(models.WrongTopic{
			OpenId:  openId,
			TopicId: topicId,
		}); err != nil {
			logrus.Error("AddWrongTopic error :", err)
			return nil, err
		}
	}
	if err = models.UpdateTopic(topic); err != nil {
		logrus.Error("UpdateTopic error :", err)
		return nil, err
	}
	//3-1 存一个今日已答题
	if err = gredis.SAdd(common.TODAY_FINISH_PREFIX+openId, string(topicId)); err != nil {
		logrus.Error("AddDoneTopic sadd error ", err)
		return nil, err
	}
	//再设一个过期时间
	if err = gredis.ExpireAt(common.TODAY_FINISH_PREFIX+openId, util.GetNextDayBegin()); err != nil {
		logrus.Error("AddDoneTopic expireAt error ", err)
		return nil, err
	}

	//4.返回
	answerResp := &AnswerResp{
		MyAnswer: myAnswer,
		Right:    right,
	}
	answerResp.Answer = topic.Answer
	answerResp.TopicAnalysis = topic.TopicAnalysis
	answerResp.WrongNum = topic.WrongNum
	answerResp.RightNum = topic.RightNum
	return answerResp, nil

}
func GetAnalysis(openId string, topicId int) (*Analysis, error) {
	//1.看下答案对不对
	topic, err := models.GetTopic(topicId)
	if err != nil {
		logrus.Error("GetTopic error :", err)
		return nil, err
	}
	if topic == nil {
		logrus.Error("GetTopic empty :", topicId)
		return nil, err
	}
	anlysis := &Analysis{}
	anlysis.Answer = topic.Answer
	anlysis.TopicAnalysis = topic.TopicAnalysis
	anlysis.WrongNum = topic.WrongNum
	anlysis.RightNum = topic.RightNum
	return anlysis, nil
}
func Collect(openId string, topicId int) error {
	if err := models.AddCollect(models.Collect{
		OpenId:  openId,
		TopicId: topicId,
	}); err != nil {
		logrus.Error("Collect error :", err)
		return err
	}
	return nil
}
