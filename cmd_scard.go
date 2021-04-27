package redis

import (
	"fmt"

	"github.com/tidwall/redcon"
)

func SCard(c *Client, cmd redcon.Command) {
	if len(cmd.Args) != 2 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, "scard"))
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
		c.Conn().WriteError("scard error!")
		return
	}

	c.Conn().WriteInt(s.Len())
}
