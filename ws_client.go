package main

import (
	"client-test/prototest"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"net/url"
)

func main() {
	u := url.URL{Scheme: "ws", Host: "101.251.201.34:8088", Path: "/connect"}
	fmt.Println("connecting to", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("connect to server %s success\n", u.String())
	defer c.Close()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			fmt.Println("read err:", err)
			return
		}
		//time.Sleep(time.Second * 2)
		p := &prototest.Person{}
		proto.Unmarshal(msg, p)
		fmt.Println("from server msg:", p)
	}

	//done := make(chan struct{})
	//
	//go func() {
	//	defer close(done)
	//	for {
	//		_, msg, err := c.ReadMessage()
	//		if err != nil {
	//			fmt.Println("read err:", err)
	//			return
	//		}
	//		fmt.Println("from server msg:", string(msg))
	//	}
	//}()
	//
	//ticker := time.NewTicker(time.Second * 5)
	//defer ticker.Stop()
	//
	//for {
	//	select {
	//	case <-done:
	//		return
	//	case t := <-ticker.C:
	//		fmt.Println("write msg:", t.String())
	//		err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
	//		if err != nil {
	//			fmt.Println("write err:", err)
	//			return
	//		}
	//	}
	//}
}
