package models

type Setting struct {
	Model
	OpenId         string `json:"open_id"`
	Region         string `json:"region"`
	ExamType       string `json:"exam_type"`
	ElementTypeOne string `json:"element_type_one"`
	ElementTypeTwo string `json:"element_type_two"`
	DailyNeedNum   int    `json:"daily_need_num"`
}

func AddSetting(setting Setting) error {
	if err := db.Create(&setting).Error; err != nil {
		return err
	}
	return nil
}
func SelectSettingByOpenId(openId string) (*Setting, error) {
	var setting Setting
	err := db.Where("open_id = ? AND flag = 0", openId).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}
func UpdateSetting(setting *Setting) error {
	data := make(map[string]interface{})

	if setting.Region != "" {
		data["region"] = setting.Region
	}
	if setting.ExamType != "" {
		data["exam_type"] = setting.ExamType
	}
	if setting.ElementTypeOne != "" {
		data["element_type_one"] = setting.ElementTypeOne
	}
	if setting.ElementTypeTwo != "" {
		data["element_type_two"] = setting.ElementTypeTwo
	}
	if setting.DailyNeedNum != 0 {
		data["daily_need_num"] = setting.DailyNeedNum
	}
	if err := db.Model(&Setting{}).Where("open_id = ? AND flag = 0 ", setting.OpenId).Updates(data).Error; err != nil {
		return err
	}
	return nil
}
