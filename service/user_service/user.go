package user_service

import (
	"github.com/sirupsen/logrus"
	"gwyApp/common"
	"gwyApp/models"
	"gwyApp/pkg/gredis"
	"strconv"
)

/**
		"rank":10,//排名
     	"total":100,//总人数
    	"left_days":,//考试剩余天数
    	"daily_need_num":,//每日需刷题数量
    	"today_practice_num":,//今日刷题数量 存redis
    	"has_learn_num":,//已学数量
    	"wrong_num":10,//做错数量 需复习数量
*/
type HomeDetail struct {
	LeftDays         int    `json:"left_days"`
	TodayPracticeNum int    `json:"today_practice_num"`
	DailyNeedNum     int    `json:"daily_need_num"`
	HasLearnNum      int    `json:"has_learn_num"`
	WrongNum         int    `json:"wrong_num"`
	TotalTopicNum    int    `json:"total_topic_num"`
	Region           string `json:"region"`
	ExamType         string `json:"exam_type"`
	ElementTypeOne   string `json:"element_type_one"`
	ElementTypeTwo   string `json:"element_type_two"`
}

func Home(openId string) (*HomeDetail, error) {
	//0.获取用户基本信息
	user, err := models.SelectUserByOpenId(openId)
	if err != nil {
		logrus.Error("select user error :", err)
		return nil, err
	}
	//TODO:3.获取剩余天数
	//4.获取今日做题数
	todayNumStr, err := gredis.Get(common.TODAY_PREFIX + openId)
	if err != nil {
		logrus.Error("redis Get error :", err)
		todayNumStr = []byte("0")
	}
	todayNum, err := strconv.Atoi(string(todayNumStr))
	if err != nil {
		logrus.Error("Atoi error :", err)
		return nil, err
	}
	//5.今日错题数
	topics, err := models.GetWrongTopics(openId)
	if err != nil {
		logrus.Error("GetWrongTopics error :", err)
		return nil, err
	}
	setting, err := models.SelectSettingByOpenId(openId)
	if err != nil {
		logrus.Error("get setting error")
		return nil, err
	}
	homeDetail := &HomeDetail{
		TodayPracticeNum: todayNum,
		DailyNeedNum:     setting.DailyNeedNum,
		HasLearnNum:      user.HasLearnNum,
		WrongNum:         len(topics),
		Region:           setting.Region,
		TotalTopicNum:    1000,
	}
	return homeDetail, nil
}

type PlanReq struct {
	AccessToken    string `schema:"accessToken"`
	Region         string `schema:"region"`
	ExamType       string `schema:"exam_type"`
	DailyNeedNum   int    `schema:"daily_need_num"`
	ElementTypeOne string `schema:"element_type_one"`
	ElementTypeTwo string `schema:"element_type_two"`
}

func Plan(req *PlanReq) error {
	if err := models.UpdateSetting(&models.Setting{
		OpenId:         req.AccessToken,
		Region:         req.Region,
		ExamType:       req.ExamType,
		DailyNeedNum:   req.DailyNeedNum,
		ElementTypeOne: req.ElementTypeOne,
		ElementTypeTwo: req.ElementTypeTwo,
	}); err != nil {
		logrus.Error("UpdateSetting error :", err)
		return err
	}
	return nil
}
