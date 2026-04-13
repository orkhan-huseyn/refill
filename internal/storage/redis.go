package storage

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(url string) RedisStore {
	opt, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opt)
	return RedisStore{
		client: client,
	}
}

var luaScript = `
local key = KEYS[1]
local max_tokens = tonumber(ARGV[1])
local refill_rate = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local amount = tonumber(ARGV[4])

local data = redis.call('HGETALL', key)
local tokens = max_tokens
local last_refill = now

if #data > 0 then
  local fields = {}
  for i = 1, #data, 2 do
    fields[data[i]] = data[i + 1]
  end
  tokens = tonumber(fields['tokens']) or max_tokens
  last_refill = tonumber(fields['last_refill']) or now
end

local elapsed = now - last_refill
local new_tokens = elapsed * refill_rate
tokens = math.min(max_tokens, tokens + new_tokens)

local allowed = 0
local remaining = tokens

if tokens >= amount then
  tokens = tokens - amount
  remaining = tokens
  allowed = 1
end

local retry_after = 0
if allowed == 0 then
    retry_after = (1 - tokens) / refill_rate
end

local reset_time = now + ((max_tokens - tokens) / refill_rate)

redis.call('HSET', key, 'tokens', tostring(tokens), 'last_refill', tostring(now))
redis.call('EXPIRE', key, math.ceil(max_tokens / refill_rate) + 1)

return { 
	tostring(allowed), 
	tostring(math.floor(remaining)), 
	tostring(retry_after), 
	tostring(reset_time) 
}
`

func (s RedisStore) Take(ctx context.Context, key string, amount int) (RateLimitResult, error) {
	var res RateLimitResult

	keys := []string{key}

	// TODO: where to get those arguments? hardcoded doesn't look good
	refillRate := 0.5
	maxTokens := 5.0
	now := float64(time.Now().UnixNano() / 1e9)

	val, err := s.client.Eval(ctx, luaScript, keys, maxTokens, refillRate, now, amount).StringSlice()
	if err != nil {
		return res, err
	}

	remaining, _ := strconv.ParseFloat(val[1], 64)
	retryAfterSecs, _ := strconv.ParseFloat(val[2], 64)
	resetTimeSecs, _ := strconv.ParseFloat(val[3], 64)

	res.Allowed = val[0] == "1"
	res.Limit = int(maxTokens)
	res.Remaining = int(remaining)
	res.RetryAfter = time.Duration(retryAfterSecs * float64(time.Second))
	res.ResetTime = time.Unix(0, int64(resetTimeSecs*1e9))

	return res, nil
}
