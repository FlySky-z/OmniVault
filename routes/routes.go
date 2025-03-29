package routes

import (
	docs "omnivault/docs"
	"omnivault/handlers"
	"omnivault/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 配置 Swagger 基本信息
	docs.SwaggerInfo.Title = "omnivault API"
	docs.SwaggerInfo.Description = "This is the API documentation for the omnivault project."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// 设置 Swagger 路由
	swag := r.Group("/swagger")
	swag.GET("/*any", ginswagger.WrapHandler(swaggerFiles.Handler))
	// 健康测试
	r.GET("/ping", handlers.PingHandler)

	// 设置authorize路由组
	authorize := r.Group("/authorize")
	{
		// 注册新用户
		authorize.POST("/register", handlers.RegisterHandler)
		authorize.POST("/login", handlers.LoginHandler)
		// 用户登录
		authorize.Use(middleware.AuthJWT())
		authorize.POST("/logout", handlers.LogoutHandler)
	}

	// 设置用户资源的 RESTful 路由
	users := r.Group("/users")
	{
		// users.GET("", handlers.ListUsers)       // 获取用户列表
		// users.POST("", handlers.CreateUser)    // 创建新用户
		users.GET("/:id", handlers.GetUser)    // 获取指定用户信息
		users.PUT("/:id", handlers.UpdateUser) // 更新指定用户信息
		// users.DELETE("/:id", handlers.DeleteUser) // 删除指定用户
	}

	// 设置文件资源的 RESTful 路由
	files := r.Group("/files")
	{
		// files.GET("", handlers.ListFiles)       // 获取文件列表
		files.POST("", handlers.UploadFile) // 上传新文件
		// files.GET("/:id", handlers.GetFile)    // 获取指定文件信息
		// files.PUT("/:id", handlers.UpdateFile) // 更新指定文件信息
		// files.DELETE("/:id", handlers.DeleteFile) // 删除指定文件
	}

	return r
}
