package redis

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

// redis server
var r = Default()

var addr = ":6380"

var ctx = context.Background()

// redis client
var c = redis.NewClient(&redis.Options{
	Addr: addr,
})

func run() {
	err := r.Run(addr)
	if err != nil {
		panic(err)
	}
}

func init() {
	go run()
}

func TestPingCommand(t *testing.T) {
	s, err := c.Ping(ctx).Result()
	assert.Equal(t, "PONG", s)
	assert.NoError(t, err)

	pingCmd := redis.NewStringCmd(ctx, "ping", "Hello,", "redis server!")
	_ = c.Process(ctx, pingCmd)
	s, err = pingCmd.Result()
	assert.Equal(t, "Hello, redis server!", s)
	assert.NoError(t, err)
}

func TestSetCommand(t *testing.T) {
	s, err := c.Set(ctx, "k", "v", 0).Result()
	assert.Equal(t, "OK", s)
	assert.NoError(t, err)

	s, err = c.Set(ctx, "k2", nil, 0).Result()
	assert.Equal(t, "OK", s)
	assert.NoError(t, err)

	s, err = c.Set(ctx, "k3", "v", 1*time.Hour).Result()
	assert.Equal(t, "OK", s)
	assert.NoError(t, err)
}

func TestGetCommand(t *testing.T) {
	_, _ = c.Set(ctx, "k", "v", 0).Result()
	s, err := c.Get(ctx, "k").Result()
	assert.Equal(t, "v", s)
	assert.NoError(t, err)
}

func TestDelCommand(t *testing.T) {
	_, _ = c.Set(ctx, "k", "v", 0).Result()
	_, _ = c.Set(ctx, "k3", "v", 1*time.Hour).Result()
	i, err := c.Del(ctx, "k", "k3").Result()
	assert.Equal(t, i, int64(2))
	assert.NoError(t, err)

	i, err = c.Del(ctx, "abc").Result()
	assert.Zero(t, i)
	assert.NoError(t, err)
}

func TestTtlCommand(t *testing.T) {
	s, err := c.Set(ctx, "aKey", "hey", 1*time.Minute).Result()
	assert.Equal(t, "OK", s)
	assert.NoError(t, err)
	s, err = c.Set(ctx, "bKey", "hallo", 0).Result()
	assert.Equal(t, "OK", s)
	assert.NoError(t, err)

	ttl, err := c.TTL(ctx, "aKey").Result()
	assert.True(t, ttl.Seconds() > 55 && ttl.Seconds() < 61, "ttl: %d", ttl)
	assert.NoError(t, err)

	ttl, err = c.TTL(ctx, "none").Result()
	assert.Equal(t, time.Duration(-2000000000), ttl)
	assert.NoError(t, err)

	ttl, err = c.TTL(ctx, "bKey").Result()
	assert.NoError(t, err)
	assert.Equal(t, time.Duration(-1000000000), ttl)
}

func TestExpiry(t *testing.T) {
	s, err := c.Set(ctx, "x", "val", 10*time.Millisecond).Result()
	assert.NoError(t, err)
	assert.Equal(t, "OK", s)

	time.Sleep(10 * time.Millisecond)

	s, err = c.Get(ctx, "x").Result()
	assert.Equal(t, "", s)
	assert.Error(t, err)
}
