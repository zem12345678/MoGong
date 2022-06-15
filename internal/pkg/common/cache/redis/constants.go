package redis

const (
	// ScriptDeleteLock 释放redis并发锁 lua脚本 判断value为本次锁的value才释放
	ScriptDeleteLock = `
if redis.call("get", KEYS[1]) == ARGV[1] then
    return redis.call("del",KEYS[1])
else
    return 0
end
`
	// lua脚本实现令牌桶算法限流
	ScriptTokenLimit = `
local rateLimit = redis.pcall('HMGET',KEYS[1],'lastTime','tokens')
local lastTime = rateLimit[1]
local tokens = tonumber(rateLimit[2])
local capacity = tonumber(ARGV[1])
local rate = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
if tokens == nil then
  tokens = capacity
else
  local deltaTokens = math.floor((now-lastTime)*rate)
  tokens = tokens+deltaTokens
  if tokens>capacity then
    tokens = capacity
  end
end
local result = false
lastTime = now
if(tokens>0) then
  result = true
  tokens = tokens-1
end
redis.call('HMSET',KEYS[1],'lastTime',lastTime,'tokens',tokens)
return result
`
)
