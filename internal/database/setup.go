package database

import (
	"errors"
	"goblog/internal/config"
	"goblog/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectAndMigrate(cfg *config.Config) (*gorm.DB, error) {
	database, err := gorm.Open(postgres.Open(cfg.Database.DSN))
	if err != nil {
		return nil, errors.New("failed connect to database")
	}

	// Auto Migrate
	err = database.AutoMigrate(entity.User{})
	err = database.AutoMigrate(entity.Topic{})
	err = database.AutoMigrate(entity.Article{})

	if err != nil {
		return nil, errors.New("failed auto migrate database")
	}

	return database, nil
}
