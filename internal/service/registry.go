package service

import (
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"hmdp-backend/internal/config"
	"hmdp-backend/internal/observability"
	"hmdp-backend/internal/utils"
)

// Registry 聚合全部业务 Service，方便注入 handler
type Registry struct {
	Blog           *BlogService
	Shop           *ShopService
	ShopType       *ShopTypeService
	Voucher        *VoucherService
	SeckillVoucher *SeckillVoucherService
	User           *UserService
	VoucherOrder   *VoucherOrderService
	Follow         *FollowService
}

// NewRegistry 构造服务注册中心
func NewRegistry(
	db *gorm.DB,
	rdb *redis.Client,
	kafkaWriter *kafka.Writer,
	kafkaRetryWriter *kafka.Writer,
	kafkaDLQWriter *kafka.Writer,
	cacheInvalidateWriter *kafka.Writer,
	cacheInvalidateDLQWriter *kafka.Writer,
	kafkaReader *kafka.Reader,
	kafkaRetryReader *kafka.Reader,
	kafkaDLQReader *kafka.Reader,
	cacheInvalidateReader *kafka.Reader,
	cacheInvalidateDLQReader *kafka.Reader,
	smtpCfg utils.SMTPConfig,
	shopCacheCfg config.ShopCacheConfig,
	seckillMetrics *observability.SeckillMetrics,
	log *zap.Logger,
) *Registry {
	if log == nil {
		log = zap.NewNop()
	}
	seckillSvc := NewSeckillVoucherService(db)
	followSvc := NewFollowService(db, rdb)
	return &Registry{
		Blog:           NewBlogService(db, rdb, followSvc),
		Shop:           NewShopService(db, rdb, cacheInvalidateWriter, cacheInvalidateDLQWriter, cacheInvalidateReader, cacheInvalidateDLQReader, smtpCfg, shopCacheCfg, log),
		ShopType:       NewShopTypeService(db, rdb),
		Voucher:        NewVoucherService(db, seckillSvc, rdb),
		SeckillVoucher: seckillSvc,
		User:           NewUserService(db, rdb),
		VoucherOrder:   NewVoucherOrderService(db, rdb, kafkaWriter, kafkaRetryWriter, kafkaDLQWriter, kafkaReader, kafkaRetryReader, kafkaDLQReader, smtpCfg, seckillMetrics, log),
		Follow:         followSvc,
	}
}
