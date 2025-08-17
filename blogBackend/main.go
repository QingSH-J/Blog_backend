package main

import (
	"fmt"
	"log"
	"project/internal/api"
	"project/internal/config"
	"project/internal/service"
	"project/internal/store"
)

func main() {
	// 加载配置
	config.LoadConfig()
	cfg := config.GetConfig()

	// database connection configuration
	db := store.NewDB(cfg.DB_DSN)

	// create stores and services
	userStore := store.NewUserStore(db)
	bookLogStore := store.NewBookLogStore(db)
	forumStore := store.NewForumStore(db)
	authService := service.NewAuthService(userStore)
	logService := service.NewLogService(bookLogStore)
	forumService := service.NewForumService(forumStore)
	// database migrations
	fmt.Println("Running database migrations...")
	if err := userStore.Migrate(); err != nil {
		log.Fatalf("Error migrating user table: %v", err)
	}
	fmt.Println("User table migration successful")

	if err := bookLogStore.Migrate(); err != nil {
		log.Fatalf("Error migrating book log table: %v", err)
	}
	fmt.Println("Book log table migration successful")

	// create API dependencies
	deps := api.HandlerDependencies{
		AuthService:  authService,
		LogService:   logService,
		ForumService: forumService,
	}

	// create router
	router := api.NewRouter(deps)

	// start server
	serverAddr := cfg.SERVER_HOST + ":" + cfg.SERVER_PORT
	fmt.Printf("\n🚀 Server starting on %s (%s environment)...\n", serverAddr, cfg.ENVIRONMENT)
	fmt.Println("📡 API endpoints available:")
	fmt.Println("   POST /api/v1/auth/register - 用户注册")
	fmt.Println("   POST /api/v1/auth/login    - 用户登录")
	fmt.Println("   POST /api/v1/new/          - 创建图书记录 (需要JWT认证)")
	fmt.Printf("\n🔐 JWT配置: Secret已设置, Token有效期: %s\n", cfg.JWT_EXPIRES_IN)
	fmt.Printf("📚 图书录入功能已启用，支持以下字段:\n")
	fmt.Printf("   - title, author, cover_url, status (必填)\n")
	fmt.Printf("   - my_rating, my_comment (可选)\n")
	fmt.Printf("🌐 CORS已启用，允许的源: %s\n", cfg.CORS_ALLOWED_ORIGINS)
	fmt.Printf("📝 日志级别: %s\n", cfg.LOG_LEVEL)
	fmt.Printf("\n💡 使用说明:\n")
	fmt.Printf("   1. 先注册/登录获取JWT token\n")
	fmt.Printf("   2. 在请求头添加: Authorization: Bearer <token>\n")
	fmt.Printf("   3. POST 图书数据到 /api/v1/new/ 端点\n")
	fmt.Println("⏳ Waiting for frontend requests...")

	if err := router.Run(":" + cfg.SERVER_PORT); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
