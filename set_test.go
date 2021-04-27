package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSAddCommand(t *testing.T) {
	s, err := c.SAdd(ctx, "k1", "1x", "2x", "3x").Result()
	assert.Equal(t, int64(3), s)
	assert.NoError(t, err)

	s, err = c.SAdd(ctx, "k2", 1, 2, 3).Result()
	assert.Equal(t, int64(3), s)
	assert.NoError(t, err)

	s, err = c.SAdd(ctx, "k3", "1x", "1x", "3x").Result()
	assert.Equal(t, int64(2), s)
	assert.NoError(t, err)
}

func TestSRemCommand(t *testing.T) {

	_, _ = c.SAdd(ctx, "k1", "1x", "2x", "3x").Result()
	s, err := c.SRem(ctx, "k1", "1x", "2x", "3x").Result()
	assert.Equal(t, int64(3), s)
	assert.NoError(t, err)

	_, _ = c.SAdd(ctx, "k2", 1, 2, 3).Result()
	s, err = c.SRem(ctx, "k2", 1, 2).Result()
	assert.Equal(t, int64(2), s)
	assert.NoError(t, err)

	_, _ = c.SAdd(ctx, "k3", "1x", "1x", "3x").Result()
	s, err = c.SRem(ctx, "k3", "1x", "1x", "4x").Result()
	assert.Equal(t, int64(1), s)
	assert.NoError(t, err)
}

func TestSMoveCommand(t *testing.T) {

	_, _ = c.SAdd(ctx, "k1", "1x", "2x", "3x").Result()
	s, err := c.SMove(ctx, "k1", "k11", "1x").Result()
	assert.Equal(t, true, s)
	assert.NoError(t, err)

	_, _ = c.SAdd(ctx, "k2", "1x", "2x", "3x").Result()
	_, _ = c.SAdd(ctx, "k3", "4x", "5x", "6x").Result()
	s, err = c.SMove(ctx, "k2", "k3", "3x").Result()
	assert.Equal(t, true, s)
	assert.NoError(t, err)

	_, _ = c.SAdd(ctx, "k4", "1x", "2x", "3x").Result()
	_, _ = c.SAdd(ctx, "k5", "4x", "5x", "6x").Result()
	s, err = c.SMove(ctx, "k4", "k5", "4x").Result()
	assert.Equal(t, false, s)
	assert.NoError(t, err)

}

func TestSMembersCommand(t *testing.T) {

	_, _ = c.SAdd(ctx, "k1", "1x", "2x", "3x").Result()
	s, err := c.SMembers(ctx, "k1").Result()
	assert.Equal(t, 3, len(s))
	assert.NoError(t, err)

	s, err = c.SMembers(ctx, "k2").Result()
	assert.Equal(t, 0, len(s))
	assert.NoError(t, err)
}

func TestSisMemberCommand(t *testing.T) {

	_, _ = c.SAdd(ctx, "k1", "1x", "2x", "3x").Result()
	s, err := c.SIsMember(ctx, "k1", "1x").Result()
	assert.Equal(t, true, s)
	assert.NoError(t, err)

	s, err = c.SIsMember(ctx, "k1", "4x").Result()
	assert.Equal(t, false, s)
	assert.NoError(t, err)

	s, err = c.SIsMember(ctx, "k2", "4x").Result()
	assert.Equal(t, false, s)
	assert.NoError(t, err)
}

func TestSCardCommand(t *testing.T) {

	_, _ = c.SAdd(ctx, "k1", "1x", "2x", "3x").Result()
	s, err := c.SCard(ctx, "k1").Result()
	assert.Equal(t, int64(3), s)
	assert.NoError(t, err)

	s, err = c.SCard(ctx, "k2").Result()
	assert.Equal(t, int64(0), s)
	assert.NoError(t, err)
}
