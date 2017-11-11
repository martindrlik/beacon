package main

import "time"

type Message struct {
	Text string
	Time time.Time
}

type rawMessage struct {
	sender   string
	receiver string
	message  Message

	response chan string
}

var (
	messageCh = make(chan rawMessage)
)

func proc(rm rawMessage) {
	defer close(rm.response)
	u, ok := users[rm.receiver]
	if !ok {
		rm.response <- "unknown receiver"
		return
	}
	u.Undelivered[rm.sender] = append(u.Undelivered[rm.sender], rm.message)
}

func init() {
	go func() {
		for rm := range messageCh {
			proc(rm)
		}
	}()
}
