package service

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"hmdp-backend/internal/model"
	"hmdp-backend/internal/utils"
)

const (
	stockKeyFmt = "seckill:stock:vid:%d"
	orderSetFmt = "order:vid:%d"
	streamKey   = "stream.orders"
	streamGroup = "g1"
	streamCons  = "c1"
)

//go:embed seckill.lua
var seckillLuaSource string

// VoucherOrderService 处理秒杀下单逻辑
type VoucherOrderService struct {
	db         *gorm.DB
	rdb        *redis.Client
	idWorker   *utils.RedisIdWorker
	seckillLua *redis.Script
	streamKey  string
}

func NewVoucherOrderService(db *gorm.DB, rdb *redis.Client) *VoucherOrderService {
	svc := &VoucherOrderService{
		db:         db,
		rdb:        rdb,
		idWorker:   utils.NewRedisIdWorker(rdb),
		streamKey:  streamKey,
		seckillLua: redis.NewScript(seckillLuaSource),
	}
	// 初始化消费组
	_ = rdb.XGroupCreateMkStream(context.Background(), svc.streamKey, streamGroup, "$").Err()

	go svc.consumeOrders(context.Background())
	return svc
}

// Seckill 下单处理：校验时间/库存，扣减库存后创建订单
func (s *VoucherOrderService) Seckill(ctx context.Context, voucherID, userID int64) (int64, error) {
	var info struct {
		ID        int64
		BeginTime time.Time
		EndTime   time.Time
		Stock     int
		Status    int
	}
	// 查询秒杀券信息
	err := s.db.WithContext(ctx).Table("tb_voucher AS v").
		Select("v.id, v.status, sv.begin_time, sv.end_time, sv.stock").
		Joins("LEFT JOIN tb_seckill_voucher sv ON v.id = sv.voucher_id").
		Where("v.id = ?", voucherID).
		Take(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, errors.New("优惠券不存在")
	}
	if err != nil {
		return 0, err
	}
	if info.Status != 1 {
		return 0, errors.New("优惠券已下架或过期")
	}

	now := time.Now()
	if now.Before(info.BeginTime) {
		return 0, errors.New("秒杀尚未开始")
	}
	if now.After(info.EndTime) {
		return 0, errors.New("秒杀已结束")
	}
	// 库存不足直接返回
	if info.Stock <= 0 {
		return 0, errors.New("库存不足")
	}

	// 生成订单ID
	orderID, err := s.idWorker.NextId(ctx, "order")
	if err != nil {
		return 0, err
	}

	stockKey := fmt.Sprintf(stockKeyFmt, voucherID)
	orderSetKey := fmt.Sprintf(orderSetFmt, voucherID)
	// 执行 Lua 脚本，完成库存校验与扣减、用户下单资格校验与标记、订单写入 Stream
	res, err := s.seckillLua.Run(ctx, s.rdb, []string{stockKey, orderSetKey, s.streamKey},
		userID, voucherID, orderID).Int()
	if err != nil {
		return 0, err
	}

	switch res {
	case 0:
		// Lua 校验成功，消息已写入 Stream，异步创建订单
		return orderID, nil
	case 1:
		return 0, errors.New("库存不足")
	case 2:
		return 0, errors.New("每人限购一单")
	default:
		return 0, errors.New("秒杀失败")
	}
}

// consumeOrders 异步创建订单
func (s *VoucherOrderService) consumeOrders(ctx context.Context) {
	for {
		// 优先处理 pending，再拉取新消息
		msgs, err := s.readStream(ctx, "0") // pending
		if err != nil {
			log.Printf("consumeOrders read pending error: %v", err)
			time.Sleep(time.Second)
			continue
		}
		if len(msgs) == 0 {
			msgs, err = s.readStream(ctx, ">") // new messages
			if err != nil {
				log.Printf("consumeOrders read new error: %v", err)
				time.Sleep(time.Second)
				continue
			}
			if len(msgs) == 0 {
				continue
			}
		}

		for _, m := range msgs {
			uid, vid, oid, err := parseOrderMsg(m)
			if err != nil {
				log.Printf("consumeOrders parse msg error: %v msg=%v", err, m)
				_ = s.rdb.XAck(ctx, s.streamKey, streamGroup, m.ID).Err()
				continue
			}

			nowTime := time.Now()
			// 事务内扣减 DB 库存并创建订单，确保数据库内一致性
			if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
				res := tx.Model(&model.SeckillVoucher{}).
					Where("voucher_id = ? AND stock > 0", vid).
					Update("stock", gorm.Expr("stock - 1"))
				if res.Error != nil {
					return res.Error
				}
				if res.RowsAffected == 0 {
					return errors.New("db stock not enough")
				}
				order := &model.VoucherOrder{
					ID:         oid,
					UserID:     uid,
					VoucherID:  vid,
					CreateTime: nowTime,
					UpdateTime: nowTime,
				}
				return tx.Create(order).Error
			}); err != nil {
				log.Printf("consumeOrders txn error: %v msgID=%s", err, m.ID)
				continue
			}
			// 处理成功，ACK
			if err := s.rdb.XAck(ctx, s.streamKey, streamGroup, m.ID).Err(); err != nil {
				log.Printf("consumeOrders ack error: %v msgID=%s", err, m.ID)
			}
		}
	}
}

// readStream 读取 Stream 消息，start 为 "0" 读取 pending，为 ">" 读取新消息
func (s *VoucherOrderService) readStream(ctx context.Context, start string) ([]redis.XMessage, error) {
	res, err := s.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    streamGroup,
		Consumer: streamCons,
		Streams:  []string{s.streamKey, start},
		Count:    10,
		Block:    time.Second,
	}).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	if len(res) == 0 || len(res[0].Messages) == 0 {
		return nil, nil
	}
	return res[0].Messages, nil
}

// parseOrderMsg 解析 Stream 消息中的订单字段
func parseOrderMsg(m redis.XMessage) (int64, int64, int64, error) {
	parse := func(v interface{}) (int64, error) {
		switch val := v.(type) {
		case string:
			return strconv.ParseInt(val, 10, 64)
		case float64:
			return int64(val), nil
		case json.Number:
			return val.Int64()
		default:
			return 0, fmt.Errorf("unexpected type %T", v)
		}
	}

	uid, err1 := parse(m.Values["userId"])
	vid, err2 := parse(m.Values["voucherId"])
	oid, err3 := parse(m.Values["orderId"])
	if err1 != nil || err2 != nil || err3 != nil {
		return 0, 0, 0, fmt.Errorf("uidErr=%v vidErr=%v oidErr=%v", err1, err2, err3)
	}
	return uid, vid, oid, nil
}
