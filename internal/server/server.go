package server

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"goblog/docs"
	"goblog/internal/usecase/users"
	"net/http"
	"strconv"
)

const UserKey = "AUTH"

// TODO move to env APP_SECRET
var secret = []byte("RHYaxoa6iqb1VTCsFtdM2PAAu8i8CYhU")

type Config struct {
	UserController *users.Controller
}

type Server struct {
	Router         *gin.Engine
	userController *users.Controller
}

type ListQuery struct {
	Page     int  `form:"page"`
	Limit    int  `form:"limit"`
	Simplify bool `form:"simplify"`
}

func New(config Config) *Server {
	s := Server{
		Router:         gin.New(),
		userController: config.UserController,
	}

	s.registerRoutes()
	return &s
}

func (s *Server) registerRoutes() {
	//Middleware
	s.Router.Use(gin.Logger())
	s.Router.Use(gin.Recovery())
	s.Router.Use(sessions.Sessions(UserKey, cookie.NewStore(secret)))

	//Swagger
	configureSwagger(s.Router)

	// K8s probe
	s.Router.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	//API /api/v1
	mainGroup := s.Router.Group("/api/v1")
	mainGroup.POST("/login", s.handleLogin)

	mainGroup.Use(AuthRequired(s.userController))
	{
		s.AddUserRoutes(mainGroup)

		mainGroup.GET("/login", s.handleProfile)
		mainGroup.GET("/logout", s.handleLogout)
	}
}

func configureSwagger(r *gin.Engine) {
	//Swagger
	docs.SwaggerInfo.Title = "DMAAS"
	docs.SwaggerInfo.Description = "Data management and analytic system"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func AuthRequired(controller *users.Controller) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userKey := session.Get(UserKey)

		if userKey == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		authUser, err := controller.GetUserByUsername(userKey.(string))

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		c.Set("user", authUser)
		c.Next()
	}

}

func (s *Server) getListQuery(c *gin.Context) (ListQuery, error) {
	pageQuery := c.Query("page")
	limitQuery := c.Query("limit")
	simplifyQuery := c.Query("simplify")

	if simplifyQuery == "" {
		simplifyQuery = "false"
	}

	page, err := strconv.Atoi(pageQuery)
	limit, err := strconv.Atoi(limitQuery)
	simplify, err := strconv.ParseBool(simplifyQuery)

	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 10
	}

	if err != nil {
		return ListQuery{}, err
	}

	query := ListQuery{
		Page:     page,
		Limit:    limit,
		Simplify: simplify,
	}
	return query, nil
}
