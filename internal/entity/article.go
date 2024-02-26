package entity

import (
	"gorm.io/gorm"
	"time"
)

type Article struct {
	ID        int `gorm:"primary_key"`
	Title     string
	Content   string `gorm:"type:text"`
	Tags      string
	IsDraft   bool
	TopicID   int
	Topic     Topic
	AuthorID  int
	Author    User
	CreatedAt time.Time
	UpdatedAt time.Time
}

type articleRepository struct {
	db *gorm.DB
}

func (r articleRepository) CreateArticle(article *Article) error {
	return r.db.Create(article).Error
}

func (r articleRepository) UpdateArticle(article *Article) error {
	return r.db.Save(article).Error
}

func (r articleRepository) GetArticleById(id int) (Article, error) {
	var article Article
	if err := r.db.Where(Article{ID: id}).Preload("Author").Preload("Topic").First(&article).Error; err != nil {
		return Article{}, err
	}

	return article, nil
}

func (r articleRepository) ListArticle(page, limit int) ([]Article, error) {
	var articles []Article
	offset := (page - 1) * limit
	if err := r.db.Offset(offset).Limit(limit).Preload("Author").Preload("Topic").Order("id desc").Find(&articles).Error; err != nil {
		return []Article{}, err
	}

	return articles, nil
}

func (r articleRepository) ListArticlesByTopic(page, limit, topicID int) ([]Article, error) {
	var articles []Article
	offset := (page - 1) * limit
	if err := r.db.
		Offset(offset).
		Limit(limit).
		Preload("Author").
		Preload("Topic").
		Order("id desc").
		Where(&Article{TopicID: topicID}).
		Find(&articles).Error; err != nil {
		return []Article{}, err
	}

	return articles, nil
}

func (r articleRepository) ListArticlesByYear(page, limit int, year string) ([]Article, error) {
	var articles []Article
	offset := (page - 1) * limit
	if err := r.db.
		Offset(offset).
		Limit(limit).
		Preload("Author").
		Preload("Topic").
		Order("id desc").
		Where("to_char(created_at, 'YYYY') = ?", year).
		Find(&articles).Error; err != nil {
		return []Article{}, err
	}

	return articles, nil
}

func (r articleRepository) ListArchiveYears() []string {
	var years []string
	if err := r.db.Model(&Article{}).Select("to_char(created_at, 'YYYY') as years").Group("to_char(created_at,'YYYY')").Scan(&years).Error; err != nil {
		return []string{}
	}

	return years
}

func (r articleRepository) GetArticlesCount() int64 {
	var count int64
	r.db.Model(&Article{}).Count(&count)
	return count
}

func (r articleRepository) GetArticlesByTopicCount(topicID int) int64 {
	var count int64
	r.db.Model(&Article{}).Where(&Article{TopicID: topicID}).Count(&count)
	return count
}

func (r articleRepository) GetArticlesByYearCount(year string) int64 {
	var count int64
	r.db.Model(&Article{}).Where("to_char(created_at, 'YYYY') = ?", year).Count(&count)
	return count
}
