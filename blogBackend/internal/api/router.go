package api

import (
	"project/internal/api/middleware"
	"project/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type HandlerDependencies struct {
	AuthService  service.AuthService
	LogService   service.LogService
	ForumService service.ForumService
	ReadService  service.ReadService
}

func NewRouter(deps HandlerDependencies) *gin.Engine {
	router := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true

	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept"}

	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}

	config.ExposeHeaders = []string{"Content-Length"}

	config.AllowCredentials = true

	router.Use(cors.New(config))

	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	authHandler := NewAuthHandler(deps.AuthService)
	logHandler := NewLogHandler(deps.LogService)
	forumHandler := NewForumHandler(deps.ForumService)
	readHandler := NewReadHandler(deps.ReadService)
	apiV1 := router.Group("/api/v1")
	{
		authGroup := apiV1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
		}
		reviewGroup := apiV1.Group("/review")
		reviewGroup.Use(middleware.AuthMiddleware())
		{
			reviewGroup.GET("/books", logHandler.GetBookLog)
			reviewGroup.GET("/books/", logHandler.GetBookLog)
		}

		logGroup := apiV1.Group("/new")
		logGroup.Use(middleware.AuthMiddleware())
		{
			logGroup.POST("/", logHandler.CreateBookLog)
		}

		booksGroup := apiV1.Group("/books")
		booksGroup.Use(middleware.AuthMiddleware())
		{
			booksGroup.GET("/:id", logHandler.GetBook)
			booksGroup.PUT("/:id", logHandler.UpdateBookLog)
		}

		searchGroup := apiV1.Group("/search")
		searchGroup.GET("", logHandler.SearchBook)

		forumGroup := apiV1.Group("/forum")
		{
			forumGroup.GET("/topics", forumHandler.GetTopics)
			forumGroup.POST("/topics", middleware.AuthMiddleware(), forumHandler.CreateTopic)
			forumGroup.GET("/topics/:id", forumHandler.GetTopicByID)
			forumGroup.POST("/topics/:id/comments", middleware.AuthMiddleware(), forumHandler.CreateComment)
			forumGroup.GET("/topics/:id/comments", forumHandler.GetComments)
		}

		ReadGroup := apiV1.Group("/readtime")
		{
			ReadGroup.POST("", middleware.AuthMiddleware(), readHandler.CreateReadTime)
		}
	}

	return router
}
