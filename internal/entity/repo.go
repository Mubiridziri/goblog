package entity

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
	userRepository
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{
		db: db,
		userRepository: userRepository{
			db: db,
		},
	}
}
