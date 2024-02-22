package entity

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        int    `gorm:"primary_key"`
	Username  string `gorm:"unique"`
	FirstName string
	LastName  string
	Email     string `gorm:"unique"`
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type userRepository struct {
	db *gorm.DB
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	err = user.HashPassword()
	return err
}

func (user *User) BeforeUpdate(tx *gorm.DB) (err error) {
	if tx.Statement.Changed("Password") {
		err = user.HashPassword()
	}

	return err
}

func (user *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return err
	}

	user.Password = string(bytes)
	return nil
}

func (user *User) IsPasswordCorrect(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (r userRepository) GetUserById(id int) (User, error) {
	var user User
	if err := r.db.Where(User{ID: id}).First(&user).Error; err != nil {
		return User{}, err
	}

	return user, nil
}

func (r userRepository) GetUserByUsername(username string) (User, error) {
	var user User
	if err := r.db.Where(User{Username: username}).First(&user).Error; err != nil {
		return User{}, err
	}

	return user, nil
}

func (r userRepository) CreateUser(user *User) error {
	return r.db.Create(user).Error
}

func (r userRepository) UpdateUser(user *User) error {
	return r.db.Save(user).Error
}

func (r userRepository) RemoveUser(user *User) error {
	return r.db.Delete(user).Error
}

func (r userRepository) ListUsers(page, limit int) ([]User, error) {
	var users []User
	offset := (page - 1) * limit
	if err := r.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return []User{}, err
	}

	return users, nil
}

func (r userRepository) GetUsersCount() int64 {
	var count int64
	r.db.Model(&User{}).Count(&count)
	return count
}
