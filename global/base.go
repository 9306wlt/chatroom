package global

import "chatroom/define"

var (
	EnteringChannel = make(chan *define.User)
	LeavingChannel  = make(chan *define.User)
	MessageChannel  = make(chan define.Message, 8)
)
