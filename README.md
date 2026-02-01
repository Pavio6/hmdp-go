## 秒杀功能说明（Seckill）

本项目的秒杀流程使用 **Redis Lua + Kafka 异步落库 + 重试/DLQ** 组合，以保证高并发下的性能与一致性。

### 功能概述
- Redis Lua 脚本负责校验库存、是否重复下单、并扣减库存（原子操作）。
- 通过 Kafka 异步写订单，提升接口吞吐与响应速度。
- 消费端写库失败后进入重试队列；超过最大次数进入 DLQ（死信队列），可人工处理或告警。

### 业务流程（请求链路）
1. **客户端请求** `/voucher-order/seckill/{voucherId}`。
2. **Redis Lua** 原子校验与扣减：
   - 库存不足 → 直接失败
   - 重复下单 → 直接失败
   - 成功 → 返回订单 ID
3. **Kafka 生产**：将订单消息写入主 Topic（按 voucherId 分区）。
4. **Kafka 消费**：
   - 事务内创建订单
   - 订单创建成功后再扣减 DB 库存（防重复消费）
5. **失败处理**：
   - 可重试错误 → 写入 retry topic，按指数退避延迟处理
   - 超过最大次数 → 写入 DLQ，发送邮件告警（可选）

### 关键设计点
- **防超卖**：DB 扣减用 `UPDATE ... SET stock = stock - 1 WHERE stock > 0` 原子条件更新。
- **幂等**：订单表唯一约束，重复消费会触发 duplicate key，直接返回成功避免重复扣库存。
- **分区有序**：Kafka 使用 `voucherId` 作为 key，同券消息落同分区。
- **重试退避**：指数退避（1s, 2s, 4s...，最大 30s），超过次数进入 DLQ。

### 代码位置
- Lua 脚本：`internal/service/seckill.lua`
- 订单生产/消费/重试逻辑：`internal/service/voucher_order_service.go`
- ID 生成：`internal/utils/redisId_worker.go`

### 本地运行依赖
需要本地或容器内提供：
- MySQL
- Redis
- Kafka

配置文件：`configs/app.yaml`

### 测试方法

#### 1) 单次下单（curl）
```bash
TOKEN="替换成你的token"
VOUCHER_ID=12
curl -X POST "http://127.0.0.1:8081/voucher-order/seckill/${VOUCHER_ID}" \
  -H "authorization: ${TOKEN}"
```

#### 2) 压测秒杀（k6）
1. 生成测试 token：
```bash
go run cmd/gen_tokens/main.go -in hmdp_tb_user.csv -out tokens.csv -redis 127.0.0.1:6379 -db 0
```

2. 启动服务：
```bash
rm -f server.log
go run cmd/server/main.go > server.log 2>&1 &
```

3. 运行压测：
```bash
k6 run -e BASE_URL=http://127.0.0.1:8081 \
  -e VOUCHER_ID=12 \
  -e TOKENS_FILE=../../tokens.csv \
  -e RAMP_WINDOW=10s \
  scripts/k6/seckill.js
```

4. 查看消费落库日志：
```bash
rg -n "handleConsume success|handleConsume failed" server.log
```

#### 3) 测试重试与 DLQ（计数开关）
设置环境变量 `FORCE_SECKILL_CONSUME_FAIL_COUNT=n`，当 `RetryCount < n` 时强制失败。

示例：失败 3 次后成功
```bash
rm -f server.log
FORCE_SECKILL_CONSUME_FAIL_COUNT=3 go run cmd/server/main.go > server.log 2>&1 &
```

观察重试/DLQ：
```bash
rg -n "publish to retry|publish to dlq|handleConsume failed" server.log
```

> 提示：若 `n > maxRetryCount`（当前为 3），消息会进入 DLQ。
