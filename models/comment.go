package models

type Comment struct {
	Model
	TopicId        int    `json:"topic_id"`
	OpenId         string `json:"open_id"`
	CommentContent string `json:"comment_content"`
}

func AddComment(comment Comment) error {
	if err := db.Create(&comment).Error; err != nil {
		return err
	}
	return nil
}

func DeleteComment(id int) error {
	data := make(map[string]interface{})
	data["id"] = id
	if err := db.Model(&Comment{}).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}
