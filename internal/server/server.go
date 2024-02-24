package server

import (
	"fmt"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"goblog/docs"
	"goblog/internal/usecase/users"
	"net/http"
	"path/filepath"
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

	//Server Side Rendering
	//s.Router.LoadHTMLGlob("web/template/**/*")
	s.Router.HTMLRender = loadTemplates("web/template/")
	s.Router.Static("/static", "web/static/")
	ui := s.Router.Group("")
	ui.Use(LoadUserIfExists(s.userController))
	{
		ui.GET("", s.renderHomePage)
		ui.GET("/login", s.loginPage)
		ui.POST("/login", s.loginPage)
		ui.GET("/register", s.registerPage)
		ui.POST("/register", s.registerPage)
		ui.GET("/u/:username", s.profilePage)
	}

	s.AddUserPagesRoutes(ui)

	//API /api/v1
	api := s.Router.Group("/api/v1")
	api.POST("/login", s.handleLogin)

	api.Use(AuthRequired(s.userController))
	{
		s.AddUserAPIRoutes(api)

		api.GET("/login", s.handleProfile)
		api.GET("/logout", s.handleLogout)
	}

	//Swagger
	configureSwagger(s.Router)
	// K8s probe
	s.Router.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })
}

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	rootPages, err := filepath.Glob(templatesDir + "/*.tmpl")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Loaded templates: ")

	for _, page := range rootPages {
		path := filepath.Base(page)
		fmt.Println("\t - " + path)
		r.AddFromFiles(path, page)
	}

	layouts, err := filepath.Glob(templatesDir + "/layouts/*.tmpl")
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob(templatesDir + "/includes/*.tmpl")
	if err != nil {
		panic(err.Error())
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, include := range includes {
		path := filepath.Base(include)
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		fmt.Println("\t - " + path)
		r.AddFromFiles(path, files...)
	}
	fmt.Println("")
	return r
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

func LoadUserIfExists(controller *users.Controller) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userKey := session.Get(UserKey)

		if userKey == nil {
			c.Next()
			return
		}

		authUser, err := controller.GetUserByUsername(userKey.(string))

		if err != nil {
			c.Next()
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

func (s *Server) renderHomePage(c *gin.Context) {
	user, exists := c.Get("user")

	params := gin.H{}

	if exists {
		params["user"] = user
	}

	c.HTML(http.StatusOK, "home.tmpl", params)
}
