package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"hmdp-backend/internal/config"
	"hmdp-backend/internal/data"
)

// TestHyperLogLogBatchInsert 验证 HyperLogLog 批量写入 100 万条数据后基数估算是否在误差范围内。
func TestHyperLogLogBatchInsert(t *testing.T) {
	ctx := context.Background()

	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = filepath.Join("..", "..", "configs", "app.yaml")
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	rdb := data.NewRedis(cfg.Redis)

	const (
		key = "test:hll:batch"
		// 整数字面量可以使用下划线分割来提升可读性
		total     = 1_000_000
		batchSize = 1000
	)
	// 清理旧数据
	_ = rdb.Del(ctx, key).Err()

	batch := make([]string, 0, batchSize)
	for i := 0; i < total; i++ {
		batch = append(batch, fmt.Sprintf("member-%d", i))
		if len(batch) == batchSize {
			if err := rdb.PFAdd(ctx, key, batch).Err(); err != nil {
				t.Fatalf("pfadd batch: %v", err)
			}
			batch = batch[:0] // 清空切片以便重用
		}
	}
	if len(batch) > 0 {
		if err := rdb.PFAdd(ctx, key, batch).Err(); err != nil {
			t.Fatalf("pfadd last batch: %v", err)
		}
	}

	count, err := rdb.PFCount(ctx, key).Result()
	if err != nil {
		t.Fatalf("pfcount: %v", err)
	}

	// HyperLogLog 典型误差约 0.81%，这里允许 ±1% 容忍范围
	lower := int64(float64(total) * 0.99)
	upper := int64(float64(total) * 1.01)
	if count < lower || count > upper {
		t.Fatalf("unexpected hll count: %d, expect around %d (±1%%)", count, total)
	}
	t.Logf("HyperLogLog count=%d within tolerance", count)
}
