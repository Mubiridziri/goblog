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
	ListArticlesByTopic(page, limit, topicID int) ([]entity.Article, error)
	ListArticlesByYear(page, limit int, year string) ([]entity.Article, error)
	GetArticlesCount() int64
	GetArticlesByTopicCount(topicID int) int64
	GetArticlesByYearCount(year string) int64
	ListArchiveYears() []string
}

type PaginatedArticleList struct {
	Total   int64     `json:"total"`
	Entries []Article `json:"entries"`
}

type Controller struct {
	Repository
}

func NewController(repo Repository) *Controller {
	return &Controller{Repository: repo}
}

func (c Controller) ListArchiveYears() []string {
	return c.Repository.ListArchiveYears()
}

func (c Controller) ListArticle(page, limit int) (PaginatedArticleList, error) {
	articleDB, err := c.Repository.ListArticle(page, limit)

	if err != nil {
		return PaginatedArticleList{}, err
	}

	var articles []Article

	for _, item := range articleDB {
		articles = append(articles, fromDBArticle(&item))
	}

	return PaginatedArticleList{
		Total:   c.Repository.GetArticlesCount(),
		Entries: articles,
	}, nil
}

func (c Controller) ListArticlesByTopic(page, limit, topicID int) (PaginatedArticleList, error) {
	articleDB, err := c.Repository.ListArticlesByTopic(page, limit, topicID)

	if err != nil {
		return PaginatedArticleList{}, err
	}

	var articles []Article

	for _, item := range articleDB {
		articles = append(articles, fromDBArticle(&item))
	}

	return PaginatedArticleList{
		Total:   c.Repository.GetArticlesByTopicCount(topicID),
		Entries: articles,
	}, nil
}

func (c Controller) ListArticlesByYear(page, limit int, year string) (PaginatedArticleList, error) {
	articleDB, err := c.Repository.ListArticlesByYear(page, limit, year)

	if err != nil {
		return PaginatedArticleList{}, err
	}

	var articles []Article

	for _, item := range articleDB {
		articles = append(articles, fromDBArticle(&item))
	}

	return PaginatedArticleList{
		Total:   c.Repository.GetArticlesByYearCount(year),
		Entries: articles,
	}, nil
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
