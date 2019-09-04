package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
	"time"
)

func main() {
	u := url.URL{Scheme: "ws", Host: "localhost:8088", Path: "/connect"}
	fmt.Println("connecting to", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				fmt.Println("read err:", err)
				return
			}
			fmt.Println("from server msg:", string(msg))
		}
	}()

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			fmt.Println("write msg:", t.String())
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				fmt.Println("write err:", err)
				return
			}
		}
	}
}
