package redis

import (
	"fmt"

	"github.com/tidwall/redcon"
)

func SRemCommand(c *Client, cmd redcon.Command) {
	if len(cmd.Args) < 3 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, "srem"))
		return
	}

	key := string(cmd.Args[1])

	db := c.Db()
	i := db.GetOrExpire(&key, true)
	if i == nil {
		c.Conn().WriteNull()
		return
	} else if i.Type() != SetType {
		c.Conn().WriteError(fmt.Sprintf("%s: key is a %s not a %s", WrongTypeErr, i.TypeFancy(), SetTypeFancy))
		return
	}

	s, ok := i.(*Set)
	if !ok {
		c.Conn().WriteError("srem error!")
		return
	}
	var length int
	c.Redis().Mu().Lock()
	for j := 2; j < len(cmd.Args); j++ {
		v := string(cmd.Args[j])
		if s.Remove(v) {
			length++
		}
	}
	c.Redis().Mu().Unlock()

	if s.Len() == 0 {
		db.Delete(&key)
	}

	c.Conn().WriteInt(length)
}
