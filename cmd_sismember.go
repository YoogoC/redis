package redis

import (
	"fmt"

	"github.com/redis-go/redcon"
)

func SIsMembersCommand(c *Client, cmd redcon.Command) {
	if len(cmd.Args) != 3 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, "sismember"))
		return
	}

	key := string(cmd.Args[1])

	db := c.Db()
	i := db.GetOrExpire(&key, true)
	if i == nil {
		c.Conn().WriteInt(0)
		return
	} else if i.Type() != SetType {
		c.Conn().WriteError(fmt.Sprintf("%s: key is a %s not a %s", WrongTypeErr, i.TypeFancy(), SetTypeFancy))
		return
	}

	s, ok := i.(*Set)
	if !ok {
		c.Conn().WriteError("sismember error!")
		return
	}

	if s.Contains(string(cmd.Args[2])) {
		c.Conn().WriteInt(1)
	} else {
		c.Conn().WriteInt(0)
	}
}
