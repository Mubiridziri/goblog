package users

import (
	"github.com/stretchr/testify/assert"
	"goblog/internal/entity"
	"testing"
	"time"
)

type mockUserRepository struct {
	user entity.User
}

func (r mockUserRepository) CreateUser(user *entity.User) error {
	return nil
}
func (r mockUserRepository) UpdateUser(user *entity.User) error {
	return nil
}
func (r mockUserRepository) RemoveUser(user *entity.User) error {
	return nil
}
func (r mockUserRepository) ListUsers(page, limit int) ([]entity.User, error) {
	return []entity.User{r.user}, nil
}
func (r mockUserRepository) GetUserById(id int) (entity.User, error) {
	return r.user, nil
}
func (r mockUserRepository) GetUserByUsername(username string) (entity.User, error) {
	return r.user, nil
}
func (r mockUserRepository) GetUsersCount() int64 {
	return 1
}

func TestController_LoginUser(t *testing.T) {
	a := assert.New(t)

	dbUser := entity.User{
		ID:        1,
		Username:  "username",
		FirstName: "FirstName",
		LastName:  "LastName",
		Email:     "username@email.com",
		Password:  "123456",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	_ = dbUser.HashPassword()

	repo := mockUserRepository{dbUser}
	controller := NewController(repo)

	input := UserLogin{
		Username: "username",
		Password: "123456",
	}

	user, err := controller.LoginUser(input)

	if err != nil {
		t.Fatal(err.Error())
	}

	a.Equal(user.ID, dbUser.ID, "id they should be equal")
	a.Equal(user.Username, dbUser.Username, "username they should be equal")
	a.Equal(user.FirstName, dbUser.FirstName, "firstname they should be equal")
	a.Equal(user.LastName, dbUser.LastName, "lastname they should be equal")
	a.Equal(user.Email, dbUser.Email, "email they should be equal")
}

func TestController_LoginUser2(t *testing.T) {
	dbUser := entity.User{
		ID:        1,
		Username:  "username",
		FirstName: "FirstName",
		LastName:  "LastName",
		Email:     "username@email.com",
		Password:  "123456",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	_ = dbUser.HashPassword()

	repo := mockUserRepository{dbUser}
	controller := NewController(repo)

	input := UserLogin{
		Username: "username",
		Password: "wrontpassword",
	}

	_, err := controller.LoginUser(input)

	if err == nil {
		t.Fatal("wrong password passed like valid")
	}
}

func TestController_CreateUser(t *testing.T) {
	a := assert.New(t)
	dbUser := entity.User{
		ID:        1,
		Username:  "username",
		FirstName: "FirstName",
		LastName:  "LastName",
		Email:     "username@email.com",
		Password:  "123456",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	repo := mockUserRepository{dbUser}
	controller := NewController(repo)

	input := CreateUser{
		FirstName: "FirstName",
		LastName:  "LastName",
		Username:  "username",
		Email:     "username@email.com",
		Password:  "123456",
	}

	user, err := controller.CreateUser(input)

	if err != nil {
		t.Fatal(err.Error())
	}

	a.Equal(user.Username, input.Username, "username they should be equal")
	a.Equal(user.FirstName, input.FirstName, "firstname they should be equal")
	a.Equal(user.LastName, input.LastName, "lastname they should be equal")
	a.Equal(user.Email, input.Email, "email they should be equal")

}

func TestController_UpdateUser(t *testing.T) {
	a := assert.New(t)
	dbUser := entity.User{
		ID:        1,
		Username:  "username",
		FirstName: "FirstName",
		LastName:  "LastName",
		Email:     "username@email.com",
		Password:  "123456",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	repo := mockUserRepository{dbUser}
	controller := NewController(repo)

	input := CreateUser{
		FirstName: "John",
		LastName:  "Smith",
		Username:  "john.smit",
		Email:     "john@smits.com",
		Password:  "123456",
	}

	user, err := controller.UpdateUser(1, input)

	if err != nil {
		t.Fatal(err.Error())
	}

	a.Equal(user.Username, input.Username, "username they should be equal")
	a.Equal(user.FirstName, input.FirstName, "firstname they should be equal")
	a.Equal(user.LastName, input.LastName, "lastname they should be equal")
	a.Equal(user.Email, input.Email, "email they should be equal")
}
