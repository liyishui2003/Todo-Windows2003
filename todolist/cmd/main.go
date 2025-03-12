package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"todolist/internal/config"
	"todolist/internal/controller"
	"todolist/internal/middleware"
	"todolist/internal/model"
	"todolist/internal/repository"
	"todolist/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 初始化数据库连接
func initDB(cfg config.Config) *gorm.DB {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	// 自动迁移模型
	db.AutoMigrate(&model.User{}, &model.Task{})
	return db
}

func initLogger() *zap.Logger {
	// 打开日志文件
	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	// 配置 Zap 日志核心
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), // JSON 格式
		zapcore.AddSync(file), // 输出到文件
		zapcore.DebugLevel,    // 日志级别
	)

	// 创建日志记录器
	logger := zap.New(core)
	return logger
}

func main() {
	//初始化日志
	logger := initLogger()
	zap.ReplaceGlobals(logger) // 替换全局日志记录器
	defer logger.Sync()        // 在程序结束时刷新日志缓冲区
	defer func() {
		if r := recover(); r != nil {
			zap.L().Panic("捕获到 panic", zap.Any("panic", r))
		}
	}()

	//读取gin配置
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("加载配置失败")
	}
	config.Cfg = cfg

	// 初始化数据库
	db := initDB(config.Cfg)
	fmt.Println("数据库连接成功:", db)
	model.DB = db

	// 初始化 Repository、Service 和 Handler
	taskRepo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := controller.NewTaskHandler(taskService)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	r := gin.Default()

	// 临时测试路由
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// 注册各种路由
	//公开路由，无需jwt令牌验证
	r.POST("/login", userController.Login)
	r.POST("/register", userController.Register)

	//受保护的 需要中间件验证
	// tasks
	r.POST("/tasks", middleware.AuthMiddleware(), taskHandler.CreateTask)
	r.GET("/tasks/allTask/:userID", middleware.AuthMiddleware(), taskHandler.GetTasksByUserID)
	r.GET("/tasks/:taskID", middleware.AuthMiddleware(), taskHandler.GetTask)
	r.GET("/tasks", middleware.AuthMiddleware(), taskHandler.GetTasksByPriority)
	r.PUT("/tasks/:taskID", middleware.AuthMiddleware(), taskHandler.EditTask)
	r.DELETE("/tasks/:taskID", middleware.AuthMiddleware(), taskHandler.DeleteTask)

	// users
	r.GET("/users/:username", middleware.AuthMiddleware(), userController.FindUserID)

	//路由分组，暂定
	protected := r.Group("/protected")
	protected.Use(middleware.AuthMiddleware()) // 使用 JWT 认证中间件
	{
		protected.GET("/info", func(ctx *gin.Context) {
			userID := ctx.MustGet("user_id").(float64)
			ctx.JSON(http.StatusOK, gin.H{"user_id": userID})
		})
	}
	// 启动服务
	r.Run(":8080")
}
