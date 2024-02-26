package app

import (
	"goblog/internal/config"
	"goblog/internal/database"
	"goblog/internal/entity"
	"goblog/internal/server"
	"goblog/internal/usecase/articles"
	"goblog/internal/usecase/topics"
	"goblog/internal/usecase/users"
	"net/http"
)

type Application struct {
	BindAddr string
}

func (a *Application) Run() error {
	configLoader := config.ConfigLoader{}

	cfg, err := configLoader.LoadConfig()
	if err != nil {
		return err
	}

	db, err := database.ConnectAndMigrate(cfg)

	if err != nil {
		return err
	}

	repo := entity.NewRepository(db)
	userController := users.NewController(repo)
	topicController := topics.NewController(repo)
	articleController := articles.NewController(repo)

	s := server.New(server.Config{
		UserController:    userController,
		TopicController:   topicController,
		ArticleController: articleController,
	})

	return http.ListenAndServe(a.BindAddr, s.Router)
}

func New(bindAddr string) *Application {
	return &Application{BindAddr: bindAddr}
}
