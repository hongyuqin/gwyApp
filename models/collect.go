package models

type Collect struct {
	Model
	OpenId  string `json:"open_id"`
	TopicId int    `json:"topic_id"`
}

func AddCollect(collect Collect) error {
	//假如存在 就不插入，直接返回nil
	var count int
	err := db.Model(&Collect{}).Where("open_id = ? AND topic_id = ? ", collect.OpenId, collect.TopicId).Count(&count).Error
	if err != nil {
		return err
	}
	if count == 1 {
		return nil
	}
	if err := db.Create(&collect).Error; err != nil {
		return err
	}
	return nil
}

func DelCollect(topicId int, openId string) error {
	if err := db.Where("topic_id = ? AND open_id = ?", topicId, openId).Delete(&Collect{}).Error; err != nil {
		return err
	}
	return nil
}

func GetCollect(openId string) (*Collect, error) {
	var (
		collect Collect
		err     error
	)
	data := make(map[string]interface{})
	data["open_id"] = openId

	err = db.Where(data).First(&collect).Error
	if err != nil {
		return nil, err
	}
	return &collect, nil
}
func GetCollects(openId string) ([]Collect, error) {
	var (
		collects []Collect
		err      error
	)
	data := make(map[string]interface{})
	data["open_id"] = openId

	err = db.Where(data).Find(&collects).Error
	if err != nil {
		return nil, err
	}
	return collects, nil
}
func ExistCollect(openId string, topicId int) (bool, error) {
	count := 0
	data := make(map[string]interface{})
	data["open_id"] = openId
	data["topic_id"] = topicId
	db.Model(&Collect{}).Where(data).Count(&count)
	if count == 1 {
		return true, nil
	}
	return false, nil
}
