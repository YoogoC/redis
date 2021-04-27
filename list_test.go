package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLPushCommand(t *testing.T) {
	i, err := c.LPush(ctx, "lpushkey", "va").Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), i)

	i, err = c.LPush(ctx, "lpushkey", "vb").Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), i)

	i, err = c.LPush(ctx, "lpushkey", "vc", "vd").Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(4), i)

	i, err = c.LPush(ctx, "lpushkey2", "1", "2").Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), i)

	_, err = c.LPush(ctx, "lpush3key").Result()
	assert.Error(t, err)
}

func TestLPopCommand(t *testing.T) {
	s, err := c.LPop(ctx, "lpop1").Result()
	assert.Zero(t, s)
	assert.Error(t, err)

	i, err := c.LPush(ctx, "list", "a", "b").Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), i)

	s, err = c.LPop(ctx, "list").Result()
	assert.NoError(t, err)
	assert.Equal(t, "b", s)

	s, err = c.LPop(ctx, "list").Result()
	assert.NoError(t, err)
	assert.Equal(t, "a", s)

	s, err = c.LPop(ctx, "list").Result()
	assert.Error(t, err)
	assert.Zero(t, s)
}

func TestLRangeCommand(t *testing.T) {
	s, err := c.LRange(ctx, "lrange", 0, 0).Result()
	assert.Error(t, err)
	assert.Zero(t, s)

	sl, err := c.Set(ctx, "works", "esfkjsefj", 0).Result()
	assert.NoError(t, err)
	assert.NotZero(t, sl)
	assert.NotEmpty(t, sl)

	i, err := c.LPush(ctx, "list2", "a", "b").Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), i)
}
