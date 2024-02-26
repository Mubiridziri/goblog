package articles

import (
	"goblog/internal/entity"
	"goblog/internal/usecase/topics"
	"goblog/internal/usecase/users"
	"time"
)

type Article struct {
	ID        int
	Title     string
	Content   string
	TopicID   int
	Topic     topics.Topic
	Tags      string
	AuthorID  int
	Author    users.User
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateArticle struct {
	Title    string `json:"title" form:"title" binding:"required"`
	Content  string `json:"content" form:"content" binding:"required"`
	TopicID  int    `json:"topic_id" form:"topic_id" binding:"required"`
	Tags     string `json:"tags" form:"tags" binding:"required"`
	AuthorID int
}

type Repository interface {
	CreateArticle(article *entity.Article) error
	UpdateArticle(article *entity.Article) error
	GetArticleById(id int) (entity.Article, error)
	ListArticle(page, limit int) ([]entity.Article, error)
	GetArticlesCount() int64
}

type Controller struct {
	Repository
}

func NewController(repo Repository) *Controller {
	return &Controller{Repository: repo}
}

func (c Controller) CreateArticle(input CreateArticle) (Article, error) {
	articleDB := entity.Article{
		Title:    input.Title,
		Content:  input.Content,
		TopicID:  input.TopicID,
		Tags:     input.Tags,
		AuthorID: input.AuthorID,
	}

	err := c.Repository.CreateArticle(&articleDB)

	if err != nil {
		return Article{}, err
	}

	return fromDBArticle(&articleDB), nil
}

func (c Controller) GetArticleById(id int) (Article, error) {
	article, err := c.Repository.GetArticleById(id)

	if err != nil {
		return Article{}, err
	}

	return fromDBArticle(&article), nil
}

func fromDBArticle(article *entity.Article) Article {
	return Article{
		ID:        article.ID,
		Title:     article.Title,
		Content:   article.Content,
		TopicID:   article.TopicID,
		Topic:     topics.FromDBTopic(&article.Topic),
		Tags:      article.Tags,
		AuthorID:  article.AuthorID,
		Author:    users.FromDBUser(&article.Author),
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
	}
}
