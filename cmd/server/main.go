package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"hmdp-backend/internal/config"
	"hmdp-backend/internal/data"
	"hmdp-backend/internal/middleware"
	"hmdp-backend/internal/router"
	"hmdp-backend/internal/service"
	"hmdp-backend/internal/utils"
	"hmdp-backend/pkg/logger"
)

func main() {
	cfgPath := os.Getenv("HMDP_CONFIG")
	if cfgPath == "" {
		cfgPath = "configs/app.yaml"
	}
	// 加载配置
	cfg := config.MustLoad(cfgPath)
	log, err := logger.New(cfg.Logging.Level)
	if err != nil {
		panic(err)
	}
	defer log.Sync()
	log.Info("loaded config", zap.String("path", cfgPath))

	// 初始化 MySQL
	db, err := data.NewMySQL(cfg.MySQL, log)
	if err != nil {
		log.Fatal("mysql init failed", zap.Error(err))
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("mysql db handle", zap.Error(err))
	}
	defer sqlDB.Close()
	log.Info("connected to mysql")

	// 初始化 Redis
	redisClient := data.NewRedis(cfg.Redis)
	if err := data.Ping(context.Background(), redisClient); err != nil {
		log.Fatal("redis ping failed", zap.Error(err))
	}
	defer redisClient.Close()
	log.Info("connected to redis", zap.String("addr", cfg.Redis.Addr))

	// 构建 Service Registry
	services := service.NewRegistry(db, redisClient)

	// 初始化 Gin 引擎
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(middleware.ErrorHandler(log))

	uploadDir := cfg.App.ImageUploadDir
	if uploadDir == "" {
		uploadDir = utils.IMAGE_UPLOAD_DIR
	}
	log.Info("configured upload directory", zap.String("path", uploadDir))
	router.RegisterRoutes(engine, services, uploadDir, redisClient)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: engine,
	}
	// 启动 HTTP 服务（异步）
	go func() {
		log.Info("starting http server", zap.String("addr", addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server run failed", zap.Error(err))
		}
	}()

	// 监听系统信号，执行优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Fatal("server shutdown failed", zap.Error(err))
	}
	log.Info("server exited")
}
