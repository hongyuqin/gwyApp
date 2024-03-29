package models

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

type Topic struct {
	Model
	TopicNo        string `json:"topic_no"`
	TopicName      string `json:"topic_name"`
	Photo          string `json:"photo"`
	OptionA        string `json:"option_a"`
	OptionB        string `json:"option_b"`
	OptionC        string `json:"option_c"`
	OptionD        string `json:"option_d"`
	Answer         string `json:"answer"`
	TopicAnalysis  string `json:"topic_analysis"`
	WrongNum       int    `json:"wrong_num"`
	RightNum       int    `json:"right_num"`
	Region         string `json:"region"`
	Year           int    `json:"year"`
	ExamType       string `json:"exam_type"`
	ElementTypeOne string `json:"element_type_one"`
	ElementTypeTwo string `json:"element_type_two"`
}

func (this *Topic) String() string {
	bytes, err := json.Marshal(this)
	if err != nil {
		return "marshal error"
	}
	return string(bytes)
}
func AddTopic(topic *Topic) error {
	if err := db.Create(&topic).Error; err != nil {
		return err
	}
	return nil
}

func GetTopic(id int) (*Topic, error) {
	var topic Topic
	if err := db.Where("id = ? AND flag = 0", id).First(&topic).Error; err != nil {
		return nil, err
	}
	return &topic, nil
}

func GetTopics(topic *Topic) ([]Topic, error) {
	var (
		topics []Topic
		err    error
	)
	data := make(map[string]interface{})
	/*if topic.Region != "" {
		data["region"] = topic.Region
	}*/
	if topic.ElementTypeOne != "" {
		data["element_type_one"] = topic.ElementTypeOne
	}
	if topic.ElementTypeTwo != "" {
		data["element_type_two"] = topic.ElementTypeTwo
	}

	err = db.Where(data).Find(&topics).Error
	if err != nil {
		return nil, err
	}
	return topics, nil
}

func GetTopicsId(topic *Topic) ([]int, error) {
	var topicIds []int
	data := make(map[string]interface{})
	if topic.ElementTypeOne != "" {
		data["element_type_one"] = topic.ElementTypeOne
	}
	if topic.ElementTypeTwo != "" {
		data["element_type_two"] = topic.ElementTypeTwo
	}
	db.Model(&Topic{}).Where(data).Pluck("id", &topicIds)
	return topicIds, nil
}

func UpdateTopic(topic *Topic) error {
	data := make(map[string]interface{})
	if topic.WrongNum > 0 {
		data["wrong_num"] = topic.WrongNum
	}
	if topic.RightNum > 0 {
		data["right_num"] = topic.RightNum
	}
	if err := db.Model(&Topic{}).Where("id = ? AND flag = 0 ", topic.ID).Updates(data).Error; err != nil {
		return err
	}
	return nil
}
func UpdateTopicByTopicNo(topic *Topic) error {
	logrus.Info("UpdateTopicByTopicNo  :", topic.TopicNo)
	data := make(map[string]interface{})
	if topic.Answer != "" {
		data["answer"] = topic.Answer
	}
	if topic.TopicAnalysis != "" {
		data["topic_analysis"] = topic.TopicAnalysis
	}
	if err := db.Model(&Topic{}).Where("topic_no = ? AND flag = 0 ", topic.TopicNo).Updates(data).Error; err != nil {
		return err
	}
	return nil
}
