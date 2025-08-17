package api

import (
	"project/internal/api/middleware"
	"project/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type HandlerDependencies struct {
	AuthService service.AuthService
	LogService  service.LogService
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
	}

	return router
}
