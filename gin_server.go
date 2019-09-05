package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

var upGrader = websocket.Upgrader{}

func main() {
	bindAddress := ":8088"
	r := gin.Default()
	r.GET("/connect", connect)
	r.Run(bindAddress)
}

func connect(c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	reqMutex := sync.Mutex{}

	go func() {
		ticker := time.NewTicker(time.Second * 1)
		defer ticker.Stop()

		for {
			select {
			case t := <-ticker.C:
				fmt.Println("write 1 msg:", t.String())
				func() {
					reqMutex.Lock()
					defer reqMutex.Unlock()
					err := ws.WriteMessage(websocket.TextMessage, []byte("1111 "+t.String()))
					if err != nil {
						fmt.Println("write 1 err:", err)
						return
					}
				}()
			}
		}
	}()
	go func() {
		ticker := time.NewTicker(time.Second * 2)
		defer ticker.Stop()

		for {
			select {
			case t := <-ticker.C:
				fmt.Println("write 2 msg:", t.String())
				func() {
					reqMutex.Lock()
					defer reqMutex.Unlock()
					err := ws.WriteMessage(websocket.TextMessage, []byte("2222 "+t.String()))
					if err != nil {
						fmt.Println("write 2 err:", err)
						return
					}
				}()
			}
		}
	}()

	for {
		mt, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("read err:", err)
			break
		}
		fmt.Println("receive:", string(msg))
		err = ws.WriteMessage(mt, []byte("ok: "+string(msg)))
		if err != nil {
			fmt.Println("write err:", err)
			break
		}
	}
}
