package redis

import (
	"fmt"
	"time"

	"github.com/tidwall/redcon"
)

func SMoveCommand(c *Client, cmd redcon.Command) {
	if len(cmd.Args) != 4 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, "smove"))
		return
	}

	srckey := string(cmd.Args[1])
	db := c.Db()
	srcitem := db.GetOrExpire(&srckey, true)
	if srcitem == nil {
		c.Conn().WriteNull()
	} else if srcitem.Type() != SetType {
		c.Conn().WriteError(fmt.Sprintf("%s: key is a %s not a %s", WrongTypeErr, srcitem.TypeFancy(), SetTypeFancy))
		return
	}

	dstkey := string(cmd.Args[2])
	dstitem := db.GetOrExpire(&dstkey, true)
	if dstitem != nil && dstitem.Type() != SetType {
		c.Conn().WriteError(fmt.Sprintf("%s: key is a %s not a %s", WrongTypeErr, dstitem.TypeFancy(), SetTypeFancy))
		return
	}

	if srcitem == dstitem {
		c.Conn().WriteInt(0)
		return
	}

	srcset, sok := srcitem.(*Set)

	if !sok {
		c.Conn().WriteError("smove error!")
		return
	}

	var dstset *Set
	if dstitem == nil {
		dstset = NewSet()
		db.Set(&dstkey, dstset, false, time.Time{})
	} else {
		dstset = dstitem.(*Set)
	}

	member := string(cmd.Args[3])
	c.Redis().Mu().Lock()
	if srcset.Remove(member) {
		dstset.Set(member)
	} else {
		c.Redis().Mu().Unlock()
		c.Conn().WriteInt(0)
		return
	}
	c.Redis().Mu().Unlock()

	if srcset.Len() == 0 {
		db.Delete(&srckey)
	}

	c.Conn().WriteInt(1)
}
