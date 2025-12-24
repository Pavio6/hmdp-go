local stockKey = KEYS[1]
local orderSetKey = KEYS[2]
local streamKey = KEYS[3]
local userId = ARGV[1]
local voucherId = ARGV[2]
local orderId = ARGV[3]
-- 获取voucher的库存值
local stock = tonumber(redis.call("get", stockKey))
-- 判断库存是否存在或已小于0
if not stock or stock <= 0 then
  return 1
end
-- 检查orderSetKey中是否已包含userId 用来判断用户是否重复下单
if redis.call("sismember", orderSetKey, userId) == 1 then
  return 2
end
-- 库存值减一
redis.call("decr", stockKey)
-- 将userId添加到集合中
redis.call("sadd", orderSetKey, userId)
-- 向Stream中追加一条消息 ID自动生成 包含userId voucherId orderId 三个字段
redis.call("xadd", streamKey, "*", "userId", userId, "voucherId", voucherId, "orderId", orderId)
return 0
