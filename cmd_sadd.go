package redis

import (
	"fmt"
	"time"

	"github.com/tidwall/redcon"
)

func SAddCommand(c *Client, cmd redcon.Command) {
	if len(cmd.Args) < 3 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, "sadd"))
		return
	}

	key := string(cmd.Args[1])

	db := c.Db()
	i := db.GetOrExpire(&key, true)
	if i == nil {
		i = NewSet()
		db.Set(&key, i, false, time.Time{})
	} else if i.Type() != SetType {
		c.Conn().WriteError(fmt.Sprintf("%s: key is a %s not a %s", WrongTypeErr, i.TypeFancy(), SetTypeFancy))
		return
	}

	s, ok := i.(*Set)
	if !ok {
		c.Conn().WriteError("sadd error!")
		return
	}
	var length int
	c.Redis().Mu().Lock()
	for j := 2; j < len(cmd.Args); j++ {
		v := string(cmd.Args[j])
		if s.Set(v) {
			length++
		}
	}
	c.Redis().Mu().Unlock()

	c.Conn().WriteInt(length)
}
