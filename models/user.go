package models

type User struct {
	Model
	UserName     string `json:"user_name"`
	OpenId       string `json:"open_id"`
	HasLearnNum  int    `json:"has_learn_num"`
	PracticeDays int    `json:"practice_days"`
	PracticeTime int    `json:"practice_time"`
	AnswerNum    int    `json:"answer_num"`
	WrongNum     int    `json:"wrong_num"`
}

func AddUser(user User) error {
	if err := db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}
func SelectUserById(id int) (*User, error) {
	var user User
	err := db.Where("id = ? AND flag = 0", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func SelectUserByOpenId(openId string) (*User, error) {
	var user User
	err := db.Where("open_id = ? AND flag = 0", openId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func CountUser(region, examType string) (int, error) {
	var count int
	if err := db.Model(&User{}).Where("region = ? AND exam_type = ? AND flag = 0", region, examType).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
func UpdateUser(user *User) error {
	data := make(map[string]interface{})
	//更新已学题目，每日需刷题数量，答题时长，练习天数，创建练习次数，答题量
	if user.HasLearnNum > 0 {
		data["has_learn_num"] = user.HasLearnNum
	}
	if user.PracticeDays > 0 {
		data["practice_days"] = user.PracticeDays
	}
	if user.PracticeTime > 0 {
		data["practice_time"] = user.PracticeTime
	}
	if user.AnswerNum > 0 {
		data["answer_num"] = user.AnswerNum
	}
	if err := db.Model(&User{}).Where("open_id = ? AND flag = 0 ", user.OpenId).Updates(data).Error; err != nil {
		return err
	}
	return nil
}
