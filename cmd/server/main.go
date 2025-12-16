package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"hmdp-backend/internal/config"
	"hmdp-backend/internal/data"
	"hmdp-backend/internal/handler"
	"hmdp-backend/internal/middleware"
	"hmdp-backend/internal/service"
	"hmdp-backend/internal/utils"
	"hmdp-backend/pkg/logger"
)

func main() {
	cfgPath := os.Getenv("HMDP_CONFIG")
	if cfgPath == "" {
		cfgPath = "configs/app.yaml"
	}
	cfg := config.MustLoad(cfgPath)
	log, err := logger.New(cfg.Logging.Level)
	if err != nil {
		panic(err)
	}
	defer log.Sync()
	log.Info("loaded config", zap.String("path", cfgPath))

	log.Info("connecting to mysql")
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

	redisClient := data.NewRedis(cfg.Redis)
	if err := data.Ping(context.Background(), redisClient); err != nil {
		log.Fatal("redis ping failed", zap.Error(err))
	}
	defer redisClient.Close()
	log.Info("connected to redis", zap.String("addr", cfg.Redis.Addr))

	services := service.NewRegistry(db)

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(middleware.ErrorHandler(log))

	shopHandler := handler.NewShopHandler(services.Shop)
	shopTypeHandler := handler.NewShopTypeHandler(services.ShopType)
	voucherHandler := handler.NewVoucherHandler(services.Voucher)
	blogHandler := handler.NewBlogHandler(services.Blog, services.User)
	uploadDir := cfg.App.ImageUploadDir
	if uploadDir == "" {
		uploadDir = utils.IMAGE_UPLOAD_DIR
	}
	uploadHandler := handler.NewUploadHandler(uploadDir)
	log.Info("configured upload directory", zap.String("path", uploadDir))
	userHandler := handler.NewUserHandler(services.User, services.UserInfo)
	voucherOrderHandler := handler.NewVoucherOrderHandler()
	blogCommentsHandler := handler.NewBlogCommentsHandler()
	followHandler := handler.NewFollowHandler()

	shopHandler.RegisterRoutes(engine)
	shopTypeHandler.RegisterRoutes(engine)
	voucherHandler.RegisterRoutes(engine)
	blogHandler.RegisterRoutes(engine)
	uploadHandler.RegisterRoutes(engine)
	userHandler.RegisterRoutes(engine)
	voucherOrderHandler.RegisterRoutes(engine)
	blogCommentsHandler.RegisterRoutes(engine)
	followHandler.RegisterRoutes(engine)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Info("starting http server", zap.String("addr", addr))
	if err := engine.Run(addr); err != nil {
		log.Fatal("server run failed", zap.Error(err))
	}
}
