package service

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"

	"hmdp-backend/internal/config"
	"hmdp-backend/internal/data"
	"hmdp-backend/internal/model"
	"hmdp-backend/internal/utils"
	"hmdp-backend/pkg/logger"
)

// TestWarmShopCacheWithLogicalExpire 测试方法：将指定 ID 的商铺写入 Redis 并设置逻辑过期时间
func TestWarmShopCacheWithLogicalExpire(t *testing.T) {

	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = filepath.Join("..", "..", "configs", "app.yaml")
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	log, err := logger.New(cfg.Logging.Level)
	if err != nil {
		t.Fatalf("init logger: %v", err)
	}
	db, err := data.NewMySQL(cfg.MySQL, log)
	if err != nil {
		t.Fatalf("init mysql: %v", err)
	}
	rdb := data.NewRedis(cfg.Redis)

	shopID := int64(1)
	if envID := os.Getenv("SHOP_ID"); envID != "" {
		parsed, parseErr := strconv.ParseInt(envID, 10, 64)
		if parseErr != nil {
			t.Fatalf("invalid SHOP_ID: %v", parseErr)
		}
		shopID = parsed
	}

	svc := NewShopService(db, rdb, log)
	key := utils.CACHE_SHOP_KEY + strconv.FormatInt(shopID, 10)
	var shop model.Shop
	if err := db.WithContext(context.Background()).First(&shop, shopID).Error; err != nil {
		t.Fatalf("query shop: %v", err)
	}
	if err := svc.saveShopWithLogicalExpire(key, &shop, time.Duration(utils.CACHE_SHOP_TTL)*time.Minute); err != nil {
		t.Fatalf("save logical expire cache: %v", err)
	}
}

func TestLoadShopData(t *testing.T) {
	ctx := context.Background()

	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = filepath.Join("..", "..", "configs", "app.yaml")
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	log, err := logger.New(cfg.Logging.Level)
	if err != nil {
		t.Fatalf("init logger: %v", err)
	}
	db, err := data.NewMySQL(cfg.MySQL, log)
	if err != nil {
		t.Fatalf("init mysql: %v", err)
	}
	rdb := data.NewRedis(cfg.Redis)

	var shops []model.Shop
	if err := db.WithContext(ctx).Find(&shops).Error; err != nil {
		t.Fatalf("query shops: %v", err)
	}
	if len(shops) == 0 {
		t.Fatalf("no shops returned from database")
	}

	// Group shops by typeId then batch write into Redis GEO sets keyed by typeId.
	grouped := make(map[int64][]*redis.GeoLocation)
	for _, shop := range shops {
		loc := &redis.GeoLocation{
			Name:      strconv.FormatInt(shop.ID, 10),
			Longitude: shop.X,
			Latitude:  shop.Y,
		}
		grouped[shop.TypeID] = append(grouped[shop.TypeID], loc)
	}

	const batchSize = 10
	for typeID, locations := range grouped {
		key := utils.SHOP_GEO_KEY + strconv.FormatInt(typeID, 10)
		if err := rdb.Del(ctx, key).Err(); err != nil {
			t.Fatalf("clear redis key %s: %v", key, err)
		}

		_, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			for start := 0; start < len(locations); start += batchSize {
				end := start + batchSize
				if end > len(locations) {
					end = len(locations)
				}
				if err := pipe.GeoAdd(ctx, key, locations[start:end]...).Err(); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			t.Fatalf("geo add for type %d: %v", typeID, err)
		}

		count, err := rdb.ZCard(ctx, key).Result()
		if err != nil {
			t.Fatalf("zcard %s: %v", key, err)
		}
		if count != int64(len(locations)) {
			t.Fatalf("unexpected member count for key %s: want %d got %d", key, len(locations), count)
		}
	}
}
