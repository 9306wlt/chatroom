package main

import (
	"bufio"
	"chatroom/define"
	"chatroom/global"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", ":2020")
	if err != nil {
		panic(err)
	}
	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Panicln(err)
			continue
		}
		go handleConn(conn)
	}
}

func broadcaster() {
	users := make(map[*define.User]struct{})
	for {
		select {
		case user := <-global.EnteringChannel:
			users[user] = struct{}{}
		case user := <-global.LeavingChannel:
			delete(users, user)
			close(user.MessageChannel)
		case msg := <-global.MessageChannel:
			for user := range users {
				if user.ID == msg.OwnerID {
					continue
				}
				user.MessageChannel <- msg.Content
			}
		}
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	user := &define.User{
		ID:             GenUserID(),
		Addr:           conn.RemoteAddr().String(),
		EnterAt:        time.Now(),
		MessageChannel: make(chan string, 8),
	}

	go sendMessage(conn, user.MessageChannel)

	msg := define.Message{
		OwnerID: user.ID,
		Content: "user:`" + strconv.Itoa(user.ID) + "` has enter",
	}

	user.MessageChannel <- "Welcome, " + user.String()
	global.MessageChannel <- msg

	global.EnteringChannel <- user

	input := bufio.NewScanner(conn)
	for input.Scan() {
		msg.Content = strconv.Itoa(user.ID) + ":" + input.Text()
		global.MessageChannel <- msg
	}

	if err := input.Err(); err != nil {
		log.Panicln("读取错误：", err)
	}
	global.LeavingChannel <- user
	msg.Content = "user: `" + strconv.Itoa(user.ID) + "` has left"
	global.MessageChannel <- msg
}

func sendMessage(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

var (
	globalID int
	idLocker sync.Mutex
)

func GenUserID() int {
	idLocker.Lock()
	defer idLocker.Unlock()
	globalID++
	return globalID
}
