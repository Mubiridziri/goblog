package topics

import "goblog/internal/entity"

type Topic struct {
	ID    int
	Title string
}

type Repository interface {
	ListTopics(page, limit int) ([]entity.Topic, error)
	GetTopicById(id int) (entity.Topic, error)
	GetTopicsCount() int64
}

type PaginatedTopicsList struct {
	Total   int64   `json:"total"`
	Entries []Topic `json:"entries"`
}

type Controller struct {
	Repository
}

func NewController(repo Repository) *Controller {
	return &Controller{Repository: repo}
}

func (c Controller) ListTopics(page, limit int) (PaginatedTopicsList, error) {

	topicsDB, err := c.Repository.ListTopics(page, limit)

	if err != nil {
		return PaginatedTopicsList{}, err
	}

	var topics []Topic

	for _, dbtopic := range topicsDB {
		topics = append(topics, FromDBTopic(&dbtopic))
	}

	return PaginatedTopicsList{
		Total:   c.Repository.GetTopicsCount(),
		Entries: topics,
	}, nil
}

func (c Controller) GetTopicById(id int) (Topic, error) {

	topicDB, err := c.Repository.GetTopicById(id)

	if err != nil {
		return Topic{}, err
	}

	return FromDBTopic(&topicDB), nil
}

func FromDBTopic(topicDB *entity.Topic) Topic {
	return Topic{
		ID:    topicDB.ID,
		Title: topicDB.Title,
	}
}
