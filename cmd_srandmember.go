package redis

import (
	"fmt"
	"math"
	"strconv"

	"github.com/tidwall/redcon"
)

func SRandMemberCommand(c *Client, cmd redcon.Command) {
	if len(cmd.Args) < 2 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, "srandmember"))
		return
	}

	if len(cmd.Args) == 3 {
		SRandMemberWithCountCommand(c, cmd)
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
		c.Conn().WriteError("srandmember error!")
		return
	}

	c.Conn().WriteString(s.GetOne().(string))
}

func SRandMemberWithCountCommand(c *Client, cmd redcon.Command) {
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
		c.Conn().WriteError("srandmember error!")
		return
	}

	count, err := strconv.Atoi(string(cmd.Args[2]))
	if err != nil {
		c.Conn().WriteError("srandmember count argument error!")
		return
	}

	if count < 0 {
		countAbs := int(math.Abs(float64(count)))
		c.Conn().WriteArray(countAbs)
		for i := 0; i < countAbs; i++ {
			c.Conn().WriteBulkString(s.GetOne().(string))
		}
	} else {
		if count > s.Len() {
			count = s.Len()
		}

		c.Conn().WriteArray(count)
		iter := s.Iter()
		for i := 0; i < count; i++ {
			v := <-iter
			c.Conn().WriteBulkString(v.(string))
		}
	}
}
