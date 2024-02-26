package entity

import "gorm.io/gorm"

type Topic struct {
	ID    int    `gorm:"primary_key"`
	Title string `gorm:"unique"`
}

type topicRepository struct {
	db *gorm.DB
}

func (r topicRepository) GetTopicById(id int) (Topic, error) {
	var topic Topic
	if err := r.db.Where(Topic{ID: id}).First(&topic).Error; err != nil {
		return Topic{}, err
	}

	return topic, nil
}

func (r topicRepository) ListTopics(page, limit int) ([]Topic, error) {
	var topics []Topic
	offset := (page - 1) * limit
	if err := r.db.Offset(offset).Limit(limit).Find(&topics).Error; err != nil {
		return []Topic{}, err
	}

	return topics, nil
}

func (r topicRepository) GetTopicsCount() int64 {
	var count int64
	r.db.Model(&Topic{}).Count(&count)
	return count
}
