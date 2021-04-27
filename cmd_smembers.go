package redis

import (
	"fmt"

	"github.com/tidwall/redcon"
)

func SMembersCommand(c *Client, cmd redcon.Command) {
	if len(cmd.Args) != 2 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, "smembers"))
		return
	}

	key := string(cmd.Args[1])

	db := c.Db()
	i := db.GetOrExpire(&key, true)
	if i == nil {
		c.Conn().WriteArray(0)
		return
	} else if i.Type() != SetType {
		c.Conn().WriteError(fmt.Sprintf("%s: key is a %s not a %s", WrongTypeErr, i.TypeFancy(), SetTypeFancy))
		return
	}

	s, ok := i.(*Set)
	if !ok {
		c.Conn().WriteError("smembers error!")
		return
	}

	c.Conn().WriteArray(s.Len())

	for obj := range s.Iter() {
		c.Conn().WriteBulkString(obj.(string))
	}
}
