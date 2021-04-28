package redis

import (
	"fmt"
	"strconv"

	"github.com/tidwall/redcon"
)

func SPopCommand(c *Client, cmd redcon.Command) {
	if len(cmd.Args) < 2 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, "spop"))
		return
	}

	key := string(cmd.Args[1])

	db := c.Db()
	i := db.GetOrExpire(&key, true)
	if i == nil {
		// c.Conn().WriteNull()
		c.Conn().WriteString("")
		return
	} else if i.Type() != SetType {
		c.Conn().WriteError(fmt.Sprintf("%s: key is a %s not a %s", WrongTypeErr, i.TypeFancy(), SetTypeFancy))
		return
	}

	s, ok := i.(*Set)
	if !ok {
		c.Conn().WriteError("spop error!")
		return
	}

	needRemoveLen := 1
	if len(cmd.Args) > 2 {
		slen := s.Len()

		psize, err := strconv.Atoi(string(cmd.Args[2]))
		if err != nil {
			c.Conn().WriteError("spop count argument error!")
			return
		}

		if slen < psize {
			needRemoveLen = slen
		} else {
			needRemoveLen = psize
		}
	}

	if int(needRemoveLen) > 1 {
		c.Conn().WriteArray(int(needRemoveLen))
		for i := 0; i <= int(needRemoveLen)-1; i++ {
			v, _ := s.Pop().(string)
			c.Conn().WriteBulkString(v)
		}
	} else {
		v, _ := s.Pop().(string)
		c.Conn().WriteString(v)
	}

	if s.Len() == 0 {
		db.Delete(&key)
	}
}
