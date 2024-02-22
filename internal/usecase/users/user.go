package users

import (
	"errors"
	"goblog/internal/entity"
	"time"
)

type UserLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUser struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password"  binding:"required"`
}

type Repository interface {
	CreateUser(user *entity.User) error
	UpdateUser(user *entity.User) error
	RemoveUser(user *entity.User) error
	ListUsers(page, limit int) ([]entity.User, error)
	GetUserById(id int) (entity.User, error)
	GetUserByUsername(username string) (entity.User, error)
	GetUsersCount() int64
}

type PaginatedUsersList struct {
	Total   int64  `json:"total"`
	Entries []User `json:"entries"`
}

type Controller struct {
	Repository
}

func NewController(repo Repository) *Controller {
	return &Controller{Repository: repo}
}

func (c Controller) LoginUser(input UserLogin) (User, error) {
	user, err := c.Repository.GetUserByUsername(input.Username)

	if err != nil {
		return User{}, err
	}

	if !user.IsPasswordCorrect(input.Password) {
		return User{}, errors.New("invalid credentials")
	}

	return fromDBUser(&user), nil
}

func (c Controller) CreateUser(input CreateUser) (User, error) {
	user := entity.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Username:  input.Username,
		Password:  input.Password,
		Email:     input.Email,
	}
	if err := c.Repository.CreateUser(&user); err != nil {
		return User{}, err
	}

	return fromDBUser(&user), nil
}

func (c Controller) UpdateUser(id int, input CreateUser) (User, error) {

	user, err := c.Repository.GetUserById(id)

	if err != nil {
		return User{}, err
	}

	user.FirstName = input.FirstName
	user.LastName = input.LastName
	user.Username = input.Username
	user.Email = input.Email

	if input.Password != "" {
		user.Password = input.Password
	}

	if err := c.Repository.UpdateUser(&user); err != nil {
		return User{}, err
	}

	return fromDBUser(&user), nil
}

func (c Controller) RemoveUser(id int) (User, error) {

	user, err := c.Repository.GetUserById(id)

	if err != nil {
		return User{}, err
	}

	if err := c.Repository.RemoveUser(&user); err != nil {
		return User{}, err
	}

	return fromDBUser(&user), nil
}

func (c Controller) ListUsers(page, limit int) (PaginatedUsersList, error) {

	users, err := c.Repository.ListUsers(page, limit)

	if err != nil {
		return PaginatedUsersList{}, err
	}

	var userViews []User

	for _, user := range users {
		userViews = append(userViews, fromDBUser(&user))
	}

	return PaginatedUsersList{
		Total:   c.Repository.GetUsersCount(),
		Entries: userViews,
	}, nil
}

func (c Controller) GetUserById(id int) (User, error) {

	user, err := c.Repository.GetUserById(id)

	if err != nil {
		return User{}, err
	}

	return fromDBUser(&user), nil
}

func (c Controller) GetUserByUsername(username string) (User, error) {
	user, err := c.Repository.GetUserByUsername(username)

	if err != nil {
		return User{}, err
	}

	return fromDBUser(&user), nil
}

func fromDBUser(user *entity.User) User {
	return User{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
