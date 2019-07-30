package models

type WrongTopic struct {
	Model
	OpenId  string `json:"open_id"`
	TopicId int    `json:"topic_id"`
}

func AddWrongTopic(wrongTopic WrongTopic) error {
	//假如存在 就不插入，直接返回nil
	var count int
	err := db.Model(&WrongTopic{}).Where("open_id = ? AND topic_id = ? ", wrongTopic.OpenId, wrongTopic.TopicId).Count(&count).Error
	if err != nil {
		return err
	}
	if count == 1 {
		return nil
	}
	if err := db.Create(&wrongTopic).Error; err != nil {
		return err
	}
	return nil
}

func DelWrongTopic(wrongTopic WrongTopic) error {
	if err := db.Where("open_id = ? AND topic_id = ?", wrongTopic.OpenId, wrongTopic.TopicId).Delete(&WrongTopic{}).Error; err != nil {
		return err
	}
	return nil
}

func GetWrongTopics(openId string) ([]WrongTopic, error) {
	var (
		topics []WrongTopic
		err    error
	)
	data := make(map[string]interface{})
	data["open_id"] = openId

	err = db.Where(data).Find(&topics).Error
	if err != nil {
		return nil, err
	}
	return topics, nil
}

func GetWrongTopic(openId string) (*WrongTopic, error) {
	var (
		topic WrongTopic
		err   error
	)
	data := make(map[string]interface{})
	data["open_id"] = openId

	err = db.Where(data).First(&topic).Error
	if err != nil {
		return nil, err
	}
	return &topic, nil
}
