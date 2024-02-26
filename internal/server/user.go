package server

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"goblog/internal/usecase/users"
	"net/http"
)

func (s *Server) profilePage(c *gin.Context) {
	username := c.Param("username")

	if username == "" {
		c.Redirect(http.StatusFound, "/")
		return
	}

	user, err := s.userController.GetUserByUsername(username)

	if err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	c.HTML(http.StatusOK, "profile.tmpl", gin.H{
		"user": user,
	})
}

func (s *Server) loginPage(c *gin.Context) {
	_, exists := c.Get("user")

	if exists {
		c.Redirect(http.StatusFound, "/")
		return
	}

	if c.Request.Method == "POST" {
		//Handle form
		var form users.UserLogin

		if err := c.ShouldBind(&form); err != nil {
			c.HTML(http.StatusBadRequest, "login.tmpl", gin.H{
				"error": err.Error(),
			})
			return
		}

		user, err := s.userController.LoginUser(form)

		if err != nil {
			c.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{
				"error": "Неверный логин или пароль",
			})
			return
		}

		session := sessions.Default(c)
		session.Set(UserKey, user.Username)
		err = session.Save()

		if err != nil {
			c.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{
				"error": "Ошибка сохранения куков",
			})
			return
		}

		c.Redirect(http.StatusFound, "/")
		return
	}
	c.HTML(http.StatusOK, "login.tmpl", gin.H{})
}

func (s *Server) registerPage(c *gin.Context) {
	_, exists := c.Get("user")

	if exists {
		c.Redirect(http.StatusFound, "/")
		return
	}

	if c.Request.Method == "POST" {
		//Handle form
		var form users.CreateUser

		if err := c.ShouldBind(&form); err != nil {
			c.HTML(http.StatusBadRequest, "register.tmpl", gin.H{
				"error": err.Error(),
			})
			return
		}

		_, err := s.userController.CreateUser(form)

		if err != nil {
			c.HTML(http.StatusInternalServerError, "register.tmpl", gin.H{
				"error": err.Error(),
			})
			return
		}

		_, err = s.userController.LoginUser(users.UserLogin{
			Username: form.Username,
			Password: form.Password,
		})

		if err != nil {
			c.HTML(http.StatusUnauthorized, "register.tmpl", gin.H{
				"error": err.Error(),
			})
			return
		}

		c.Redirect(http.StatusFound, "/")
		return
	}
	c.HTML(http.StatusOK, "register.tmpl", gin.H{})
}
