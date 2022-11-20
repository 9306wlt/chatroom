package define

import (
	"strconv"
	"time"
)

type User struct {
	ID             int
	Addr           string
	EnterAt        time.Time
	MessageChannel chan string
}

func (u User) String() string {
	return u.Addr + ",UID = " + strconv.Itoa(u.ID) + " ,Enter time" + u.EnterAt.Format("2022-11-19:10:46:50")
}

type Message struct {
	OwnerID int
	Content string
}
