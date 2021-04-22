package redis

import (
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

// redis server
var r = Default()

var addr = ":6380"

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
	s, err := c.Ping().Result()
	assert.Equal(t, "PONG", s)
	assert.NoError(t, err)

	pingCmd := redis.NewStringCmd("ping", "Hello,", "redis server!")
	_ = c.Process(pingCmd)
	s, err = pingCmd.Result()
	assert.Equal(t, "Hello, redis server!", s)
	assert.NoError(t, err)
}

func TestSetCommand(t *testing.T) {
	s, err := c.Set("k", "v", 0).Result()
	assert.Equal(t, "OK", s)
	assert.NoError(t, err)

	s, err = c.Set("k2", nil, 0).Result()
	assert.Equal(t, "OK", s)
	assert.NoError(t, err)

	s, err = c.Set("k3", "v", 1*time.Hour).Result()
	assert.Equal(t, "OK", s)
	assert.NoError(t, err)
}

func TestGetCommand(t *testing.T) {
	_, _ = c.Set("k", "v", 0).Result()
	s, err := c.Get("k").Result()
	assert.Equal(t, "v", s)
	assert.NoError(t, err)
}

func TestDelCommand(t *testing.T) {
	_, _ = c.Set("k", "v", 0).Result()
	_, _ = c.Set("k3", "v", 1*time.Hour).Result()
	i, err := c.Del("k", "k3").Result()
	assert.Equal(t, i, int64(2))
	assert.NoError(t, err)

	i, err = c.Del("abc").Result()
	assert.Zero(t, i)
	assert.NoError(t, err)
}

func TestTtlCommand(t *testing.T) {
	s, err := c.Set("aKey", "hey", 1*time.Minute).Result()
	assert.Equal(t, "OK", s)
	assert.NoError(t, err)
	s, err = c.Set("bKey", "hallo", 0).Result()
	assert.Equal(t, "OK", s)
	assert.NoError(t, err)

	ttl, err := c.TTL("aKey").Result()
	assert.True(t, ttl.Seconds() > 55 && ttl.Seconds() < 61, "ttl: %d", ttl)
	assert.NoError(t, err)

	ttl, err = c.TTL("none").Result()
	assert.Equal(t, time.Duration(-2000000000), ttl)
	assert.NoError(t, err)

	ttl, err = c.TTL("bKey").Result()
	assert.NoError(t, err)
	assert.Equal(t, time.Duration(-1000000000), ttl)
}

func TestExpiry(t *testing.T) {
	s, err := c.Set("x", "val", 10*time.Millisecond).Result()
	assert.NoError(t, err)
	assert.Equal(t, "OK", s)

	time.Sleep(10 * time.Millisecond)

	s, err = c.Get("x").Result()
	assert.Equal(t, "", s)
	assert.Error(t, err)
}

func TestSAddCommand(t *testing.T) {
	s, err := c.SAdd("k1", "1x", "2x", "3x").Result()
	assert.Equal(t, int64(3), s)
	assert.NoError(t, err)

	s, err = c.SAdd("k2", 1, 2, 3).Result()
	assert.Equal(t, int64(3), s)
	assert.NoError(t, err)

	s, err = c.SAdd("k3", "1x", "1x", "3x").Result()
	assert.Equal(t, int64(2), s)
	assert.NoError(t, err)
}

func TestSRemCommand(t *testing.T) {

	_, _ = c.SAdd("k1", "1x", "2x", "3x").Result()
	s, err := c.SRem("k1", "1x", "2x", "3x").Result()
	assert.Equal(t, int64(3), s)
	assert.NoError(t, err)

	_, _ = c.SAdd("k2", 1, 2, 3).Result()
	s, err = c.SRem("k2", 1, 2).Result()
	assert.Equal(t, int64(2), s)
	assert.NoError(t, err)

	_, _ = c.SAdd("k3", "1x", "1x", "3x").Result()
	s, err = c.SRem("k3", "1x", "1x", "4x").Result()
	assert.Equal(t, int64(1), s)
	assert.NoError(t, err)
}
