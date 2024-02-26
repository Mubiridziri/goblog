package server

import (
	"github.com/gin-gonic/gin"
	"goblog/internal/usecase/articles"
	"goblog/internal/usecase/users"
	"net/http"
	"strconv"
)

func (s *Server) NewArticlePage(c *gin.Context) {
	user, exists := c.Get("user")

	if !exists { //TODO middleware?
		c.Redirect(http.StatusUnauthorized, "/")
		return
	}

	topicsList, err := s.topicController.ListTopics(1, 10) //ajax instead this

	if err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	if c.Request.Method == "POST" {
		//Handle form

		var form articles.CreateArticle
		if bindErr := c.ShouldBind(&form); bindErr != nil {
			c.HTML(http.StatusOK, "editor.tmpl", gin.H{
				"user":   user,
				"topics": topicsList.Entries,
				"error":  bindErr.Error(),
				"form":   form,
			})
			return
		}
		form.AuthorID = user.(users.User).ID
		article, createErr := s.articleController.CreateArticle(form)

		if createErr != nil {
			c.HTML(http.StatusBadRequest, "editor.tmpl", gin.H{
				"user":   user,
				"topics": topicsList.Entries,
				"error":  createErr.Error(),
				"form":   form,
			})
			return
		}

		c.HTML(http.StatusOK, "editor.tmpl", gin.H{
			"user":    user,
			"topics":  topicsList.Entries,
			"article": article,
			"form":    form,
		})
		return
	}

	c.HTML(http.StatusOK, "editor.tmpl", gin.H{
		"user":   user,
		"topics": topicsList.Entries,
	})
}

func (s *Server) ViewArticlePage(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		//Render Not Found
		c.Redirect(http.StatusNotFound, "/")
		return
	}

	article, err := s.articleController.GetArticleById(id)

	if err != nil {
		//Render Not Found
		c.Redirect(http.StatusNotFound, "/")
		return
	}

	user, _ := c.Get("user")

	topicsList, err := s.topicController.ListTopics(1, 10)

	c.HTML(http.StatusOK, "article.tmpl", gin.H{
		"user":    user,
		"article": article,
		"topics":  topicsList.Entries,
	})

}

func (s *Server) AddArticlesPagesRoutes(g *gin.RouterGroup) {
	grp := g.Group("/articles")
	grp.GET("/new", s.NewArticlePage)
	grp.POST("/new", s.NewArticlePage)

	grp.GET("/:id", s.ViewArticlePage)

	//grp.PUT("/:id", s.handleUpdateUser)
	//grp.GET("/:id", s.handleDetailUser)
	//grp.DELETE("/:id", s.handleDeleteUser)
}
